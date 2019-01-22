package urls

import (
	"gopkg.in/src-d/go-billy.v4"
)

// Cloner defines how classes should clone repositories.
type Cloner interface {
	Clone(url *URL, depth int, fs billy.Filesystem) error
}
