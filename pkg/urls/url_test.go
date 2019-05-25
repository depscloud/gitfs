package urls_test

import (
	"testing"

	"github.com/deps-cloud/gitfs/pkg/urls"
	"github.com/stretchr/testify/require"
)

func parse(t *testing.T, url string, vcs urls.VCS) error {
	parsedUrl, err := urls.ParseURL(url)
	if err != nil {
		return err
	}

	require.Equal(t, vcs, parsedUrl.VCS)
	require.Equal(t, url, parsedUrl.String())

	return nil
}

func Test_ParseUrl(t *testing.T) {
	require.NoError(t, parse(t, "git@github.com:deps-cloud/gitfs.git", urls.GIT))
	require.NoError(t, parse(t, "https://github.com/deps-cloud/gitfs.git", urls.GIT))
	require.NoError(t, parse(t, "http://github.com/deps-cloud/gitfs.git", urls.GIT))

	require.NoError(t, parse(t, "svn://github.com/deps-cloud/gitfs", urls.SVN))
	require.NoError(t, parse(t, "svn+ssh://github.com/deps-cloud/gitfs", urls.SVN))

	require.NoError(t, parse(t, "deps-cloud/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "deps-cloud/gitfs#master", urls.MERCURIAL))
	require.NoError(t, parse(t, "file://deps-cloud/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "file://deps-cloud/gitfs#master", urls.MERCURIAL))
	require.NoError(t, parse(t, "https://github.com/deps-cloud/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "https://github.com/deps-cloud/gitfs#master", urls.MERCURIAL))
	require.NoError(t, parse(t, "http://github.com/deps-cloud/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "http://github.com/deps-cloud/gitfs#master", urls.MERCURIAL))
	require.NoError(t, parse(t, "ssh://github.com/deps-cloud/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "ssh://github.com/deps-cloud/gitfs#master", urls.MERCURIAL))
}
