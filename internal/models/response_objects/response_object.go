package response_objects

type ResponseObject struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponseObject(status bool, message string, data interface{}) *ResponseObject {
	return &ResponseObject{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
