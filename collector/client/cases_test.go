package client

import (
	"encoding/json"
	"errors"
	"exercism-mentoring-request-notifier/request"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var testCasesGetAllMentoringRequests = []struct {
	description string
	result      request.MentoringRequestsResults
	expectError bool
	getClient   func(request.MentoringRequestsResults) *ExercismHTTPClient
}{
	{
		description: "valid response with two requests",
		result:      request.MentoringRequestsResults{MentoringRequests: []request.MentoringRequest{{UUID: "e01333ce426b474981391f0566b7a78d", TrackTitle: "C++", ExerciseIconURL: "https://dg8krxphbh767.cloudfront.net/exercises/nth-prime.svg", ExerciseTitle: "Nth Prime", StudentHandle: "Student", StudentAvatarURL: "https://dg8krxphbh767.cloudfront.net/placeholders/user-avatar.svg", UpdatedAt: time.Now().UTC(), HaveMentoredPreviously: false, IsFavorited: false, Status: nil, URL: "https://exercism.org/mentoring/requests/e01333ce426b474981391f0566b7a78d"}, {UUID: "f01333ce426b474981391f0566b7a78d", TrackTitle: "C++", ExerciseIconURL: "https://dg8krxphbh767.cloudfront.net/exercises/atbash-cipher.svg", ExerciseTitle: "Atbash Cipher", StudentHandle: "Student", StudentAvatarURL: "https://dg8krxphbh767.cloudfront.net/placeholders/user-avatar.svg", UpdatedAt: time.Now().UTC(), HaveMentoredPreviously: false, IsFavorited: false, Status: nil, URL: "https://exercism.org/mentoring/requests/f01333ce426b474981391f0566b7a78d"}}, Meta: request.Meta{CurrentPage: 1, TotalCount: 2, TotalPages: 1, UnscopedTotal: 2}},
		expectError: false,
		getClient: func(result request.MentoringRequestsResults) *ExercismHTTPClient {
			json, _ := json.Marshal(result)
			log.Println(string(json))
			return &ExercismHTTPClient{Client: &mockClient{expectedURL: "https://exercism.org/api/v2/mentoring/requests?track_slug=go&order=recent&page=", statusCode: http.StatusOK, expectedToken: "3ccb0df8-933c-4985-b1f8-b3c354b81313", responseBody: []io.ReadCloser{io.NopCloser(strings.NewReader(string(json)))}}, Token: "3ccb0df8-933c-4985-b1f8-b3c354b81313"}
		},
	},
	{
		description: "valid response with two pages",
		result:      request.MentoringRequestsResults{MentoringRequests: []request.MentoringRequest{{UUID: "e01333ce426b474981391f0566b7a78d", TrackTitle: "C++", ExerciseIconURL: "https://dg8krxphbh767.cloudfront.net/exercises/nth-prime.svg", ExerciseTitle: "Nth Prime", StudentHandle: "Student", StudentAvatarURL: "https://dg8krxphbh767.cloudfront.net/placeholders/user-avatar.svg", UpdatedAt: time.Now().UTC(), HaveMentoredPreviously: false, IsFavorited: false, Status: nil, URL: "https://exercism.org/mentoring/requests/e01333ce426b474981391f0566b7a78d"}, {UUID: "f01333ce426b474981391f0566b7a78d", TrackTitle: "C++", ExerciseIconURL: "https://dg8krxphbh767.cloudfront.net/exercises/atbash-cipher.svg", ExerciseTitle: "Atbash Cipher", StudentHandle: "Student", StudentAvatarURL: "https://dg8krxphbh767.cloudfront.net/placeholders/user-avatar.svg", UpdatedAt: time.Now().UTC(), HaveMentoredPreviously: false, IsFavorited: false, Status: nil, URL: "https://exercism.org/mentoring/requests/f01333ce426b474981391f0566b7a78d"}}, Meta: request.Meta{CurrentPage: 1, TotalCount: 2, TotalPages: 2, UnscopedTotal: 2}},
		expectError: false,
		getClient: func(result request.MentoringRequestsResults) *ExercismHTTPClient {
			mentoringRequests := result.MentoringRequests
			result.MentoringRequests = []request.MentoringRequest{mentoringRequests[0]}
			jsonPage1, _ := json.Marshal(result)

			result.MentoringRequests = []request.MentoringRequest{mentoringRequests[1]}
			result.Meta.CurrentPage = 2
			jsonPage2, _ := json.Marshal(result)

			result.MentoringRequests = mentoringRequests

			return &ExercismHTTPClient{Client: &mockClient{expectedURL: "https://exercism.org/api/v2/mentoring/requests?track_slug=go&order=recent&page=", statusCode: http.StatusOK, expectedToken: "3ccb0df8-933c-4985-b1f8-b3c354b81313",
				responseBody: []io.ReadCloser{
					io.NopCloser(strings.NewReader(string(jsonPage1))),
					io.NopCloser(strings.NewReader(string(jsonPage2)))},
			}, Token: "3ccb0df8-933c-4985-b1f8-b3c354b81313"}
		},
	},
	{
		description: "client returns error",
		result:      request.MentoringRequestsResults{},
		expectError: true,
		getClient: func(result request.MentoringRequestsResults) *ExercismHTTPClient {
			return &ExercismHTTPClient{Client: &mockClient{err: errors.New("test-error")}, Token: "3ccb0df8-933c-4985-b1f8-b3c354b81313"}
		},
	},
	{
		description: "invalid response-body (uuid is boolean in json instead of string)",
		result:      request.MentoringRequestsResults{},
		expectError: true,
		getClient: func(result request.MentoringRequestsResults) *ExercismHTTPClient {
			json := "{ \"results\": [ { \"uuid\": false } ] }"
			return &ExercismHTTPClient{Client: &mockClient{expectedURL: "https://exercism.org/api/v2/mentoring/requests?track_slug=go&order=recent&page=", statusCode: http.StatusOK, expectedToken: "3ccb0df8-933c-4985-b1f8-b3c354b81313", responseBody: []io.ReadCloser{io.NopCloser(strings.NewReader(json))}}, Token: "3ccb0df8-933c-4985-b1f8-b3c354b81313"}
		},
	},
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) { return 0, errors.New("test error") }

func (errReader) Close() error { return nil }

var testCasesGetRequest = []struct {
	description string
	expected    []byte
	expectError bool
	url         string
	client      ExercismHTTPClient
}{
	{
		description: "two valid mentoring requests",
		expected:    []byte("{\"results\":[{\"uuid\":\"e01333ce426b474981391f0566b7a78d\",\"track_title\":\"C++\",\"exercise_icon_url\":\"https://dg8krxphbh767.clou dfront.net/exercises/nth-prime.svg\",\"exercise_title\":\"Nth Prime\",\"student_handle\":\"Student\",\"student_avatar_url\":\"https://dg8krxphbh767.cloudfront.net/placeholders/user-ava tar.svg\",\"updated_at\":\"2022-06-21T17:20:27.8324484Z\",\"have_mentored_previously\":false,\"is_favorited\":false,\"status\":null,\"tooltip_url\":\"\",\"url\":\"https://exercism.org/ mentoring/requests/e01333ce426b474981391f0566b7a78d\"},{\"uuid\":\"f01333ce426b474981391f0566b7a78d\",\"track_title\":\"C++\",\"exercise_icon_url\":\"https://dg8krxphbh767.cloudfront.n et/exercises/atbash-cipher.svg\",\"exercise_title\":\"Atbash Cipher\",\"student_handle\":\"Student\",\"student_avatar_url\":\"https://dg8krxphbh767.cloudfront.net/placeholders/user-ava tar.svg\",\"updated_at\":\"2022-06-21T17:20:27.8324484Z\",\"have_mentored_previously\":false,\"is_favorited\":false,\"status\":null,\"tooltip_url\":\"\",\"url\":\"https://exercism.org/ mentoring/requests/f01333ce426b474981391f0566b7a78d\"}],\"meta\":{\"current_page\":1,\"total_count\":2,\"total_pages\":1,\"unscoped_total\":2}}"),
		expectError: false,
		url:         "https://exercism-mentoring-request-notifier/api/v2/mentoring/requests?track_slug=cpp?page=1",
		client:      ExercismHTTPClient{Client: &mockClient{expectedURL: "https://exercism-mentoring-request-notifier/api/v2/mentoring/requests?track_slug=cpp?page=", statusCode: http.StatusOK, expectedToken: "3ccb0df8-933c-4985-b1f8-b3c354b81313", responseBody: []io.ReadCloser{io.NopCloser(strings.NewReader("{\"results\":[{\"uuid\":\"e01333ce426b474981391f0566b7a78d\",\"track_title\":\"C++\",\"exercise_icon_url\":\"https://dg8krxphbh767.clou dfront.net/exercises/nth-prime.svg\",\"exercise_title\":\"Nth Prime\",\"student_handle\":\"Student\",\"student_avatar_url\":\"https://dg8krxphbh767.cloudfront.net/placeholders/user-ava tar.svg\",\"updated_at\":\"2022-06-21T17:20:27.8324484Z\",\"have_mentored_previously\":false,\"is_favorited\":false,\"status\":null,\"tooltip_url\":\"\",\"url\":\"https://exercism.org/ mentoring/requests/e01333ce426b474981391f0566b7a78d\"},{\"uuid\":\"f01333ce426b474981391f0566b7a78d\",\"track_title\":\"C++\",\"exercise_icon_url\":\"https://dg8krxphbh767.cloudfront.n et/exercises/atbash-cipher.svg\",\"exercise_title\":\"Atbash Cipher\",\"student_handle\":\"Student\",\"student_avatar_url\":\"https://dg8krxphbh767.cloudfront.net/placeholders/user-ava tar.svg\",\"updated_at\":\"2022-06-21T17:20:27.8324484Z\",\"have_mentored_previously\":false,\"is_favorited\":false,\"status\":null,\"tooltip_url\":\"\",\"url\":\"https://exercism.org/ mentoring/requests/f01333ce426b474981391f0566b7a78d\"}],\"meta\":{\"current_page\":1,\"total_count\":2,\"total_pages\":1,\"unscoped_total\":2}}"))}}, Token: "3ccb0df8-933c-4985-b1f8-b3c354b81313"},
	},
	{
		description: "invalid port as ! is directly added to the port (fails to create new http-request)",
		expected:    nil,
		expectError: true,
		url:         "https://exercism-mentoring-request-notifier:150!",
		client:      ExercismHTTPClient{Client: &mockClient{expectedToken: "3ccb0df8-933c-4985-b1f8-b3c354b81313", responseBody: []io.ReadCloser{io.NopCloser(strings.NewReader("{\"results\":[{\"uuid\":\"e01333ce426b474981391f0566b7a78d\",\"track_title\":\"C++\",\"exercise_icon_url\":\"https://dg8krxphbh767.clou dfront.net/exercises/nth-prime.svg\",\"exercise_title\":\"Nth Prime\",\"student_handle\":\"Student\",\"student_avatar_url\":\"https://dg8krxphbh767.cloudfront.net/placeholders/user-ava tar.svg\",\"updated_at\":\"2022-06-21T17:20:27.8324484Z\",\"have_mentored_previously\":false,\"is_favorited\":false,\"status\":null,\"tooltip_url\":\"\",\"url\":\"https://exercism.org/ mentoring/requests/e01333ce426b474981391f0566b7a78d\"},{\"uuid\":\"f01333ce426b474981391f0566b7a78d\",\"track_title\":\"C++\",\"exercise_icon_url\":\"https://dg8krxphbh767.cloudfront.n et/exercises/atbash-cipher.svg\",\"exercise_title\":\"Atbash Cipher\",\"student_handle\":\"Student\",\"student_avatar_url\":\"https://dg8krxphbh767.cloudfront.net/placeholders/user-ava tar.svg\",\"updated_at\":\"2022-06-21T17:20:27.8324484Z\",\"have_mentored_previously\":false,\"is_favorited\":false,\"status\":null,\"tooltip_url\":\"\",\"url\":\"https://exercism.org/ mentoring/requests/f01333ce426b474981391f0566b7a78d\"}],\"meta\":{\"current_page\":1,\"total_count\":2,\"total_pages\":1,\"unscoped_total\":2}}"))}, err: nil}, Token: "3ccb0df8-933c-4985-b1f8-b3c354b81313"},
	},
	{
		description: "client.Do returns error",
		expected:    nil,
		expectError: true,
		url:         "https://exercism-mentoring-request-notifier/api/v2/mentoring/requests?track_slug=cpp?page=2",
		client:      ExercismHTTPClient{Client: &mockClient{err: errors.New("testerror: client.Do failed")}},
	},
	{
		description: "read of body fails",
		expected:    nil,
		expectError: true,
		url:         "https://exercism-mentoring-request-notifier/api/v2/mentoring/requests?track_slug=cpp?page=1",
		client:      ExercismHTTPClient{Client: &mockClient{expectedURL: "https://exercism-mentoring-request-notifier/api/v2/mentoring/requests?track_slug=cpp?page=", responseBody: []io.ReadCloser{errReader(0)}, expectedToken: "3ccb0df8-933c-4985-b1f8-b3c354b81313"}, Token: "3ccb0df8-933c-4985-b1f8-b3c354b81313"},
	},
	{
		description: "invalid status-code (404)",
		expected:    nil,
		expectError: true,
		url:         "https://exercism-mentoring-request-notifier/unknown/path?page=1",
		client:      ExercismHTTPClient{Client: &mockClient{expectedURL: "https://exercism-mentoring-request-notifier/unknown/path?page=", statusCode: http.StatusNotFound, expectedToken: "3ccb0df8-933c-4985-b1f8-b3c354b81313", responseBody: []io.ReadCloser{io.NopCloser(strings.NewReader("given path not found"))}}, Token: "3ccb0df8-933c-4985-b1f8-b3c354b81313"},
	},
}
