package main

import (
	"encoding/json"
	"log"
	"os"
)

type jsonSerializer struct {
	filename string
}

func NewJsonSerializer(filename string) *jsonSerializer {
	return &jsonSerializer{
		filename: filename,
	}
}

func (s *jsonSerializer) Load() (RepoList, error) {
	f, err := os.Open(s.filename)
	if err != nil {
		log.Println("open file failed!", err)
		return nil, err
	}
	defer f.Close()

	var repos RepoList
	err = json.NewDecoder(f).Decode(&repos)
	if err != nil {
		log.Println("decode externals json failed!", err)
		return nil, err
	}

	return repos, nil
}

func (s *jsonSerializer) Save(repos RepoList) error {
	f, err := os.Create(s.filename)
	if err != nil {
		log.Println("open file failed!", err)
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "    ")
	err = enc.Encode(&repos)
	if err != nil {
		log.Println("encode to externals json failed!", err)
		return err
	}

	return nil
}
