package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	batchSize     = 1000 // Process 1000 users at a time
	maxGoroutines = 10   // Limit to 10 concurrent uploads
)

var wg sync.WaitGroup
var semaphore = make(chan struct{}, maxGoroutines) // To control concurrency

func main() {
	// Load environment variables for AWS
	awsBucket := os.Getenv("AWS_BUCKET")
	bucketPrefix := "mybl-tests/" // Define your prefix
	awsEndpoint := os.Getenv("AWS_ENDPOINT")

	// Initialize the database connection
	db, err := InitDB()
	if err != nil {
		log.Fatal("Failed to initialize the database connection:", err)
	}
	defer db.Close()

	// Track the last processed ID
	lastID := 0

	// Create the directory if it doesn't exist
	outputDir := "./images" // Current working directory (project root)
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal("Error creating directory: ", err)
	}

	for {
		// Step 1: Fetch a batch of users
		rows, err := db.Query("SELECT id, name, username, profile_image_base64 FROM users WHERE profile_image_base64 != '' AND id > ? LIMIT ?", lastID, batchSize)
		if err != nil {
			log.Fatal("Error fetching data from MySQL: ", err)
		}

		// Check if there are no more rows to process
		if !rows.Next() {
			rows.Close()
			break
		}

		// Step 2: Iterate through the result set and log each row
		for rows.Next() {
			var id int
			var name, username, profile_image_base64 string

			// Scan the row into variables
			err = rows.Scan(&id, &name, &username, &profile_image_base64)
			if err != nil {
				log.Println("Error scanning row: ", err)
				continue
			}

			// Log the fetched data
			log.Printf("Fetched user - ID: %d, Name: %s, Username: %s\n", id, name, username)

			// Decode the base64 string.env
			imgData, err := base64.StdEncoding.DecodeString(profile_image_base64)
			if err != nil {
				log.Println("Error decoding base64 string: ", err)
				continue
			}

			// Prepare the S3 key
			s3Key := fmt.Sprintf("%suser_%d.png", bucketPrefix, id)

			// Concurrency control with semaphore and wait group
			semaphore <- struct{}{}
			wg.Add(1)

			go func(id int, imgData []byte, s3Key string) {
				defer wg.Done()
				defer func() { <-semaphore }()

				// Save the image locally
				imgFilePath := fmt.Sprintf("%s/user_%d.png", outputDir, id)
				err = ioutil.WriteFile(imgFilePath, imgData, 0644)
				if err != nil {
					log.Printf("Error saving image for user ID %d: %v\n", id, err)
					return
				}

				log.Printf("Image saved to: %s\n", imgFilePath)


				// Upload the image to S3
				err := uploadToS3(awsBucket, s3Key, imgData)
				if err != nil {
					log.Printf("Error uploading image for user ID %d: %v\n", id, err)
					return
				}

				log.Printf("Successfully uploaded image for user ID %d to S3: %s\n", id, s3Key)

				// Generate S3 URL
				s3URL := fmt.Sprintf("https://%s/%s/%s", awsEndpoint, awsBucket, s3Key)

				

				// Update the database with the S3 URL
				updateQuery := "UPDATE users SET profile_image = ? WHERE id = ?"
				_, err = db.Exec(updateQuery, s3URL, id)
				if err != nil {
					log.Printf("Error updating database for user ID %d: %v\n", id, err)
					return
				}

				log.Printf("Updated profile image path in database for user ID: %d\n", id)
			}(id, imgData, s3Key)

			// Update the last processed ID
			lastID = id
		}

		// Check for any error encountered during iteration
		if err = rows.Err(); err != nil {
			log.Fatal("Error iterating through rows: ", err)
		}

		rows.Close()
	}

	// Wait for all goroutines to finish before exiting
	wg.Wait()

	log.Println("All user images processed. Exiting the program...")
	os.Exit(0) // Exit with status code 0, indicating success
}

// Function to upload data to S3
func uploadToS3(bucket, key string, data []byte) error {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(os.Getenv("AWS_DEFAULT_REGION")),
		Credentials:      credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Endpoint:         aws.String(os.Getenv("AWS_ENDPOINT")),
		S3ForcePathStyle: aws.Bool(os.Getenv("AWS_USE_PATH_STYLE_ENDPOINT") == "true"),
	})
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("image/png"), // Change if necessary
		ACL:         aws.String("public-read"), // Set ACL as needed
	})
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}
