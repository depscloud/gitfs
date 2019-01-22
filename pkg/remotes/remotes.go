package remotes

import (
	"fmt"

	"github.com/mjpitz/gitfs/pkg/config"
)

// ParseConfig is used to parse the account configuration and construct the
// necessary remote endpoint based on the configuration object.
func ParseConfig(configuration *config.Configuration) (Remote, error) {
	remotes := make([]Remote, len(configuration.Accounts))

	for i, account := range configuration.Accounts {
		var remote Remote
		var err error

		if generic := account.GetGeneric(); generic != nil {
			remote = NewGenericRemote(generic)
		} else if bitbucket := account.GetBitbucket(); bitbucket != nil {
			err = fmt.Errorf("upsupported: bitbucket")
		} else if github := account.GetGithub(); github != nil {
			remote, err = NewGithubRemote(github)
		} else if gitlab := account.GetGitlab(); gitlab != nil {
			err = fmt.Errorf("upsupported: gitlab")
		} else {
			err = fmt.Errorf("unrecognized account")
		}

		if err != nil {
			return nil, err
		}

		remotes[i] = remote
	}

	return NewCompositeRemote(remotes...), nil
}
