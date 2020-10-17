package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var client *mongo.Client

type Participants struct {
	Name string `json:"Name" bson:"Name"`
    Email string `json:"Email" son:"Email"`
    RSVP string `json:"RSVP" bson:"RSVP"`
}

type Meetings struct {
	Id primitive.ObjectID `json:"Id" bson:"Id"`
    Title string `json:"Title" bson:"Title"`
    Participant Participants `json:"Participants" bson:"Participants"`
    Start_Time Time `json:"Start Time" bson:"Start Time"`
    End_Time Time `json:"End Time" bson:"End Time"`
    Creation Timestamp `json:"Creation" bson:"Creation"`
}


func CreateMeeting(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var Meetings meeting
	_ = json.NewDecoder(request.Body).Decode(&Meetings)
	collection := client.Database("Database").Collection("Meetings")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, meeting)
	json.NewEncoder(response).Encode(result)
}

func GetMeeting(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meeting []Meetings
	collection := client.Database("Database").Collection("Meetings")
	fmt.Println("Enter ID of meeting")
	var Id string
	fmt.Scanln(&Id)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"Id"=Id})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
}

func GetMeeting_participant(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meeting []Meetings
	collection := client.Database("Database").Collection("Meetings")
	fmt.Println("Enter Name of the person")
	var name string
	fmt.Scanln(&name)
	a := []int{}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"Participants.Name"==name})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
}

func GetMeeting_within_time(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meeting []Meetings
	collection := client.Database("Database").Collection("Meetings")
	fmt.Println("Enter start time of duration:")
	var t1 Time
	fmt.Scanln(&Id)
	fmt.Println("Enter end time of duration:")
	var t2 Time
	fmt.Scanln(&Id)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{"Start_Time"<=t1 && "End_Time">=t2})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
}


defer func() {
    if err = client.Disconnect(ctx); err != nil {
        panic(err)
    }
}()

func main() {
	fmt.Println("Starting the application...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://chirag_435:<1xWIQ0XrYdkNj3NJ>@cluster0.0psuo.gcp.mongodb.net/<Database>?retryWrites=true&w=majority",))
    if err != nil { log.Fatal(err) }
	http.HandleFunc("/meetings", CreateMeeting).Methods("POST")
	http.HandleFunc("/meetings/{Id}", GetMeeting).Methods("GET")
	http.HandleFunc("/meetings?start={Start_Time}&end={End_Time}", GetMeeting_within_time).Methods("GET")
	http.HandleFunc("/meetings?participant={Email}", GetMeeting_participant).Methods("GET")
	http.ListenAndServe(":12345", router)
}
