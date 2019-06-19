package repo

import (
	"encoding/json"
	"log"
	"os"
)

type Info struct {
	Branch  string              `json:"branch"`
	Ref     string              `json:"ref"`
	Targets map[string][]string `json:"targets"`
}

type externals struct {
	filename string
	infos    map[string]*Info
	workdirs map[string]string
}

func NewExternals(filename string) *externals {
	return &externals{
		filename: filename,
		workdirs: make(map[string]string),
	}
}

func (e *externals) Load() error {
	f, err := os.Open(e.filename)
	if err != nil {
		log.Println("open file failed!", err)
		return err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&e.infos)
	if err != nil {
		log.Println("decode externals json failed!", err)
		return err
	}

	for url, info := range e.infos {
		for src, dsts := range info.Targets {
			if src == "./" {
				e.workdirs[url] = dsts[0]
			}
		}
	}

	// log.Println("debug externals:", e.infos)
	// log.Println("debug workdirs:", e.workdirs)
	return nil
}

func (e *externals) Save() error {
	return nil
}

func (e *externals) Foreach(fn func(url string, workdir string, info *Info)) error {
	for url, info := range e.infos {
		fn(url, e.workdirs[url], info)
	}

	return nil
}
