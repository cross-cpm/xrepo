package repo

func NewExecutor(url string, info *Info) *gitExecutor {
	var workdir string
	for src, dsts := range info.Targets {
		if src == "./" {
			workdir = dsts[0]
		}
	}

	return newGitExecutor(workdir, url, info)
}
