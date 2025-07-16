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

type StartupPayload struct {
	Title       *string `json:"title" validate:"required"`
	AuthorId    *string `json:"author_id" validate:"required"`
	Description *string `json:"description"`
	Category    *string `json:"category" validate:"required"`
	Image       *string `json:"image"`
	Pitch       *string `json:"pitch"`
}

type PlaylistPayload struct {
	Title *string `json:"title" validate:"required"`
}
