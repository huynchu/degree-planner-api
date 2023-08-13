package degree

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// How Course looks in MongoDB
type DegreeDB struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Semesters []Semester         `bson:"semesters" json:"semesters"`
	Owner     primitive.ObjectID `bson:"owner,omitempty" json:"owner,omitempty"`
}

type Semester struct {
	Name    string               `bson:"name" json:"name"`
	Courses []primitive.ObjectID `bson:"courses" json:"courses"`
}

type DegreeStorage struct {
	// cache map[string]*Course (this would be redis)
	db *mongo.Database
}

func NewDegreeStorage(db *mongo.Database) *DegreeStorage {
	return &DegreeStorage{
		db: db,
	}
}

func (d *DegreeStorage) CreateDegree(name string) (string, error) {
	collection := d.db.Collection("degree")

	newDegree := DegreeDB{
		Name:      name,
		Semesters: []Semester{},
	}

	// Insert the course
	insertResult, err := collection.InsertOne(context.Background(), newDegree)
	if err != nil {
		return "", err
	}

	return insertResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d *DegreeStorage) FindDegreeByID(id string) (*DegreeDB, error) {
	collection := d.db.Collection("degree")

	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find the course
	var degree DegreeDB
	err = collection.FindOne(context.Background(), primitive.M{"_id": objId}).Decode(&degree)
	if err != nil {
		return nil, err
	}

	return &degree, nil
}

func (d *DegreeStorage) AddSemester(degreeID string, semesterName string, semesterIndex int) error {
	collection := d.db.Collection("degree")

	objId, err := primitive.ObjectIDFromHex(degreeID)
	if err != nil {
		return err
	}

	newSemester := Semester{
		Name:    semesterName,
		Courses: []primitive.ObjectID{},
	}

	// Add the semester
	filter := primitive.M{"_id": objId}
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		primitive.M{
			"$push": primitive.M{
				"semesters": primitive.M{
					"$each":     []Semester{newSemester},
					"$position": semesterIndex,
				},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (d *DegreeStorage) UpdateSemesters(degreeID string, semesters []Semester) error {
	collection := d.db.Collection("degree")

	objId, err := primitive.ObjectIDFromHex(degreeID)
	if err != nil {
		return err
	}

	// Update the semesters
	filter := primitive.M{"_id": objId}
	_, err = collection.UpdateOne(
		context.Background(),
		filter,
		primitive.M{
			"$set": primitive.M{
				"semesters": semesters,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}
