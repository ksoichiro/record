package forms

// RecordCreateForm is a form for the creation of a record.
type RecordCreateForm struct {
	TaskID *int `json:"task_id" binding:"exists"`
	Amount *int `json:"amount"`
}
