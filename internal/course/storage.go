package course

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	COURSE_COLLECTION = "courses"
)

// How Course looks in MongoDB
type CourseDB struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name          string             `bson:"name" json:"name"`
	Code          string             `bson:"code" json:"code"`
	Prerequisites [][]string         `bson:"prerequisites" json:"prerequisites"`
	Corequisites  []string           `bson:"corequisites" json:"corequisites"`
	CrossListings []string           `bson:"crossListings" json:"crossListings"`
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

func (s *CourseStorage) FindCourseByNameOrCode(query string, limit int) ([]CourseDB, error) {
	collection := s.db.Collection(COURSE_COLLECTION)

	// Find the course by name or code
	filter := bson.M{
		"$or": bson.A{
			bson.M{"name": primitive.Regex{Pattern: query, Options: "i"}},
			bson.M{"code": primitive.Regex{Pattern: query, Options: "i"}},
		},
	}

	findOptions := options.Find().SetLimit(int64(limit))
	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.Background())

	// Decode the results
	var courses []CourseDB
	err = cursor.All(context.Background(), &courses)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (s *CourseStorage) FindCourseByID(id string) (*CourseDB, error) {
	collection := s.db.Collection(COURSE_COLLECTION)

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find the course by ID
	var course CourseDB
	err = collection.FindOne(context.Background(), bson.M{"_id": objId}).Decode(&course)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

// Course functions
func (c *CourseDB) Equal(other *CourseDB) bool {
	if c.Name != other.Name {
		return false
	}
	if c.Code != other.Code {
		return false
	}
	if len(c.Prerequisites) != len(other.Prerequisites) {
		return false
	}
	for i := range c.Prerequisites {
		if len(c.Prerequisites[i]) != len(other.Prerequisites[i]) {
			return false
		}
		for j := range c.Prerequisites[i] {
			if c.Prerequisites[i][j] != other.Prerequisites[i][j] {
				return false
			}
		}
	}
	if len(c.Corequisites) != len(other.Corequisites) {
		return false
	}
	for i := range c.Corequisites {
		if c.Corequisites[i] != other.Corequisites[i] {
			return false
		}
	}
	if len(c.CrossListings) != len(other.CrossListings) {
		return false
	}
	for i := range c.CrossListings {
		if c.CrossListings[i] != other.CrossListings[i] {
			return false
		}
	}
	return true
}

// for print
func (c *CourseDB) String() string {
	return fmt.Sprintf("%s %s %v %v %v", c.Code, c.Name, c.Prerequisites, c.Corequisites, c.CrossListings)
}
