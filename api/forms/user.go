package forms

// UserCreateForm is a form for the creation of a user.
type UserCreateForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserLoginForm is a form for login.
type UserLoginForm UserCreateForm
