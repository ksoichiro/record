package forms

type RecordCreateForm struct {
	TaskID *int `json:"task_id" binding:"exists"`
	Amount *int `json:"amount"`
}
