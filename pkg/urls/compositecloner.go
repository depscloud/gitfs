package urls

import (
	"fmt"

	"gopkg.in/src-d/go-billy.v4"
)

// NewCloner encapsulates the functionality for creating the default cloner.
// The default cloner maintains a route table that uses the url info to map to the appropriate cloner.
func NewCloner() Cloner {
	cloners := make(map[VCS]Cloner)
	cloners[GIT] = &gitcloner{}

	return &compositecloner{
		cloners: cloners,
	}
}

var _ Cloner = &compositecloner{}

type compositecloner struct {
	cloners map[VCS]Cloner
}

func (cc *compositecloner) Clone(url *URL, depth int, fs billy.Filesystem) error {
	cloner, ok := cc.cloners[url.VCS]
	if !ok {
		return fmt.Errorf("unsupported vcs: %s", url.VCS)
	}

	return cloner.Clone(url, depth, fs)
}
