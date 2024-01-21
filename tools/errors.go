package tools

type Err struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Extra   interface{} `json:"extra,omitempty"`
}

func (e *Err) Error() string {
	return e.Message
}
