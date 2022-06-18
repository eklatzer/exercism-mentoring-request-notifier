package collector

import "net/http"

type ExercismHttpClient struct {
	Client *http.Client
	Token  string
}
