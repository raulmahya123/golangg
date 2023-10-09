package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GeoJSONFeature struct {
	Type       string            `json:"type"`
	Geometry   GeoJSONGeometry   `json:"geometry"`
	Properties GeoJSONProperties `json:"properties"`
}

type GeoJSONGeometry struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

type GeoJSONProperties struct {
	Name string `json:"name"`
}

func CORSEnabledFunction(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Connect to your MongoDB instance
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://raulgantengbanget:0nGCVlPPoCsXNhqG@ac-oilbpwk-shard-00-00.9ofhjs3.mongodb.net:27017,ac-oilbpwk-shard-00-01.9ofhjs3.mongodb.net:27017,ac-oilbpwk-shard-00-02.9ofhjs3.mongodb.net:27017/test?replicaSet=atlas-13x7kp-shard-0&ssl=true&authSource=admin"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	// Access your MongoDB collection
	collection := client.Database("giss").Collection("giss")

	// Query your MongoDB collection to get the GeoJSON data
	var result GeoJSONFeature
	filter := bson.M{} // You can specify a filter here if needed
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serialize the GeoJSONFeature to JSON and write it to the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// Register the CORSEnabledFunction with the functions framework
	funcframework.RegisterHTTPFunction("/", CORSEnabledFunction)

	// Start the functions framework, listening on port 8080
	if err := funcframework.Start("8080"); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
