package testdata

import (
	"encoding/json"
	"errors"
	"goexamples/poem-stream/pb"
	"log"
	"os"
)

type DB map[string]*pb.Poem

func (db DB) Load(file string) error {
	if file == "" {
		return errors.New("file is empty")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	p := []*pb.Poem{}
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	for _, poem := range p {
		db[poem.Title] = poem
	}
	return nil
}

func (db DB) GetPoem(title string) (*pb.Poem, error) {
	if v, ok := db[title]; ok {
		return v, nil
	}
	return nil, errors.New("poem not found")
}

func (db DB) GetPoemCollection() []*pb.Poem {
	poems := make([]*pb.Poem, len(db))
	i := 0
	for _, v := range db {
		poems[i] = v
		i++
	}
	return poems
}

func (db DB) SetPoem(title string, poem *pb.Poem) {
	db[title] = poem
}

func NewDB(file string) DB {
	db := make(DB)
	if err := db.Load(file); err != nil {
		log.Fatalf("failed to load db file[%s]: %v", file, err)
	}
	return db
}
