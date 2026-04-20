package errors

// Problem is the RFC 7807 problem-details JSON representation.
type Problem struct {
	Type     string         `json:"type"`
	Title    string         `json:"title"`
	Status   int            `json:"status"`
	Detail   string         `json:"detail,omitempty"`
	Instance string         `json:"instance,omitempty"`
	Domain   string         `json:"domain,omitempty"`
	Code     string         `json:"code,omitempty"`
	Errors   map[string]any `json:"errors,omitempty"`
}

// ToProblem renders e as an RFC 7807 problem document.
func (e *Error) ToProblem(instance string) Problem {
	return Problem{
		Type:     "https://helixgitpx.dev/errors/" + e.Code.String(),
		Title:    e.Code.String(),
		Status:   e.HTTPStatus(),
		Detail:   e.Message,
		Instance: instance,
		Domain:   e.Domain,
		Code:     e.Code.String(),
		Errors:   e.Details,
	}
}
