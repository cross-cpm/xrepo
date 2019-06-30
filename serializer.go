package main

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
	return NewJsonSerializer(filename)
}
