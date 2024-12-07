package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title       string            `bson:"title" json:"title"`
    Description string            `bson:"description" json:"description"`
    Instructor  string            `bson:"instructor" json:"instructor"`
    CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
}

type Student struct {
    ID              primitive.ObjectID   `bson:"_id" json:"id"`
    Name            string              `bson:"name" json:"name"`
    Email           string              `bson:"email" json:"email"`
    EnrolledCourses []primitive.ObjectID `bson:"enrolled_courses" json:"enrolled_courses"`
    CreatedAt       time.Time           `bson:"created_at" json:"created_at"`
}