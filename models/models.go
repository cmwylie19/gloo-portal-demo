package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name,omitempty" bson:"name,omitempty"`
	Status string             `json:"status,omitempty" bson:"status,omitempty"`
}

type UpdatedTask struct {
	ID     string
	Name   string
	Status string
}
