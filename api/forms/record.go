package forms

// RecordCreateForm is a form for the creation of a record.
type RecordCreateForm struct {
	TaskID *int  `json:"task_id" binding:"exists"`
	Done   *bool `json:"done"`
	Amount *int  `json:"amount"`
}

// RecordUpdateForm is a form for the update of the record.
type RecordUpdateForm struct {
	ID     *int  `json:"id" binding:"exists"`
	Done   *bool `json:"done"`
	Amount *int  `json:"amount"`
}
