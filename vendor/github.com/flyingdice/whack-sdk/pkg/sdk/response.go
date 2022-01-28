package sdk

type response struct {
	success interface{}
	err     error
}

func (r *response) Success() interface{} { return r.success }
func (r *response) Error() error         { return r.err }

func Success(val interface{}) *response {
	return &response{
		success: val,
	}
}

func Error(err error) *response {
	return &response{
		err: err,
	}
}
