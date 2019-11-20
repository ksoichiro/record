package forms

// TaskCreateForm is a form for the creation of a task.
type TaskCreateForm struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        *bool  `json:"done"`
	Type        *int   `json:"type"`
	Amount      *int   `json:"amount"`
}

// TaskUpdateForm is a form for the update of the task.
type TaskUpdateForm struct {
	ID          *int    `json:"id" binding:"exists"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
	Type        *int    `json:"type"`
	Amount      *int    `json:"amount"`
}
