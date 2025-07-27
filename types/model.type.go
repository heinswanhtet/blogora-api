package types

import (
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

type Author struct {
	ID        *string    `json:"id,omitempty"`
	Name      *string    `json:"name,omitempty"`
	Username  *string    `json:"username,omitempty"`
	Email     *string    `json:"email,omitempty"`
	Image     *string    `json:"image,omitempty"`
	Bio       *string    `json:"bio,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitzero"`
	UpdatedAt *time.Time `json:"updated_at,omitzero"`
	CreatedBy *string    `json:"created_by,omitempty"`
	UpdatedBy *string    `json:"updated_by,omitempty"`
	Deleted   *string    `json:"deleted,omitempty"`
}

type Startup struct {
	ID          *string    `json:"id"`
	Title       *string    `json:"title"`
	Slug        *string    `json:"slug"`
	AuthorId    *string    `json:"author_id"`
	Author      Author     `json:"author,omitzero"`
	Views       *int       `json:"views"`
	Description *string    `json:"description"`
	Category    *string    `json:"category"`
	Image       *string    `json:"image"`
	Pitch       *string    `json:"pitch"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedBy   *string    `json:"created_by"`
	UpdatedBy   *string    `json:"updated_by"`
	Deleted     *string    `json:"deleted"`
}

type Playlist struct {
	ID        *string    `json:"id"`
	Title     *string    `json:"title"`
	Slug      *string    `json:"slug"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	CreatedBy *string    `json:"created_by"`
	UpdatedBy *string    `json:"updated_by"`
	Deleted   *string    `json:"deleted"`
}

type StartupPlaylist struct {
	StartupID  *string    `json:"startup_id"`
	PlaylistId *string    `json:"playlist_id"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	CreatedBy  *string    `json:"created_by"`
	UpdatedBy  *string    `json:"updated_by"`
	Deleted    *string    `json:"deleted"`
}
