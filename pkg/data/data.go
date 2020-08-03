package data

import (
	"io/ioutil"
	"log"
	"os"
)

type DataType interface {
	GetID()
	SetID()
	SaveToStore()
}

type Store struct {
	Name string
	Path string
}

func New(storeName string) *Store {
	return &Store{storeName, ""}
}

func (s *Store) GetFilePath() string {
	return s.Path
}

func (s *Store) SetFilePath(pathValue string) {
	s.Path = pathValue
}

func (s *Store) GetContent() []byte {
	jsonFile, err := os.Open(s.GetFilePath())
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()
	fileContent, _ := ioutil.ReadAll(jsonFile)
	return fileContent
}

func (s *Store) Write(data []byte) {
	err := ioutil.WriteFile(s.GetFilePath(), data, 660)
	if err != nil {
		log.Fatal(err)
	}
}
