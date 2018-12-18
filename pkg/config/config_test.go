package config_test

import (
	"bytes"
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/stretchr/testify/require"

	"github.com/ghodss/yaml"
)

func test_common(t *testing.T, cfg *config.Configuration) {
	require.Equal(t, 8, len(cfg.Accounts))

	require.NotNil(t, cfg.Accounts[0].GetDarwin())
	require.NotNil(t, cfg.Accounts[1].GetDarwin())

	require.NotNil(t, cfg.Accounts[2].GetGitlab())
	require.NotNil(t, cfg.Accounts[3].GetGitlab())

	require.NotNil(t, cfg.Accounts[4].GetGithub())
	require.NotNil(t, cfg.Accounts[5].GetGithub())

	require.NotNil(t, cfg.Accounts[6].GetBitbucket())
	require.NotNil(t, cfg.Accounts[7].GetBitbucket())
}

func Test_proto(t *testing.T) {
	pto, err := ioutil.ReadFile("../../hack/config/full.txt")
	require.NoError(t, err)

	cfg := &config.Configuration{}

	err = proto.UnmarshalText(string(pto), cfg)
	require.NoError(t, err)

	test_common(t, cfg)
}

func Test_yaml(t *testing.T) {
	yml, err := ioutil.ReadFile("../../hack/config/full.yaml")
	require.NoError(t, err)

	jsn, err := yaml.YAMLToJSON(yml)
	require.NoError(t, err)

	cfg := &config.Configuration{}

	err = jsonpb.Unmarshal(bytes.NewReader(jsn), cfg)
	require.NoError(t, err)

	test_common(t, cfg)
}

func Test_json(t *testing.T) {
	jsn, err := ioutil.ReadFile("../../hack/config/full.json")
	require.NoError(t, err)

	cfg := &config.Configuration{}

	err = jsonpb.Unmarshal(bytes.NewReader(jsn), cfg)
	require.NoError(t, err)

	test_common(t, cfg)
}
