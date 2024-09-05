package dto

type SignUpUser struct {
	Name     string `validate:"isAlpha" example:"John"`
	Email    string `validate:"required,isEmail" example:"john@gmail.com"`
	Password string `validate:"required,isStrongPassword" example:"12345678@Tc"`
}

type FindUser struct {
	Name string `validate:"required,isAlpha" query:"name" example:"ac"`
}
