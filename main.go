package main

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection, postCollection = ConnecttoDB()

func main() {
	//Setting Up All The Routes
	http.HandleFunc("/user", createUser)
	http.HandleFunc("/user/", findUser)
	http.HandleFunc("/posts", createPost)
	http.HandleFunc("/posts/", findPost)
	http.HandleFunc("/posts/users/", userPost)

	// set our port address as 8081
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func ConnecttoDB() (*mongo.Collection, *mongo.Collection) {

	// Set client options
	//change the URI according to your database
	clientOptions := options.Client().ApplyURI("mongodb+srv://harshita_jaiswani_:abc_123@cluster0.v1ctf.mongodb.net/myFirstDatabase?retryWrites=true&w=majority")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	//Error Handling
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	//DB collection address which we are going to use
	//available to functions of all scope
	collection := client.Database("Appointy").Collection("NewsUsers")
	postCollection := client.Database("Appointy").Collection("NewPosts")
	return collection, postCollection
}

//Function to create a new User in Database

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Checking if the method is the post method or not
	if r.Method == "POST" {
		//Parsing the Values
		r.ParseForm()

		//Hashing the Passowrd
		pass := md5.Sum([]byte(r.Form.Get("email")))

		//Pushing the data into DB
		result, err := collection.InsertOne(context.TODO(), bson.D{
			{Key: "name", Value: r.Form.Get("name")},
			{Key: "passowrd", Value: string(pass[:])},
			{Key: "email", Value: r.Form.Get("email")},
			{Key: "userId", Value: r.Form.Get("id")},
		})

		//Error Check
		if err != nil {
			log.Fatal(err)
		}

		//Sending the id to back
		json.NewEncoder(w).Encode(result)
	} else {
		//If it is not a post Request
		//Then Showing appropiate error
		fmt.Fprintf(w, "Only Post Request Accepted")
		return
	}

}

//Function to Find User
func findUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Extracting the Path of URL
	url := r.URL.RequestURI()
	key := path.Base(url)

	var result bson.M

	//Finding the user
	err := collection.FindOne(context.TODO(), bson.M{"userId": key}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	//Sending the data Back
	json.NewEncoder(w).Encode(result)
}

//Function to create a POST
func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Checking the Method is POST
	if r.Method == "POST" {
		//Parsing the Form
		r.ParseForm()

		//Adding the timestamp
		dt := time.Now()
		dt.Format("01-02-2006 15:04:05")

		//Pushing the Data into DB
		result, err := postCollection.InsertOne(context.TODO(), bson.D{
			{Key: "id", Value: r.Form.Get("pid")},
			{Key: "caption", Value: r.Form.Get("caption")},
			{Key: "url", Value: r.Form.Get("url")},
			{Key: "time", Value: dt.String()},
			{Key: "userId", Value: r.Form.Get("userId")},
		})

		if err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(result)
	} else {
		fmt.Fprintf(w, "Only Post Request Accepted")
		return
	}

}

//Function to find a POST
func findPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Extracting the ID of the post
	url := r.URL.RequestURI()
	key := path.Base(url)

	var result bson.M

	//Finding the post Based on that KEY
	err := postCollection.FindOne(context.TODO(), bson.M{"id": key}).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(result)
}

func userPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Extracting the UserId
	url := r.URL.RequestURI()
	key := path.Base(url)

	//Finding the posts made by that user
	filterCursor, err := postCollection.Find(context.TODO(), bson.M{"userId": key})
	if err != nil {
		log.Fatal(err)
	}

	//Appending it to the array and sending it back
	var postFilter []bson.M
	if err = filterCursor.All(context.TODO(), &postFilter); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(postFilter)

}
