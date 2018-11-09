package remotes

import (
	"fmt"
	"indeed/delivery/darwin/go/libdarwin"
	"indeed/gophers/rlog"
)

func NewDarwinRemote() *DarwinRemote {
	return &DarwinRemote{
		client: libdarwin.New("https://darwin.sandbox.indeed.net"),
	}
}

type DarwinRemote struct {
	client *libdarwin.DarwinClient
}

func (darwin *DarwinRemote) ListRepositories() ([]string, error) {
	result := make([]string, 0)

	for i := 0; true; i += 1 {
		rlog.Infof("Fetching page %d", i+1)
		page, err := darwin.client.GetProjects(i, 100)

		if err != nil {
			return nil, fmt.Errorf("failed to fetch projects from darwin: %v", err)
		}

		pageUrls := make([]string, len(page))
		for j, project := range page {
			pageUrls[j] = project["ssh_url_to_repo"].(string)
		}

		result = append(result, pageUrls...)

		if len(page) < 100 {
			break
		}
	}

	return result, nil
}
