package urls_test

import (
	"github.com/mjpitz/gitfs/pkg/urls"
	"github.com/stretchr/testify/require"
	"testing"
)


func parse(t *testing.T, url string, vcs urls.VCS) error {
	parsedUrl, err := urls.ParseUrl(url)
	if err != nil {
		return err
	}

	require.Equal(t, vcs, parsedUrl.VCS)
	require.Equal(t, url, parsedUrl.String())

	return nil
}

func Test_ParseUrl(t *testing.T) {
	require.NoError(t, parse(t, "git@github.com:mjpitz/gitfs.git", urls.GIT))
	require.NoError(t, parse(t, "https://github.com/mjpitz/gitfs.git", urls.GIT))
	require.NoError(t, parse(t, "http://github.com/mjpitz/gitfs.git", urls.GIT))

	require.NoError(t, parse(t, "svn://github.com/mjpitz/gitfs", urls.SVN))
	require.NoError(t, parse(t, "svn+ssh://github.com/mjpitz/gitfs", urls.SVN))

	require.NoError(t, parse(t, "mjpitz/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "mjpitz/gitfs#master", urls.MERCURIAL))
	require.NoError(t, parse(t, "file://mjpitz/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "file://mjpitz/gitfs#master", urls.MERCURIAL))
	require.NoError(t, parse(t, "https://github.com/mjpitz/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "https://github.com/mjpitz/gitfs#master", urls.MERCURIAL))
	require.NoError(t, parse(t, "http://github.com/mjpitz/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "http://github.com/mjpitz/gitfs#master", urls.MERCURIAL))
	require.NoError(t, parse(t, "ssh://github.com/mjpitz/gitfs", urls.MERCURIAL))
	require.NoError(t, parse(t, "ssh://github.com/mjpitz/gitfs#master", urls.MERCURIAL))
}