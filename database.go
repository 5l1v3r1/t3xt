package main

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
)

var indexFilename = "index.json"

const copyBufferSize = 0x1000
const headLineCount = 5

type Database struct {
	lock  sync.RWMutex
	path  string
	index *index
}

func OpenDatabase(path string) (*Database, error) {
	statInfo, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return createDatabase(path)
	} else if err != nil {
		return nil, err
	}
	if !statInfo.IsDir() {
		return nil, errors.New("DB is not directory: " + path)
	}
	indexData, err := ioutil.ReadFile(filepath.Join(path, indexFilename))
	if err != nil {
		return nil, err
	}
	var index index
	if err := json.Unmarshal(indexData, &index); err != nil {
		return nil, err
	}
	return &Database{path: path, index: &index}, nil
}

// LatestEntries returns a certain number of the latest entries.
func (d *Database) LatestEntries(count int) []DatabaseEntry {
	return d.EntriesBefore(-1, count)
}

// EntriesBefore returns a certain number of entries whose IDs descend from a starting ID.
// If an entry exists in the database for startId, that entry will be included.
// The results are sorted in ascending order by ID.
func (d *Database) EntriesBefore(startId, count int) []DatabaseEntry {
	res := make([]DatabaseEntry, 0, count)
	d.lock.RLock()
	defer d.lock.RUnlock()
	if startId == -1 {
		startId = d.index.CurrentID - 1
	}
	for i := startId; i >= 0 && len(res) < count; i-- {
		if entry, ok := d.index.entryForID(i); ok {
			res = append(res, entry)
		}
	}
	return res
}

// EntriesAfter returns a certain number of entries whose IDs ascend from a starting ID.
// If an entry exists in the database for startId, that entry will be included.
// The results are sorted in ascending order by ID.
func (d *Database) EntriesAfter(startId, count int) []DatabaseEntry {
	res := make([]DatabaseEntry, count)
	remainingCount := len(res)
	d.lock.RLock()
	defer d.lock.RUnlock()
	for i := startId; i < d.index.CurrentID && remainingCount > 0; i++ {
		if entry, ok := d.index.entryForID(i); ok {
			remainingCount--
			res[remainingCount] = entry
		}
	}
	return res[remainingCount:]
}

func (d *Database) OpenEntry(shareID string) (e DatabaseEntry, r io.ReadCloser, err error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	id, ok := d.index.ShareIDToID[shareID]
	if !ok {
		err = errors.New("unknown share ID: " + shareID)
		return
	}
	e, ok = d.index.entryForID(id)
	if !ok {
		err = errors.New("unknown ID: " + strconv.Itoa(id))
		return
	}
	dataPath := filepath.Join(d.path, strconv.Itoa(e.ID))
	r, err = os.Open(dataPath)
	return
}

// Head returns the first few lines of an entry's contents.
func (d *Database) Head(id int) (string, error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	dataPath := filepath.Join(d.path, strconv.Itoa(id))
	f, err := os.Open(dataPath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	lines := make([]string, 0, headLineCount)
	for i := 0; i < headLineCount; i++ {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", err
		}
		if err == io.EOF {
			if line != "" {
				lines = append(lines, line)
			}
			break
		}
		lines = append(lines, line[:len(line)-1])
	}
	return strings.Join(lines, "\n"), nil
}

func (d *Database) CreateEntry(info DatabaseEntry,
	body io.Reader) (entry DatabaseEntry, err error) {
	tempFile, err := ioutil.TempFile("", "t3xt")
	if err != nil {
		return
	}
	lineCount, err := copyAndCountLines(tempFile, body)
	tempFile.Close()
	if err != nil {
		os.Remove(tempFile.Name())
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	entry = info
	entry.ShareID = randomShareID()
	entry.ID = d.index.CurrentID
	entry.LineCount = lineCount
	d.index.CurrentID++

	dataPath := filepath.Join(d.path, strconv.Itoa(entry.ID))
	err = os.Rename(tempFile.Name(), dataPath)
	if err != nil {
		os.Remove(tempFile.Name())
		return
	}
	d.index.IDToEntry[strconv.Itoa(entry.ID)] = entry
	d.index.ShareIDToID[entry.ShareID] = entry.ID
	err = d.saveIndex()
	if err != nil {
		delete(d.index.IDToEntry, strconv.Itoa(entry.ID))
		delete(d.index.ShareIDToID, entry.ShareID)
		os.Remove(dataPath)
	}
	return
}

func (d *Database) DeleteEntry(e DatabaseEntry) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	oldEntry, _ := d.index.entryForID(e.ID)
	delete(d.index.IDToEntry, strconv.Itoa(e.ID))
	delete(d.index.ShareIDToID, e.ShareID)
	if err := d.saveIndex(); err != nil {
		d.index.IDToEntry[strconv.Itoa(e.ID)] = oldEntry
		d.index.ShareIDToID[e.ShareID] = oldEntry.ID
		return err
	}
	dataPath := filepath.Join(d.path, strconv.Itoa(e.ID))
	os.Remove(dataPath)
	return nil
}

func (d *Database) saveIndex() error {
	encoded, _ := json.Marshal(d.index)
	indexPath := filepath.Join(d.path, indexFilename)
	return ioutil.WriteFile(indexPath, encoded, 0755)
}

func createDatabase(path string) (*Database, error) {
	if err := os.Mkdir(path, 0755); err != nil {
		return nil, err
	}
	newIndex := &index{
		IDToEntry:   map[string]DatabaseEntry{},
		ShareIDToID: map[string]int{},
		CurrentID:   0,
	}
	indexData, _ := json.Marshal(newIndex)
	indexFile := filepath.Join(path, indexFilename)
	db := &Database{path: path, index: newIndex}
	return db, ioutil.WriteFile(indexFile, indexData, 0755)
}

func copyAndCountLines(dst io.Writer, src io.Reader) (int, error) {
	buf := make([]byte, copyBufferSize)
	lines := 0
	for {
		n, err := src.Read(buf)
		if err != nil && err != io.EOF {
			return 0, err
		}
		if n != 0 {
			if _, err := dst.Write(buf[:n]); err != nil {
				return 0, err
			}
		}
		for _, ch := range buf[:n] {
			if ch == '\n' {
				lines++
			}
		}
		if err == io.EOF {
			break
		}
	}
	return lines, nil
}

type DatabaseEntry struct {
	ID      int
	ShareID string

	Language  string
	PostDate  time.Time
	PosterIP  string
	LineCount int
}

type index struct {
	IDToEntry   map[string]DatabaseEntry
	ShareIDToID map[string]int
	CurrentID   int
}

func (i *index) entryForID(id int) (entry DatabaseEntry, ok bool) {
	entry, ok = i.IDToEntry[strconv.Itoa(id)]
	return
}

func randomShareID() string {
	key := securecookie.GenerateRandomKey(16)
	return strings.ToLower(hex.EncodeToString(key))
}
