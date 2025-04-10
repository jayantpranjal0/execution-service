package queries

import (
	"context"
	"execution-service/internal/models"
	"log"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func AddEntry(collection *mongo.Collection, entry models.ExecutedJob) error {
	// Insert the entry into the collection
	result, err := collection.InsertOne(context.TODO(),entry)
	if err != nil {
		log.Printf("Error adding entry: %v", err)
		return err
	}

	log.Printf("Added entry with ID: %v", result.InsertedID)
	return nil
}

// // UpdateJobStatus updates the status of a job in the database
// func UpdateJobStatus(collection *mongo.Collection, jobID string, newStatus string) error {
// 	// Create a filter to find the job by job_id
// 	filter := bson.M{"job_id": jobID}

// 	// Define the update to set the new status
// 	update := bson.M{
// 		"$set": bson.M{
// 			"job_status": newStatus,
// 		},
// 	}

// 	// Perform the update operation
// 	result, err := collection.UpdateOne(context.TODO(), filter, update, options.Update())
// 	if err != nil {
// 		log.Printf("Error updating job status for job_id %s: %v", jobID, err)
// 		return err
// 	}

// 	// Log the result
// 	if result.MatchedCount == 0 {
// 		log.Printf("No job found with job_id %s", jobID)
// 	} else {
// 		log.Printf("Updated job status for job_id %s to %s", jobID, newStatus)
// 	}

// 	return nil
// }


