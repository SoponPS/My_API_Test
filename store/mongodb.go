package store

import (
	"context"

	"github.com/Cnes-Consulting/backend_assignment/todo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDBStore struct {
	*mongo.Collection
}

func NewMongoDBStore(col *mongo.Collection) *mongoDBStore {
	return &mongoDBStore{Collection: col}
}

func (s *mongoDBStore) New(todo *todo.Todo) (*mongo.InsertOneResult, error) {
	result, err := s.Collection.InsertOne(context.Background(), todo)
	return result, err
}

func (s *mongoDBStore) Finding(id int) *mongo.SingleResult {
	result := s.Collection.FindOne(context.Background(), bson.M{"p_id":id})
	return result
}

func (s *mongoDBStore) FindAll() (*mongo.Cursor, error) {
	cur, err := s.Collection.Find(context.Background(), bson.M{})
	return cur, err
}

func (s *mongoDBStore) Deleting(id int) (*mongo.DeleteResult, error) {
	result, err := s.Collection.DeleteOne(context.Background(), bson.M{"p_id":id})
	return result, err
}

func (s *mongoDBStore) Updating(id int, todo *todo.Todo) (*mongo.UpdateResult, error) {
	updateData := bson.M{"$set": bson.M{"title":todo.Title, "iscomplete":todo.IsComplete}}
	result, err := s.Collection.UpdateOne(context.Background(), bson.M{"p_id":id}, updateData)
	return result, err
}

func (s *mongoDBStore) NewMany(todos []todo.Todo) (*mongo.InsertManyResult, error) {
	tasks := make([]interface{}, len(todos))
	for i := range todos {
		tasks[i] = todos[i]
	}
	result, err := s.Collection.InsertMany(context.Background(), tasks)
	return result, err
}