package client

import (
	"encoding/json"
	"exercism-mentoring-request-notifier/request"
	"fmt"
	"io"
	"net/http"
)

const (
	exercismAPIBasePath      = "https://exercism.org/api/v2"
	getMentoringRequestsPath = "/mentoring/requests?track_slug=%s&order=recent&page=%d"
)

type ExercismHTTPClient struct {
	Client HTTPClient
	Token  string
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *ExercismHTTPClient) GetAllMentoringRequests(trackSlug string) ([]request.MentoringRequest, error) {
	var mentoringRequest []request.MentoringRequest
	for page := 1; true; page++ {
		requestURL := fmt.Sprintf("%s%s", exercismAPIBasePath, fmt.Sprintf(getMentoringRequestsPath, trackSlug, page))
		body, err := c.getRequest(requestURL)
		if err != nil {
			return nil, err
		}
		var requests = &request.MentoringRequestsResults{}
		err = json.Unmarshal(body, requests)
		if err != nil {
			return nil, err
		}

		mentoringRequest = append(mentoringRequest, requests.MentoringRequests...)
		if page >= requests.Meta.TotalPages {
			break
		}
	}
	return mentoringRequest, nil
}

func (c *ExercismHTTPClient) getRequest(requestURL string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http-request failed, status-code: %d, response: %s", resp.StatusCode, body)
	}
	return body, nil
}
