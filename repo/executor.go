package repo

func NewExecutor(url string, workdir string, info *Info) *gitExecutor {
	return newGitExecutor(url, workdir, info)
}
