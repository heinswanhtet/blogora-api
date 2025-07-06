package types

type RegisterUserPayload struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthorPayload struct {
	Name     *string `json:"name" validate:"required"`
	Username *string `json:"username" validate:"required"`
	Email    *string `json:"email" validate:"required,email"`
	Image    *string `json:"image"`
	Bio      *string `json:"bio"`
}
