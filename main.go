package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"

	"github.com/gorilla/mux"
)

// Models
type Cource struct {
	CourceID string `json:"cource_id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Language string `json:"language"`
}

const (
	awsAccessKey = "your AWS Access Key"
	awsSecretKey = "your AWS Secret Key"
	awsRegion    = "your AWS Region"
)

// APM Setup

// Event Struct
type Event struct {
	Details    Transaction
	DetailType string
	Source     string
}

// Transaction

type Transaction struct {
	TransactionID string    `json:"transaction_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Segments      []Segment `json:"segments"`
}

// Segment struct
type Segment struct {
	SegmentID string    `json:"segment_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

func StartTransaction() *Transaction {
	return &Transaction{
		TransactionID: fmt.Sprintf("tx-%d", time.Now().UnixNano()),
		StartTime:     time.Now(),
		Segments:      []Segment{},
	}
}

func EndTransaction(tx *Transaction) {
	tx.EndTime = time.Now()

	// Log the transaction
	event := Event{
		Details:    *tx,
		DetailType: "TransactionEvent",
		Source:     "ubc-local",
	}

	client := getCloudWatchLogsClient()
	logGroupName := "ubclogs"
	logStreamName := "ubc"

	// Log event
	if err := sendEventToCloudWatchLogs(client, logGroupName, logStreamName, event); err != nil {
		log.Println("ERROR #EndTransaction #sending transaction event to CloudWatch Logs: #", err)
	}
}

func StartSegment(tx *Transaction) *Segment {
	segment := Segment{
		SegmentID: fmt.Sprintf("seg-%d", time.Now().UnixNano()),
		StartTime: time.Now(),
	}

	tx.Segments = append(tx.Segments, segment)

	return &segment
}

func EndSegment(segment *Segment) {
	segment.EndTime = time.Now()
}

// Start Transaction
func StartTransactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx := StartTransaction()
		ctx := context.WithValue(r.Context(), "transaction", tx)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// End Transaction
func EndTransactionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, ok := r.Context().Value("transaction").(*Transaction)
		if ok {
			EndTransaction(tx)
		}
		next.ServeHTTP(w, r)
	})
}

// get logs from cloudwatch
func getLogs(logGroupName string, logStreamNames []string) {
	client := getCloudWatchLogsClient()

	endTime := time.Now()
	startTime := endTime.Add(-1 * time.Hour)

	for _, logStreamName := range logStreamNames {
		input := &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:   &logGroupName,
			LogStreamNames: []string{logStreamName},
			StartTime:      aws.Int64(startTime.Unix() * 1000),
			EndTime:        aws.Int64(endTime.Unix() * 1000),
		}

		resp, err := client.FilterLogEvents(context.Background(), input)
		if err != nil {
			log.Printf("Error fetching logs for log stream '%s': %v\n", logStreamName, err)
			continue
		}

		fmt.Printf("Logs for log stream '%s':\n", logStreamName)
		for _, event := range resp.Events {
			fmt.Println(*event.Message)
		}
		fmt.Println()
	}
}

// Getting client for logs
func getCloudWatchLogsClient() *cloudwatchlogs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		log.Fatal("ERROR #AWS config : #", err)
	}

	cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, ""))

	return cloudwatchlogs.NewFromConfig(cfg)
}

// Getting client for event
// func getCloudWatchEventsClient() *cloudwatchevents.Client {
// 	cfg, err := config.LoadDefaultConfig(context.TODO(),
// 		config.WithRegion(awsRegion),
// 	)
// 	if err != nil {
// 		log.Fatal("Unable to load AWS config:", err)
// 	}

// 	cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, ""))

// 	return cloudwatchevents.NewFromConfig(cfg)
// }

// Simple logging to AWS
/*------------------------------------------------------------------------------------------------------------------*/
func logToCloudWatchLogs(client *cloudwatchlogs.Client, logGroupName, logStreamName, message string) error {
	input := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &logGroupName,
		LogStreamName: &logStreamName,
		LogEvents: []types.InputLogEvent{
			{
				Message:   &message,
				Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
			},
		},
	}

	_, err := client.PutLogEvents(context.Background(), input)
	return err
}
func sendEventToCloudWatchLogs(client *cloudwatchlogs.Client, logGroupName, logStreamName string, eventData Event) error {
	eventDetails, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	input := &cloudwatchlogs.PutLogEventsInput{
		LogGroupName:  &logGroupName,
		LogStreamName: &logStreamName,
		LogEvents: []types.InputLogEvent{
			{
				Message:   aws.String(string(eventDetails)),
				Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
			},
		},
	}

	_, err = client.PutLogEvents(context.Background(), input)
	return err
}

// Fake DataBase
var cources []Cource

// Middleware and Helper
func (c *Cource) isEmpty() bool {
	return c.CourceID == "" && c.Name == ""
}
func main() {
	fmt.Println("API Routes")

	r := mux.NewRouter()
	r.Use(StartTransactionMiddleware)
	r.Use(EndTransactionMiddleware)
	cources = append(cources, Cource{
		CourceID: "1",
		Name:     "Go",
		Price:    10,
		Language: "Go",
	})
	cources = append(cources, Cource{
		CourceID: "2",
		Name:     "Python",
		Price:    20,
		Language: "Python",
	})
	r.HandleFunc("/", serverhome).Methods("GET")
	r.HandleFunc("/cources", getAllCorces).Methods("GET")

	// Listen Port
	log.Fatal(http.ListenAndServe(":8000", r))
}

// Controllers
//////// Server Home Route

func serverhome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is running"))
}

func getAllCorces(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get All Courses")
	w.Header().Set("Content-Type", "application/json")

	// Start a new segment and also end by defer
	tx, ok := r.Context().Value("transaction").(*Transaction)
	if ok {
		segment := StartSegment(tx)
		defer EndSegment(segment)
	}

	// Log to CloudWatch Logs
	logGroupName := "ubclogs"
	logStreamName := "ubc"
	logMessage := "log is running"

	event := Event{
		Details:    Transaction{TransactionID: tx.TransactionID},
		DetailType: "TransactionEvent",
		Source:     "YourApp",
	}

	client := getCloudWatchLogsClient()

	// Log message
	if err := logToCloudWatchLogs(client, logGroupName, logStreamName, logMessage); err != nil {
		log.Println("Error logging message to CloudWatch Logs:", err)
	}

	// Log event
	if err := sendEventToCloudWatchLogs(client, logGroupName, logStreamName, event); err != nil {
		log.Println("Error sending event to CloudWatch Logs:", err)
	}

	json.NewEncoder(w).Encode(cources)

	logStreamNames := []string{logStreamName}

	getLogs(logGroupName, logStreamNames)
}
