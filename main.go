package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Cnes-Consulting/backend_assignment/router"
	"github.com/Cnes-Consulting/backend_assignment/store"
	"github.com/Cnes-Consulting/backend_assignment/todo"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variable: %s\n", err)
	}

	/* db, err := gorm.Open(mysql.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("fail to connect database")
	} */

	//db.AutoMigrate(&todo.Todo{})

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("DB_CONN")))
	if err != nil {
		panic("fail to connect database")
	}
	collection := client.Database("myapp").Collection("tasks")

	r := router.NewFiberRouter()
	//gormStore := store.NewGormStore(db)
	mongoStore := store.NewMongoDBStore(collection)
	handler := todo.NewTodoHandler(mongoStore)
	r.POST("/v1/tasks", handler.NewTask)
	r.GET("/v1/tasks", handler.ListTask)
	r.GET("/v1/tasks/:p_id", handler.ShowOneTask)
	r.DELETE("/v1/tasks/:p_id", handler.DeleteOneTask)
	r.PUT("/v1/tasks/:p_id", handler.UpdateTask)

	if err := r.Listen(":"+os.Getenv("PORT")); err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}