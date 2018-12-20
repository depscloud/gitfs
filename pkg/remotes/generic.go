package remotes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mjpitz/gitfs/pkg/config"
	"github.com/pkg/errors"
	"github.com/savaki/jq"
	"github.com/sirupsen/logrus"
)

func NewGenericRemote(config *config.Generic) *GenericRemote {
	return &GenericRemote{
		config: config,
	}
}

var _ Remote = &GenericRemote{}

type GenericRemote struct {
	config *config.Generic
}

func parseJQResult(jqResult string) []string {
	lines := strings.Split(jqResult, "\n")

	parsed := make([]string, len(lines))
	for i, line := range lines {
		cleaned := strings.TrimPrefix(line, "\"")
		cleaned = strings.TrimSuffix(cleaned, "\"")
		parsed[i] = cleaned
	}
	return parsed
}

func (r *GenericRemote) ListRepositories() ([]string, error) {
	op, err := jq.Parse(r.config.Selector)
	if err != nil {
		return nil, err
	}

	repositories := make([]string, 0)
	for page := 1; true; page++ {
		fullUrl := fmt.Sprintf(
			"%s%s?%s=%d&%s=%d",
			r.config.BaseUrl,
			r.config.Path,
			r.config.PageParameter,
			page,
			r.config.PerPageParameter,
			r.config.PageSize,
		)

		resp, err := http.Get(fullUrl)
		if err != nil {
			return nil, errors.Wrap(err,
				fmt.Sprintf("failed to get url: %s", fullUrl))
		}

		if resp.StatusCode == http.StatusNotFound {
			logrus.Infof("encountered a 404. assuming end of data")
			break
		}

		body, err := ioutil.ReadAll(resp.Body)

		selected, err := op.Apply(body)
		if err != nil {
			return nil, errors.Wrap(err,
				fmt.Sprintf("failed to extract content using selector: %s", r.config.Selector))
		}

		currentPage := parseJQResult(string(selected))
		repositories = append(repositories, currentPage...)

		if int32(len(currentPage)) < r.config.PageSize {
			logrus.Infof("encountered an incomplete page. assuming end of data")
			break
		}
	}

	return repositories, nil
}
