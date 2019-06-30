package main

import (
	"log"
	"path/filepath"
)

type RepoInfo struct {
	Cvs     string              `json:"cvs"`
	Branch  string              `json:"branch"`
	Ref     string              `json:"ref"`
	Targets map[string][]string `json:"targets"`
}

type RepoList map[string]*RepoInfo

type Serializer interface {
	Load() (RepoList, error)
	Save(RepoList) error
}

func NewSerializer(filename string) Serializer {
	extname := filepath.Ext(filename)
	switch extname {
	case ".yaml":
		return NewYamlSerializer(filename)
	case ".json":
		return NewJsonSerializer(filename)
	default:
		log.Fatal("unsupport ext name:", extname)
	}

	return nil
}
