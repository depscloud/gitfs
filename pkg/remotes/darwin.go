package remotes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/pkg/errors"
	rlog "github.com/sirupsen/logrus"
)

// NewDarwinRemote constructs a Remote pointing at the base url, using the semantics of Indeed's internal
// darwin api. Eventually, this will support systems like gitlab, github, and bitbucket.
func NewDarwinRemote(baseUrl string) *DarwinRemote {
	return &DarwinRemote{
		baseUrl: baseUrl,
	}
}

var _ Remote = &DarwinRemote{}

// DarwinRemote implements the Remote interface. The implementation calls out to our internal darwin api
// which uses a similar api definition to Gitlab. Because the service is written in Java, it changes the
// formatting for datetime, making it incompatible with the go-gitlab client.
type DarwinRemote struct {
	baseUrl string
}

// fetch defines a common way to perform HTTP GET requests. Since darwin is readonly, we do not need to
// support other operations. This function encapsulates the URL construction, and the construction of the
// decoder.
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

// getProjects calls the /api/projects endpoint for the requested page. Additional wrapping logic, as seen
// in the ListRepositories, is required for paginating results.
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
