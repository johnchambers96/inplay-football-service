package main

import (
	"fmt"
	sportMonks "inplay-football-service/api"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest() {
	// Fetch current fixture from external api
	res, err := sportMonks.GetTodaysFixtures()
	if err != nil {
		fmt.Println("Error in GetTodaysFixtures:", err)
		return
	}
	fmt.Printf("Results: %v\n", res)
	return
}

func main() {
	lambda.Start(handleRequest)
}
