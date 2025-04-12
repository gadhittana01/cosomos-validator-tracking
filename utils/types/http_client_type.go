package types

type HTTPResponse struct {
	StatusCode int
	Body       string
	Headers    map[string][]string
}
