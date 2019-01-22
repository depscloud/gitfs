package urls

import (
	"fmt"
	"net/url"
	"strings"
)

var _ fmt.Stringer = &URL{}

// URL defines the address where a repository is located.
type URL struct {
	VCS
	URL *url.URL
}

func (u *URL) String() string {
	if u.URL.Scheme == "git" {
		// git ssh
		path := strings.TrimLeft(u.URL.Path, "/")
		return fmt.Sprintf("git@%s:%s", u.URL.Host, path)
	}

	return u.URL.String()
}

// ParseURL accepts a url string and parses it into a well structured url.
//
// git:
// git@<<HOST>>:<<PATH>>.git
// https://<<HOST>><<PATH>>.git
// http://<<HOST>><<PATH>>.git
//
// svn:
// svn://<<HOST>><<PATH>>
// svn+ssh://<<HOST>>/<PATH>
//
// hg:
// local/filesystem/path[#revision]
// file://local/filesystem/path[#revision]
// http://[user[:pass]@]host[:port]/[path][#revision]
// https://[user[:pass]@]host[:port]/[path][#revision]
// ssh://[user@]host[:port]/[path][#revision]
func ParseURL(urlString string) (*URL, error) {
	gitSSH := strings.HasPrefix(urlString, "git@")
	gitRepo := strings.HasSuffix(urlString, ".git")

	if gitSSH && gitRepo {
		idx := strings.LastIndex(urlString, ":")

		urlString = strings.Replace(urlString, ":", "/", idx)
		urlString = "git://" + urlString
	}

	uri, err := url.Parse(urlString)
	if err != nil {
		return nil, fmt.Errorf("invalid urlString: %s", uri)
	}

	vcs := MERCURIAL
	if uri.Scheme == "git" || strings.HasSuffix(uri.Path, ".git") {
		vcs = GIT
	} else if uri.Scheme == "svn+ssh" || uri.Scheme == "svn" {
		vcs = SVN
	}

	return &URL{
		VCS: vcs,
		URL: uri,
	}, err
}
