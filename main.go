package main

import (
	"fmt"
    "encoding/json"
	"log"
	"strings"
    "math/rand"
	//"path"
	"net/http"
	"net/url"
	"time"
	"context"
	//"io/ioutil"
    "strconv"
	// "github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
    //"go.mongodb.org/mongo-driver/mongo/readpref"
)

type server struct{}

// Book struct (Model)
type Meeting struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Participants  []Participant 
	StartTime  string  `json:"startTime"`
	EndTime  string  `json:"endTime"`
	CreationTimeStamp string `json:"creationTimeStamp"`
}

// Author struct
type Participant struct {
	Name string `json:"name"`
	Email  string `json:"email"`
	RSVP  string `json:"rsvp"`
}

// var participantsArray []Participants

// Init books var as a slice Book struct
var meetings []Meeting


var collection *mongo.Collection = ConnectDB()
func ConnectDB() *mongo.Collection {

	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://anmol:test123@cluster0.xod5m.mongodb.net/go_rest_api?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("go_rest_api").Collection("meetings")

	return collection
}


func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"message": "hello world"}`))
}

func allmeetings(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
	// Get all meetings
    case "GET":
        w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "get called"}`))
		
		w.Header().Set("Content-Type", "application/json")

		var meetings []Meeting
		cur, err := collection.Find(context.TODO(), bson.M{})

		if err != nil {
			log.Fatal(err)
			return
		}

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var meeting Meeting
			// & character returns the memory address of the following variable.
			err := cur.Decode(&meeting) // decode similar to deserialize process.
			if err != nil {
				log.Fatal(err)
			}
	
			// add item our array
			meetings = append(meetings, meeting)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}
	
		json.NewEncoder(w).Encode(meetings)
	// Post a meeting
    case "POST":
        w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "post called"}`))
		
		w.Header().Set("Content-Type", "application/json")
		var meeting Meeting
		_ = json.NewDecoder(r.Body).Decode(&meeting)
		meeting.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe

		result, err := collection.InsertOne(context.TODO(), meeting)

		if err != nil {	
			log.Fatal(err)
			w.Write([]byte(`{"message": "some database error occured"}`))
		}
		//meetings = append(meetings, meeting)

		json.NewEncoder(w).Encode(result)

    case "PUT":
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(`{"message": "put called"}`))
    case "DELETE":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "delete called"}`))
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "not found"}`))
    }
}

func singlemeeting(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
	// Get single meeting using id
    case "GET":
        w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "single get called"}`))
		
		w.Header().Set("Content-Type", "application/json")
		//params := mux.Vars(r) // Gets params
		//id := r.URL.Query().Get("id")
		parts := strings.Split(r.URL.String(), "/")
		//fmt.Printf(parts[2]);
		meetid := parts[2]
		var meeting Meeting

		id := meetid
		filter := bson.M{"id": id}
		err := collection.FindOne(context.TODO(), filter).Decode(&meeting)


		if err != nil {
			log.Fatal(err)
			return
		}

	    json.NewEncoder(w).Encode(meeting)
	
    case "POST":
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"message": "post called"}`))
    case "PUT":
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(`{"message": "put called"}`))
    case "DELETE":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "delete called"}`))
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "not found"}`))
    }
}


func specificMeetings(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
	// Get single meeting using id
    case "GET":
        w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "meetings with given email get called"}`))
		
		w.Header().Set("Content-Type", "application/json")

		u, _ := url.Parse(r.URL.RequestURI())
 
		values, _ := url.ParseQuery(u.RawQuery)
		
		email := values.Get("participant")
		fmt.Println("email:", email)
		
		st := values.Get("start")
		fmt.Println("st:", st)
		
		et := values.Get("end")
		fmt.Println("et:", et)

		currentTime := time.Now()
		fmt.Printf("Current time is: ", currentTime)

		var meetings []Meeting
		cur, err := collection.Find(context.TODO(), bson.M{})

		if err != nil {
			log.Fatal(err)
			return
		}

		defer cur.Close(context.TODO())

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var meeting Meeting
			// & character returns the memory address of the following variable.
			err := cur.Decode(&meeting) // decode similar to deserialize process.
			if err != nil {
				log.Fatal(err)
			}
	
			// add item our array
			meetings = append(meetings, meeting)

		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		for _, meeting := range meetings {
	
			for _, participant := range meeting.Participants{
				if participant.Email == email {
					json.NewEncoder(w).Encode(meeting)
					return
				}
			}
		}
	

	json.NewEncoder(w).Encode(meetings)

	
    case "POST":
        w.WriteHeader(http.StatusCreated)
        w.Write([]byte(`{"message": "post called"}`))
    case "PUT":
        w.WriteHeader(http.StatusAccepted)
        w.Write([]byte(`{"message": "put called"}`))
    case "DELETE":
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "delete called"}`))
    default:
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "not found"}`))
    }
}

// Main function
func main() {

	// Hardcoded data - @todo: add database


	fileServer := http.FileServer(http.Dir("./static")) // New code
	http.Handle("/", fileServer) // New code

    http.HandleFunc("/meetings", allmeetings)
    http.HandleFunc("/meeting/", singlemeeting)
    //http.HandleFunc("/meetings?startTime=st&endTime:et", specificMeetings)
    http.HandleFunc("/meetings/", specificMeetings)


    fmt.Printf("Starting server at port 8080\n")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
	}

}