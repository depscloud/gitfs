package remotes

import (
	"fmt"
	"github.com/mjpitz/gitfs/pkg/config"
)

func ParseConfig(configuration *config.Configuration) (Remote, error) {
	remotes := make([]Remote, len(configuration.Accounts))
	for i, account := range configuration.Accounts {
		if generic := account.GetGeneric(); generic != nil {
			remotes[i] = NewGenericRemote(generic)
		} else if bitbucket := account.GetBitbucket(); bitbucket != nil {
			return nil, fmt.Errorf("upsupported: bitbucket")
		} else if github := account.GetGithub(); github != nil {
			return nil, fmt.Errorf("upsupported: github")
		} else if gitlab := account.GetGitlab(); gitlab != nil {
			return nil, fmt.Errorf("upsupported: gitlab")
		} else {
			return nil, fmt.Errorf("unrecognized account")
		}
	}
	return NewCompositeRemote(remotes...), nil
}
