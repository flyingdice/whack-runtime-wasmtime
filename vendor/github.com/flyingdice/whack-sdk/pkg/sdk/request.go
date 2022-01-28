package sdk

import "encoding/json"

type DeserializedRequest struct {
	ID      string
	Payload interface{}
}

type request struct {
	id      string
	payload interface{}
}

func (r *request) Bytes() ([]byte, error) {
	dsr := &DeserializedRequest{ID: "4444", Payload: r.payload}
	return json.Marshal(dsr)
	//switch p := r.payload.(type) {
	//case []byte:
	//	return p, nil
	//case string:
	//	return []byte(p), nil
	//default:
	//	return json.Marshal(p)
	//}
}

// convert bytes to sdk mem

func NewRequest(payload interface{}) *request {
	return &request{
		payload: payload,
	}
}
