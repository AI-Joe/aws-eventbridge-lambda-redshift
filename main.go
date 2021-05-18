package main

import (
	"context"

	"github.com/AI-Joe/aws-eventbridge-lambda-redshift/handlers"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(ctx context.Context, event handlers.RedshiftEvent) error {
	return event.ExecuteStatements(ctx)
}

func main() {
	lambda.Start(Handler)
}
