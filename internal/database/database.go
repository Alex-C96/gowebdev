package database

import (
	"encoding/json"
	"log"
	"os"
	"sort"
	"sync"
)

const DatabaseString = "database.json"

type DB struct {
	path    string
	mux     *sync.RWMutex
	ChirpId int
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	err := ensureDB()
	if err != nil {
		if os.IsNotExist(err) {
			os.Create(path)
		} else {
			return &DB{}, err
		}
	}
	return &DB{
		DatabaseString,
		&sync.RWMutex{},
		1,
	}, nil
}

func ensureDB() error {
	file, err := os.Open(DatabaseString)

	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	DbStruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirp := Chirp{
		Body: body,
		Id:   db.ChirpId,
	}
	DbStruct.Chirps[db.ChirpId] = chirp
	db.writeDB(DbStruct)
	db.ChirpId++

	return chirp, err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{make(map[int]Chirp)}, err
	}
	log.Println("Read from file successfully")
	DbStruct := DBStructure{make(map[int]Chirp)}
	if len(data) > 1 {
		err = json.Unmarshal(data, &DbStruct)
		if err != nil {
			return DBStructure{make(map[int]Chirp)}, err
		}
		log.Println("Unmarshaled database.json file")
	}
	return DbStruct, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	writeData, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	log.Println("Marshaled new file data")
	err = os.WriteFile(db.path, writeData, 0666)
	if err != nil {
		return err
	}
	log.Println("Wrote to the new file")
	return nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	DbStruct, err := db.loadDB()
	if err != nil {
		return make([]Chirp, 0), err
	}
	chirps := make([]Chirp, 0)
	for _, chirp := range DbStruct.Chirps {
		chirps = append(chirps, chirp)
	}
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })
	return chirps, nil
}
