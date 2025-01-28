package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserSession struct {
	ExpiresAt int64 `json:"expires_at" bson:"expires_at"`
	IsAborted bool `json:"is_aborted" bson:"is_aborted"`
}

type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	Session UserSession `json:"-" bson:"session"`

}

type UserRegisterPayload struct {
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLoginPayload struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserStore interface{
	GetUserByID (id primitive.ObjectID) (*User, error)
	GetUserByEmail (email string) (*User, error)
	SaveUser (user *User) error
	UpdateUser (user *User) error
}
