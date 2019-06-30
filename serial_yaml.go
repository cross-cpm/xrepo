package main

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type yamlSerializer struct {
	filename string
}

func NewYamlSerializer(filename string) *yamlSerializer {
	return &yamlSerializer{
		filename: filename,
	}
}

func (s *yamlSerializer) Load() (RepoList, error) {
	f, err := os.Open(s.filename)
	if err != nil {
		log.Println("open file failed!", err)
		return nil, err
	}
	defer f.Close()

	var repos RepoList
	err = yaml.NewDecoder(f).Decode(&repos)
	if err != nil {
		log.Println("decode externals yaml failed!", err)
		return nil, err
	}

	return repos, nil
}

func (s *yamlSerializer) Save(repos RepoList) error {
	f, err := os.Create(s.filename)
	if err != nil {
		log.Println("open file failed!", err)
		return err
	}
	defer f.Close()

	enc := yaml.NewEncoder(f)
	err = enc.Encode(&repos)
	if err != nil {
		log.Println("encode to externals yaml failed!", err)
		return err
	}

	return nil
}
