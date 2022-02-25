package store

import (
	"context"

	"github.com/Cnes-Consulting/backend_assignment/todo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDBStore struct {
	*mongo.Collection
}

func NewMongoDBStore(col *mongo.Collection) *mongoDBStore {
	return &mongoDBStore{Collection: col}
}

func (s *mongoDBStore) New(todo *todo.Todo) error {
	_, err := s.Collection.InsertOne(context.Background(), todo)
	return err
}

func (s *mongoDBStore) FindAfterCreated(id int) (*mongo.Cursor, error) {
	opts := options.Find().SetProjection(bson.M{"p_id": 0, "title": 0, "is_complete": 0})
	cur, err := s.Collection.Find(context.Background(), bson.M{"p_id":id}, opts)
	return cur, err
}

func (s *mongoDBStore) Finding(id int) (*mongo.Cursor, error) {
	cur, err := s.Collection.Find(context.Background(), bson.M{"p_id":id})
	return cur, err
}

func (s *mongoDBStore) FindAll() (*mongo.Cursor, error) {
	cur, err := s.Collection.Find(context.Background(), bson.M{})
	return cur, err
}

func (s *mongoDBStore) Deleting(id int) error {
	_, err := s.Collection.DeleteOne(context.Background(), bson.M{"p_id":id})
	return err
}

func (s *mongoDBStore) Updating(id int, todotitle string) (*mongo.UpdateResult, error) {
	result, err := s.Collection.UpdateOne(context.Background(), bson.M{"p_id":id}, bson.M{"$set": bson.M{"title":todotitle}})
	return result, err
}