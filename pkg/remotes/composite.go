package remotes

import "github.com/sirupsen/logrus"

// NewCompositeRemote wraps the supplied remotes in a composite wrapper
// which logs errors and continues processing remote endpoints.
func NewCompositeRemote(remotes ...Remote) Remote {
	return &compositeRemote{
		remotes: remotes,
	}
}

var _ Remote = &compositeRemote{}

type compositeRemote struct {
	remotes []Remote
}

func (r *compositeRemote) ListRepositories() ([]string, error) {
	repositories := make([]string, 0)
	for _, remote := range r.remotes {
		repos, err := remote.ListRepositories()

		if err != nil {
			logrus.Errorf("[remotes.composite] failed to list repositories from remote: %v", err)
		} else {
			repositories = append(repositories, repos...)
		}
	}
	return repositories, nil
}
