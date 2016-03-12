package main

import (
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

// EntriesBefore returns a certain number of entries whose IDs descent from a starting ID.
func (d *Database) EntriesBefore(startId, count int) []DatabaseEntry {
	res := make([]DatabaseEntry, 0, count)
	d.lock.RLock()
	defer d.lock.RUnlock()
	if startId == -1 {
		startId = d.index.CurrentId - 1
	}
	for i := startId; i >= 0 && len(res) < count; i-- {
		if entry, ok := d.index.IDToEntry[i]; ok {
			res = append(res, entry)
		}
	}
	return res
}

func (d *Database) OpenEntry(shareID string) (e DatabaseEntry, r io.Reader, err error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	id, ok := d.index.ShareIDToID[shareID]
	if !ok {
		err = errors.New("unknown share ID: " + shareID)
		return
	}
	e, ok = d.index.IDToEntry[id]
	if !ok {
		err = errors.New("unknown ID: " + strconv.Itoa(id))
		return
	}
	dataPath := filepath.Join(d.path, strconv.Itoa(e.ID))
	r, err = os.Open(dataPath)
	return
}

func (d *Database) CreateEntry(info DatabaseEntry,
	body io.Reader) (entry DatabaseEntry, err error) {
	tempFile, err := ioutil.TempFile("", "t3xt")
	if err != nil {
		return
	}
	_, err = io.Copy(tempFile, body)
	tempFile.Close()
	if err != nil {
		os.Remove(tempFile.Name())
		return
	}

	d.lock.Lock()
	defer d.lock.Unlock()

	entry = info
	entry.ShareID = randomShareID()
	entry.ID = d.index.CurrentId
	d.index.CurrentId++

	dataPath := filepath.Join(d.path, strconv.Itoa(entry.ID))
	err = os.Rename(tempFile.Name(), dataPath)
	if err != nil {
		os.Remove(tempFile.Name())
		return
	}
	d.index.IDToEntry[entry.ID] = entry
	d.index.ShareIDToID[entry.ShareID] = entry.ID
	err = d.saveIndex()
	if err != nil {
		delete(d.index.IDToEntry, entry.ID)
		delete(d.index.ShareIDToID, entry.ShareID)
		os.Remove(dataPath)
	}
	return
}

func (d *Database) DeleteEntry(e DatabaseEntry) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	oldEntry := d.index.IDToEntry[e.ID]
	delete(d.index.IDToEntry, e.ID)
	delete(d.index.ShareIDToID, e.ShareID)
	if err := d.saveIndex(); err != nil {
		d.index.IDToEntry[e.ID] = oldEntry
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
		IDToEntry:   map[int]DatabaseEntry{},
		ShareIDToID: map[string]int{},
		CurrentId:   0,
	}
	indexData, _ := json.Marshal(newIndex)
	indexFile := filepath.Join(path, indexFilename)
	db := &Database{path: path, index: newIndex}
	return db, ioutil.WriteFile(indexFile, indexData, 0755)
}

type DatabaseEntry struct {
	ID      int
	ShareID string

	Language string
	PostDate time.Time
	PosterIP string
}

type index struct {
	IDToEntry   map[int]DatabaseEntry
	ShareIDToID map[string]int
	CurrentId   int
}

func randomShareID() string {
	key := securecookie.GenerateRandomKey(16)
	return strings.ToLower(hex.EncodeToString(key))
}
