package repo

type externals struct {
}

func NewExternals() *externals {
	return &externals{}
}

func (e *externals) Load(filename string) error {
	return nil
}

func (e *externals) Save(filename string) error {
	return nil
}

func (e *externals) ForEach() error {
	return nil
}
