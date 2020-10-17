package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Participant struct {
	Name string `json:"firstname,omitempty" bson:"firstname,omitempty"`     
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	  }


type Meeting struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title string `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Participants []Participant `json:"participants,omitempty" bson:"participants,omitempty"`
	Starttime string `json:"starttime,omitempty" bson:"starttime,omitempty"`
	Endtime string `json:"endtime,omitempty" bson:"endtime,omitempty"`
	Creationtime string `json:"creationtime,omitempty" bson:"creationtime,omitempty"`
	}


func CreateMeetingEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meeting Meeting
	_ = json.NewDecoder(request.Body).Decode(&meeting)
	collection := client.Database("AppointyMeetings").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, meeting)
	json.NewEncoder(response).Encode(result)
	fmt.Println(meeting)
}

func GetMeetingbyIdEndpoint(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		response.Header().Set("content-type", "application/json")
		fmt.Println(path.Base(request.URL.Path))
		id, _ := primitive.ObjectIDFromHex(path.Base(request.URL.Path))
		var meeting Meeting
		collection := client.Database("AppointyMeetings").Collection("meetings")
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := collection.FindOne(ctx, Meeting{Id: id}).Decode(&meeting)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
			return
		}
		json.NewEncoder(response).Encode(meeting)
	}
}






func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	http.HandleFunc("/meeting", CreateMeetingEndpoint)
	http.HandleFunc("/meeting/{id}", GetMeetingbyIdEndpoint)
	http.ListenAndServe(":12345", nil)
}

