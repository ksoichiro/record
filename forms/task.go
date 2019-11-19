package forms

type TaskCreateForm struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        *bool  `json:"done"`
	Type        *int   `json:"type"`
	Amount      *int   `json:"amount"`
}
