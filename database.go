package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
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

func (d *Database) LatestEntries(count int) []DatabaseEntry {
	res := make([]DatabaseEntry, 0, count)
	d.lock.RLock()
	defer d.lock.RUnlock()
	for i := d.index.CurrentId; i >= 0 && len(res) < count; i-- {
		if entry, ok := d.index.IDToEntry[i]; ok {
			res = append(res, entry)
		}
	}
	return res
}

func (d *Database) EntriesInRange(start, end int) []DatabaseEntry {
	res := make([]DatabaseEntry, 0, start-end+1)
	d.lock.RLock()
	defer d.lock.RUnlock()
	for i := end; i >= start && i >= 0; i-- {
		if entry, ok := d.index.IDToEntry[i]; ok {
			res = append(res, entry)
		}
	}
	return res
}

func (d *Database) OpenEntry(shareID string) (e DatabaseEntry, r io.Reader, err error) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	e, ok := d.index.ShareIDToEntry[shareID]
	if !ok {
		err = errors.New("unknown share ID: " + shareID)
		return
	}
	dataPath := filepath.Join(d.path, strconv.Itoa(e.ID))
	r, err = os.Open(dataPath)
	return
}

func (d *Database) CreateEntry(info DatabaseEntry) (entry DatabaseEntry, err error) {
	// TODO: upload the file, allocate an ID, and add it to the DB index.
	err = errors.New("not yet implemented")
	return
}

func createDatabase(path string) (*Database, error) {
	if err := os.Mkdir(path, 0755); err != nil {
		return nil, err
	}
	newIndex := &index{
		IDToEntry:      map[int]DatabaseEntry{},
		ShareIDToEntry: map[string]DatabaseEntry{},
		CurrentId:      0,
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
	IDToEntry      map[int]DatabaseEntry
	ShareIDToEntry map[string]DatabaseEntry
	CurrentId      int
}
