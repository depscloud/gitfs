package config_test

import (
	"testing"

	"github.com/deps-cloud/gitfs/pkg/config"
	"github.com/stretchr/testify/require"
)

func testCloneOverride(t *testing.T, override *config.CloneOverride) {
	require.NotNil(t, override)

	require.NotNil(t, override.RepositoryRoot)
	require.Equal(t, "repository_root", override.RepositoryRoot.Value)

	require.NotNil(t, override.Depth)
	require.Equal(t, int32(0), override.Depth.Value)
}

func testClone(t *testing.T, clone *config.CloneConfiguration) {
	require.NotNil(t, clone)
	require.NotNil(t, clone.RepositoryRoot)
	require.Equal(t, "repository_root", clone.RepositoryRoot.Value)

	require.NotNil(t, clone.Depth)
	require.Equal(t, int32(1), clone.Depth.Value)

	{
		override, ok := clone.Overrides["regex.*"]
		require.True(t, ok)
		testCloneOverride(t, override)
	}

	{
		override, ok := clone.Overrides["string-match"]
		require.True(t, ok)
		testCloneOverride(t, override)
	}
}

func testCommon(t *testing.T, cfg *config.Configuration) {
	require.Len(t, cfg.Accounts, 9)

	testClone(t, cfg.Clone)
}

func Test_proto(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.prototxt")
	require.NoError(t, err)
	testCommon(t, cfg)
}

func Test_yaml(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.yaml")
	require.NoError(t, err)
	testCommon(t, cfg)
}

func Test_json(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.json")
	require.NoError(t, err)
	testCommon(t, cfg)
}
