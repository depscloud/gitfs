package clone_test

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/mjpitz/gitfs/pkg/clone"
	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_resolve(t *testing.T) {
	defaultRoot := "/default_root"
	customRoot := "/custom_root"

	cfg := &config.CloneConfiguration{
		RepositoryRoot: &wrappers.StringValue{
			Value: defaultRoot,
		},
		Overrides: map[string]*config.CloneOverride{
			"github.com:mjpitz/.*": {
				RepositoryRoot: &wrappers.StringValue{
					Value: customRoot,
				},
			},
			"git@github.com:test/depth-only.git": {
				Depth: &wrappers.Int32Value{
					Value: 0,
				},
			},
			"git@github.com:test/both.git": {
				RepositoryRoot: &wrappers.StringValue{
					Value: customRoot,
				},
				Depth: &wrappers.Int32Value{
					Value: 0,
				},
			},
		},
	}

	cloner := clone.NewCloner(cfg)

	data := [][]interface{}{
		{"git@github.com:alexellis/k8s-on-raspbian.git", defaultRoot, int32(1)},
		{"git@github.com:mjpitz/gitfs.git", customRoot, int32(1)},
		{"git@github.com:test/depth-only.git", defaultRoot, int32(0)},
		{"git@github.com:test/both.git", customRoot, int32(0)},
	}

	for _, entry := range data {
		url := entry[0].(string)
		root, depth, _ := cloner.Resolve(url)
		require.Equal(t, entry[1], root, url)
		require.Equal(t, entry[2], depth, url)
	}
}
