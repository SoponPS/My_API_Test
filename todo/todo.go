package todo

import (
	"context"
	"net/http"
	"strconv"
	_"fmt"

	"go.mongodb.org/mongo-driver/bson"
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
	New(*Todo) error
	FindAfterCreated(int) (*mongo.Cursor, error)
	Finding(int) (*mongo.Cursor, error)
	FindAll() (*mongo.Cursor, error)
	Deleting(int) error
	Updating(int, string) (*mongo.UpdateResult, error)
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

	err := t.store.New(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	//show result
	var result []bson.M
	cur, err := t.store.FindAfterCreated(int(todo.P_ID))
	//defer cur.Close(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	cur.All(context.Background(), &result)
	c.JSON(http.StatusCreated, result)
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
	if idParam == "" {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "idparam error",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	var result []bson.M
	cur, err := t.store.Finding(id)
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

	err = t.store.Deleting(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, map[string]interface{}{
		"status": "success",
	})
}

// Update 1 Task
func (t *TodoHandler) UpdateTask(c Context) {
	idParam := c.Param("p_id")
	if idParam == "" {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "idparam error",
		})
		return
	}

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

	todotitle := todo.Title
	result, err := t.store.Updating(id, todotitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusNoContent, result)
}