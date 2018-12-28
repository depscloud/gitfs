package config_test

import (
	"testing"

	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/stretchr/testify/require"
)

func test_basic(t *testing.T, basic *config.Basic) {
	require.NotNil(t, basic)
	require.Equal(t, "username", basic.Username)

	require.NotNil(t, basic.Password)
	require.Equal(t, "password", basic.Password.Value)
}

func test_oauth(t *testing.T, token *config.OAuthToken) {
	require.NotNil(t, token)
	require.Equal(t, "token", token.Token)

	if token.ApplicationId != nil {
		require.Equal(t, "application_id", token.ApplicationId.Value)
	}
}

func test_oauth2(t *testing.T, token *config.OAuth2Token) {
	require.NotNil(t, token)
	require.Equal(t, "token", token.Token)

	require.NotNil(t, token.TokenType)
	require.Equal(t, "token_type", token.TokenType.Value)

	require.NotNil(t, token.RefreshToken)
	require.Equal(t, "refresh_token", token.RefreshToken.Value)

	require.NotNil(t, token.Expiry)
	require.Equal(t, "expiry", token.Expiry.Value)
}

func test_generic(t *testing.T, config *config.Generic) {
	require.NotNil(t, config)
	require.Equal(t, "base_url", config.BaseUrl)
	require.Equal(t, "path", config.Path)
	require.Equal(t, "per_page_parameter", config.PerPageParameter)
	require.Equal(t, "page_parameter", config.PageParameter)
	require.Equal(t, int32(20), config.PageSize)
	require.Equal(t, "selector", config.Selector)
}

func test_gitlab(t *testing.T, config *config.Gitlab) {
	require.NotNil(t, config)
	require.NotNil(t, config.BaseUrl)

	require.Equal(t, "base_url", config.BaseUrl.Value)
}

func test_github(t *testing.T, config *config.Github) {
	require.NotNil(t, config)
	require.NotNil(t, config.BaseUrl)
	require.Equal(t, "base_url", config.BaseUrl.Value)

	require.NotNil(t, config.UploadUrl)
	require.Equal(t, "upload_url", config.UploadUrl.Value)

	require.NotNil(t, config.Organizations)
	require.Len(t, config.Organizations, 1)
	require.Contains(t, config.Organizations, "org1")

	require.NotNil(t, config.Users)
	require.Len(t, config.Users, 1)
	require.Contains(t, config.Users, "user1")
}

func test_clone_override(t *testing.T, override *config.CloneOverride) {
	require.NotNil(t, override)

	require.NotNil(t, override.RepositoryRoot)
	require.Equal(t, "repository_root", override.RepositoryRoot.Value)

	require.NotNil(t, override.Depth)
	require.Equal(t, int32(0), override.Depth.Value)
}

func test_clone(t *testing.T, clone *config.CloneConfiguration) {
	require.NotNil(t, clone)
	require.NotNil(t, clone.RepositoryRoot)
	require.Equal(t, "repository_root", clone.RepositoryRoot.Value)

	require.NotNil(t, clone.Depth)
	require.Equal(t, int32(1), clone.Depth.Value)

	{
		override, ok := clone.Overrides["regex.*"]
		require.True(t, ok)
		test_clone_override(t, override)
	}

	{
		override, ok := clone.Overrides["string-match"]
		require.True(t, ok)
		test_clone_override(t, override)
	}
}

func test_common(t *testing.T, cfg *config.Configuration) {
	require.Len(t, cfg.Accounts, 8)

	{
		generic := cfg.Accounts[0].GetGeneric()
		test_generic(t, generic)
	}

	{
		generic := cfg.Accounts[1].GetGeneric()
		test_generic(t, generic)
		test_basic(t, generic.Basic)
	}

	{
		gitlab := cfg.Accounts[2].GetGitlab()
		test_gitlab(t, gitlab)
		test_oauth(t, gitlab.Private)
	}

	{
		gitlab := cfg.Accounts[3].GetGitlab()
		test_gitlab(t, gitlab)
		test_oauth(t, gitlab.Oauth)
	}

	{
		github := cfg.Accounts[4].GetGithub()
		test_github(t, github)
	}

	{
		github := cfg.Accounts[5].GetGithub()
		test_github(t, github)
		test_oauth2(t, github.Oauth2)
	}

	{
		bitbucket := cfg.Accounts[6].GetBitbucket()
		require.NotNil(t, bitbucket)
		test_basic(t, bitbucket.Basic)
	}

	{
		bitbucket := cfg.Accounts[7].GetBitbucket()
		require.NotNil(t, bitbucket)
		test_oauth(t, bitbucket.Oauth)
	}

	test_clone(t, cfg.Clone)
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
