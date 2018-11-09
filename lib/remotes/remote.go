package remotes

type Remote interface {
	ListRepositories() ([]string, error)
}
