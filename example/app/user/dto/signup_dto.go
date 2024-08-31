package dto

type SignUpUser struct {
	Name     string `validate:"required,isAlpha"`
	Email    string `validate:"required,isEmail"`
	Password string `validate:"required,isStrongPassword"`
}

type FindUser struct {
	Name string `validate:"isAlpha" query:"name"`
}
