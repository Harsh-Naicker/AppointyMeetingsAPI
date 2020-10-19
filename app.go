package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Participant struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
  Email     string             `json:"email,omiempty" bson:"email,omiempty"`
  RSVP      string             `json:"rsvp,omiempty" bson:"rsvp,omiempty"`

}
type Meeting struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string             `json:"title,omitempty" bson:"title,omitempty"`
	Participants  string             `json:"participants,omitempty" bson:"participants,omitempty"`
  StartTime     string             `json:"starttime,omiempty" bson:"starttime,omiempty"`
  EndTime      string             `json:"endtime,omiempty" bson:"endtime,omiempty"`

}
func CreateMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meeting Meeting
	_ = json.NewDecoder(request.Body).Decode(&meeting)
	collection := client.Database("themeetingapi").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, meeting)
	json.NewEncoder(response).Encode(result)
}
func GetMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var meeting Meeting
	collection := client.Database("themeetingapi").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&meeting)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
}
func GetMeetingsEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meetings []Meeting
	collection := client.Database("themeetingapi").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		meetings = append(meetings, meeting)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meetings)
}


func CreatePersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var participant Participant
	_ = json.NewDecoder(request.Body).Decode(&participant)
	collection := client.Database("themeetingapi").Collection("participants")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, participant)
	json.NewEncoder(response).Encode(result)
}
func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var participant Participant
	collection := client.Database("themeetingapi").Collection("participants")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Participant{ID: id}).Decode(&participant)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(participant)
}
func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var people []Participant
	collection := client.Database("themeetingapi").Collection("participants")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var participant Participant
		cursor.Decode(&participant)
		people = append(people, participant)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/participant", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/participant/{id}", GetPersonEndpoint).Methods("GET")
  router.HandleFunc("/meeting", CreateMeetingEndpoint).Methods("POST")
	router.HandleFunc("/meetings", GetMeetingsEndpoint).Methods("GET")
	router.HandleFunc("/meeting/{id}", GetMeetingEndpoint).Methods("GET")
	http.ListenAndServe(":8000", router)
}
