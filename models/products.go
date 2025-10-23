package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `bson:"name" json:"name" validate:"required,min=2"`
	Price float64            `bson:"price" json:"price" validate:"required,gt=0"`
	Stock int                `bson:"stock" json:"stock" validate:"required,gte=0"`
}
