package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cmwylie19/gloo-portal-demo/models"
	"github.com/cmwylie19/gloo-portal-demo/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Connection mongoDB with utils class
var collection = utils.ConnectDB()

/*
 * GET /api/tasks
 * description: get all tasks
 */
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var tasks []models.Task

	cur, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		utils.GetError(err, w)
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var task models.Task

		err := cur.Decode(&task)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	json.NewEncoder(w).Encode(tasks)
}

/*
 * GET /api/tasks/{id}
 * description: get a task
 */
func getTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task models.Task
	var params = mux.Vars(r)

	id, _ := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(task)
}

/*
 * POST /api/tasks
 * description: create a task
 */
func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task models.Task
	_ = json.NewDecoder(r.Body).Decode(&task)
	result, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(result)
}

/*
 * PUT /api/tasks/{id}
 * description: update a task
 */
func updateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task models.Task
	var params = mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	filter := bson.M{"_id": id}
	_ = json.NewDecoder(r.Body).Decode(&task)
	fmt.Println("User set status", task.Status)
	status := task.Status
	update := bson.D{
		{"$set", bson.D{
			{"Name", task.Name},
			{"Status", task.Status},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&task)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	task.ID = id
	task.Status = status
	json.NewEncoder(w).Encode(task)

}

/*
 * DELETE /api/tasks/{id}
 * description: delete a task
 */
func deleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params = mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	filter := bson.M{"_id": id}
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
}

/*
 * DELETE /api/tasks
 * description: delete all tasks
 */
func deleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	filter := bson.M{}
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		utils.GetError(err, w)
		return
	}
	json.NewEncoder(w).Encode(deleteResult)
}

type Remote struct {
	XFF string `json:"x-forwarded-for"`
}

/*
 * GET /remote
 * description: get x-forwarded-for header
 */
func getRemote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	remote := Remote{XFF: r.Header.Get("x-forwarded-for")}
	result, err := json.Marshal(remote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/tasks", getTasks).Methods("GET")
	r.HandleFunc("/api/tasks/{id}", getTask).Methods("GET")
	r.HandleFunc("/api/tasks", createTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/api/tasks", deleteAll).Methods("DELETE")
	r.HandleFunc("/api/tasks/{id}", deleteTask).Methods("DELETE")
	r.HandleFunc("/remote", getRemote).Methods("GET")

	config := utils.GetConfiguration()

	log.Fatal(http.ListenAndServe(config.Port, r))

}
