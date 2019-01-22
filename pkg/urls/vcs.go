package urls

// VCS defines the version control system used to manage the repo.
type VCS = string

const (
	// GIT defines the constant used to represent the `git` version control system.
	GIT VCS = "git"
	// SVN defines the constant used to represent the `svn` version control system.
	SVN VCS = "svn"
	// MERCURIAL defines the constant used to represent the `hg` version control system.
	MERCURIAL VCS = "hg"
)
