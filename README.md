# AWS APM Setup with CloudWatch Logs Integration

## Overview

This code provides a simple API for managing courses with AWS APM (Application Performance Monitoring) setup using CloudWatch Logs. It uses the Gorilla Mux router for handling HTTP requests and AWS SDK for Go for CloudWatch Logs integration.

## Prerequisites

Before running this code, make sure you have the following:

- AWS Access Key
- AWS Secret Key
- AWS Region

Replace the placeholder values in the code with your actual AWS Access Key, AWS Secret Key, and AWS Region.

## Setup

1. **Install dependencies:**
   ```bash
   go get -u github.com/aws/aws-sdk-go-v2
   go get -u github.com/gorilla/mux
## Configure AWS credentials:
Replace the placeholder values in the const block with your actual AWS Access Key, AWS Secret Key, and AWS Region.

## Create CloudWatch Logs:
Create a CloudWatch Logs group named "ubclogs" and a log stream named "ubc" in the AWS Management Console.

## Usage
Run the application:

bash
Copy code
go run main.go
Access the API:
Visit http://localhost:8000.

## API Routes
- GET /cources: Get a list of all courses.

## Middleware
- StartTransactionMiddleware
- Begins a new transaction for each incoming HTTP request.
- EndTransactionMiddleware
- Ends the transaction for each HTTP request, logging the transaction details to CloudWatch Logs.

## Controllers
- serverhome
- Responds with a simple message indicating that the API is running.
- getAllCorces
- Retrieves and returns a list of all courses.
- Logs messages and events to CloudWatch Logs, demonstrating CloudWatch Logs integration.


## APM Setup
- The code includes functions for starting and ending transactions and segments, as well as logging events and messages to --CloudWatch Logs.
- Ensure that your AWS credentials have the necessary permissions to write to CloudWatch Logs.

## Fake Database
- A simple in-memory database is used to store course information.
Logging
- The code includes functions for logging messages and events to CloudWatch Logs.
Note
- This code is for educational purposes and may need additional error handling and security measures for production use.
- Feel free to customize and extend the code based on your specific requirements and use case.
