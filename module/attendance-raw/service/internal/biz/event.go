package biz

type Event struct {
	Type string `json:"type"` // create, update, delete
	Key  string `json:"key"`
	Data string `json:"data"` // JSON string atau bisa pakai proto serialized
}
