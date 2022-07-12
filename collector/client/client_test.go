package client

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"exercism-mentoring-request-notifier/config"
	"exercism-mentoring-request-notifier/request"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	expectedToken string
	responseBody  []io.ReadCloser
	err           error
	statusCode    int
	expectedURL   string
	requestCount  int
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	m.requestCount++
	if m.err != nil {
		return nil, m.err
	}

	if fmt.Sprintf("%s%d", m.expectedURL, m.requestCount) != req.URL.String() {
		return nil, fmt.Errorf("invalid request-url, expected: %s, got: %s", fmt.Sprintf("%s%d", m.expectedURL, m.requestCount), req.URL.String())
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader != fmt.Sprintf("Bearer %s", m.expectedToken) {
		return nil, fmt.Errorf("invalid auth token presented")
	}

	return &http.Response{StatusCode: m.statusCode, Body: m.responseBody[m.requestCount-1]}, nil
}

func TestGetMentoringRequestsForAllTracks(t *testing.T) {
	for _, testCase := range testCasesGetMentoringRequestsForAllTracks {
		mentoringRequestsPerTrack, err := testCase.getClient(testCase.result).GetMentoringRequestsForAllTracks(map[string]config.TrackConfig{"go": {}})
		if err != nil {
			return
		}
		assertError(t, err, testCase.expectError)

		var expectedResult = map[string][]request.MentoringRequest{"go": testCase.result.MentoringRequests}
		assert.Equal(t, expectedResult, mentoringRequestsPerTrack)
	}
}

func TestGetAllMentoringRequests(t *testing.T) {
	for _, testCase := range testCasesGetAllMentoringRequests {
		t.Run(testCase.description, func(t *testing.T) {
			mentoringRequests, err := testCase.getClient(testCase.result).getAllMentoringRequests("go")
			assertError(t, err, testCase.expectError)
			assert.Equal(t, testCase.result.MentoringRequests, mentoringRequests)
		})
	}
}

func TestGetRequest(t *testing.T) {
	for _, testCase := range testCasesGetRequest {
		t.Run(testCase.description, func(t *testing.T) {
			responseBody, err := testCase.client.getRequest(testCase.url)
			assertError(t, err, testCase.expectError)
			assert.Equal(t, testCase.expected, responseBody)
		})
	}
}

func assertError(t *testing.T, err error, expectError bool) {
	if expectError {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
	}
}
