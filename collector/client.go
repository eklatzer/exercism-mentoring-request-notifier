package collector

import "net/http"

type ExercismHTTPClient struct {
	Client *http.Client
	Token  string
}
