package todo

import (
	"context"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	P_ID uint `json:"p_id"`
	Title string `json:"text" binding:"required"`
	IsComplete bool `json:"is_complete" binding:"required"`
}

func (Todo) TableName() string {
	return "tasks"
}

type storer interface {
	New(*Todo) (*mongo.InsertOneResult, error)
	FindAfterCreated(int) *mongo.SingleResult
	Finding(int) *mongo.SingleResult
	FindAll() (*mongo.Cursor, error)
	Deleting(int) (*mongo.DeleteResult, error)
	Updating(int, *Todo) (*mongo.UpdateResult, error)
	//NewMany([]Todo) (*mongo.InsertManyResult, error)
}

type TodoHandler struct {
	store storer
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store: store}
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	Param(string) string
}

// Create New Task
func (t *TodoHandler) NewTask(c Context) {
	// create
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	insertResult, err := t.store.New(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, insertResult.InsertedID)
}

//Show List Task
func (t *TodoHandler) ListTask(c Context) {
	var result []Todo
	cur, err := t.store.FindAll()
	//defer cur.Close(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	cur.All(context.Background(), &result)
	c.JSON(http.StatusOK, result)
}

// Show 1 Task
func (t *TodoHandler) ShowOneTask(c Context) {
	idParam := c.Param("p_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var result Todo
	showResult := t.store.Finding(id)
	if showResult.Err() != nil {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": showResult.Err().Error(),
		})
		return
	}

	err = showResult.Decode(&result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

// Delete 1 Task
func (t *TodoHandler) DeleteOneTask(c Context) {
	idParam := c.Param("p_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	deleteResult, err := t.store.Deleting(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if deleteResult.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "no document matching with this id",
		})
		return
	}

	c.JSON(http.StatusNoContent, map[string]interface{}{
		"DeleteCount": deleteResult.DeletedCount,
	})
}

// Update 1 Task
func (t *TodoHandler) UpdateTask(c Context) {
	idParam := c.Param("p_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	updateResult, err := t.store.Updating(id, &todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if updateResult.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"error": "no document matching with this id",
		})
		return
	}

	c.JSON(http.StatusNoContent, updateResult)
}

// Create New Many Tasks
/* func (t *TodoHandler) NewManyTask(c Context) {
	// create
	var todos []Todo
	if err := c.Bind(&todos); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	fmt.Println(todos)

	result, err := t.store.NewMany(todos)
	if err != nil {
		c.JSON(http.StatusNoContent, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusCreated, result.InsertedIDs)
} */