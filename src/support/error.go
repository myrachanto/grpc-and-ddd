package support

type HttpError struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
	Errors  string `json:"errors,omitempty"`
	Results bool   `json:"results,omitempty"`
}
