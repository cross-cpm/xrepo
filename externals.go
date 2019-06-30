package main

type externals struct {
	serializer Serializer
	repos      RepoList
}

func NewExternals(filename string) *externals {
	return &externals{
		serializer: NewSerializer(filename),
	}
}

func (e *externals) Load() error {
	repos, err := e.serializer.Load()
	if err != nil {
		return err
	}

	e.repos = repos
	return nil
}

func (e *externals) Save() error {
	return e.serializer.Save(e.repos)
}

func (e *externals) Count() int {
	return len(e.repos)
}

func (e *externals) Foreach(fn func(url string, info *RepoInfo)) error {
	for url, info := range e.repos {
		fn(url, info)
	}

	return nil
}
