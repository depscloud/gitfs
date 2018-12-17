package remotes

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	rlog "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
)

func NewDarwinRemote(baseUrl string) *DarwinRemote {
	return &DarwinRemote{
		baseUrl: baseUrl,
	}
}

type DarwinRemote struct {
	baseUrl string
}

func (darwin *DarwinRemote) fetch(path string, qs url.Values) (*json.Decoder, error) {
	fullUrl := fmt.Sprintf("%s%s", darwin.baseUrl, path)

	if qs != nil && len(qs) > 0 {
		fullUrl = fmt.Sprintf("%s?%s", fullUrl, qs.Encode())
	}

	resp, err := http.Get(fullUrl)
	if err != nil {
		return nil, err
	}

	return json.NewDecoder(resp.Body), nil
}


func (darwin *DarwinRemote) getProjects(page, pageSize int) ([]map[string]interface{}, error) {
	qs := url.Values{}
	qs.Set("page", strconv.Itoa(page))
	qs.Set("per_page", strconv.Itoa(pageSize))

	body, err := darwin.fetch("/api/projects", qs)
	if err != nil {
		return nil, err
	}

	//var payload []*interface{}
	payload := make([]map[string]interface{}, 0)
	err = body.Decode(&payload)
	return payload, err

}

func (darwin *DarwinRemote) ListRepositories() ([]string, error) {
	result := make([]string, 0)

	for i := 0; true; i += 1 {
		rlog.Infof("Fetching page %d", i+1)
		page, err := darwin.getProjects(i, 100)

		if err != nil {
			return nil, errors.Wrap(err, "failed to fetch projects from darwin")
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
