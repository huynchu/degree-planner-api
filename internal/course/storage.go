package course

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// How Course looks in MongoDB
type CourseDB struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string             `bson:"name" json:"name"`
	Code          string             `bson:"code" json:"code"`
	Prerequisites []string           `bson:"prerequisites" json:"prerequisites"`
	Corequisites  []string           `bson:"corequisites" json:"corequisites"`
}

type CourseStorage struct {
	// cache map[string]*Course (this would be redis)
	db *mongo.Database
}

func NewCourseStorage(db *mongo.Database) *CourseStorage {
	return &CourseStorage{
		db: db,
	}
}

func (s *CourseStorage) CreateCourse(course *CourseDB) (string, error) {
	collection := s.db.Collection("Courses")

	// Insert the course
	insertResult, err := collection.InsertOne(context.Background(), course)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *CourseStorage) FindCourseByID(id string) (*CourseDB, error) {
	collection := s.db.Collection("Courses")

	// Find the course by ID
	filter := bson.M{"_id": id}
	var course CourseDB
	err := collection.FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		return nil, err
	}

	return &course, nil
}
