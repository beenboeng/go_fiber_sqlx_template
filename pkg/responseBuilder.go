package pkg

type ResBuild[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

func Null() interface{} {
	return nil
}

func ResBuilder[T any](status int, message string, data T) ResBuild[T] {
	return ResBuild[T]{
		StatusCode: status,
		Message:    message,
		Data:       data,
	}
}
