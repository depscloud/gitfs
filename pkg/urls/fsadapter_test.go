package urls_test

import (
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/deps-cloud/gitfs/pkg/config"
	"github.com/deps-cloud/gitfs/pkg/urls"
	"github.com/stretchr/testify/require"
)

func Test_resolve(t *testing.T) {
	defaultRoot := "/default_root"
	customRoot := "/custom_root"

	cfg := &config.CloneConfiguration{
		RepositoryRoot: &wrappers.StringValue{
			Value: defaultRoot,
		},
		Overrides: map[string]*config.CloneOverride{
			"github.com:deps-cloud/.*": {
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

	cloner := urls.NewFileSystemAdapter(cfg)

	data := [][]interface{}{
		{
			"git@github.com:alexellis/k8s-on-raspbian.git",
			"4FDR9G7H91H539H6D2S1H3G1NI5R7DEE2P51MHGHJC2N5KA8CR4G",
			defaultRoot,
			int32(1),
		},
		{
			"git@github.com:deps-cloud/gitfs.git",
			"67U80TE6S1MFPESL7M5JMNP0R3JAK5BS9ITEA0M405M59AURSV70",
			customRoot,
			int32(1),
		},
		{
			"git@github.com:test/depth-only.git",
			"LJS784D56MT1U26K2TJLV6ROKNQ74FFV4GH77D2PI7FAJ7MDR90G",
			defaultRoot,
			int32(0),
		},
		{
			"git@github.com:test/both.git",
			"FQ017M6QS1UBQCGD11L9S2Q47NGPCG010MOGCKPE8O6A890UNOHG",
			customRoot,
			int32(0),
		},
	}

	for _, entry := range data {
		url := entry[0].(string)
		root, bucket, depth, _ := cloner.Resolve(url)

		require.Equal(t, entry[1], bucket, url)
		require.Equal(t, entry[2], root, url)
		require.Equal(t, entry[3], depth, url)
	}
}
