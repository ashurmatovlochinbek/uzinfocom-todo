package request_objects

type RequestTask struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}
type RequestTaskForUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
	IsDone      bool   `json:"isDone"`
}
