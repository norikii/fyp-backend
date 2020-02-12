package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `bson:"_id"`
	Email string `bson:"email"`
	Password string `bson:"password"`
	Orders []Order `bson:"orders"`
	RememberToken string `bson:"remember_token"`
	IsStaff bool `bson:"is_staff"`
	EmailVerifiedAt int64 `bson:"email_verification_at"`
	InteractedAt int64 `bson:"interacted_at"`
	CreatedAt int64 `bson:"created_at"`
	UpdatedAt int64 `bson:"updated_at"`
	DeletedAt int64 `bson:"deleted_at"`
}
