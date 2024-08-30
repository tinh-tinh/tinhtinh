package dto

type SignUpUser struct {
	Name string `validate:"required,isAlpha"`
}

type FindUser struct {
	Name string `validate:"isAlpha" query:"name"`
}
