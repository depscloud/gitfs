package urls

import (
	"gopkg.in/src-d/go-billy.v4"
)

type Cloner interface {
	Clone(url *URL, depth int, fs billy.Filesystem) error
}
