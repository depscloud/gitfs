package config_test

import (
	"testing"

	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/stretchr/testify/require"
)

func test_common(t *testing.T, cfg *config.Configuration) {
	require.Equal(t, 8, len(cfg.Accounts))

	require.NotNil(t, cfg.Accounts[0].GetGeneric())
	require.NotNil(t, cfg.Accounts[1].GetGeneric())

	require.NotNil(t, cfg.Accounts[2].GetGitlab())
	require.NotNil(t, cfg.Accounts[3].GetGitlab())

	require.NotNil(t, cfg.Accounts[4].GetGithub())
	require.NotNil(t, cfg.Accounts[5].GetGithub())

	require.NotNil(t, cfg.Accounts[6].GetBitbucket())
	require.NotNil(t, cfg.Accounts[7].GetBitbucket())
}

func Test_proto(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.prototxt")
	require.NoError(t, err)
	test_common(t, cfg)
}

func Test_yaml(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.yaml")
	require.NoError(t, err)
	test_common(t, cfg)
}

func Test_json(t *testing.T) {
	cfg, err := config.Load("../../hack/config/full.json")
	require.NoError(t, err)
	test_common(t, cfg)
}
