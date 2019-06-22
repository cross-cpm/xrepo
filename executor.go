package main

func NewExecutor(url string, info *Info) *gitExecutor {
	var workdir string
	for src, dsts := range info.Targets {
		if src == "./" {
			workdir = dsts[0]
		}
	}

	// TODO: check vcs type then return executor of git or svn

	return newGitExecutor(workdir, url, info)
}
