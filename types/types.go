package types

type ApiError struct {
	Message   string `json:"message"`
	Body      []byte `json:"body"`
	ErrorCode int    `json:"error_code"`
}

func (a *ApiError) Error() string {
	//TODO implement me
	panic("implement me")
}
