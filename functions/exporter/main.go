package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/eriklupander/powertracker/functions/exporter/timestream"
)

var chiLambda *chiadapter.ChiLambda

// handler is invoked whenever this lambda executes.
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return chiLambda.ProxyWithContext(ctx, req)
}

// main is called when a new lambda is constructed, when this happens is up to the underlying AWS machinery so
// don't rely on it happening on every invocation.
func main() {
	chiLambda = chiadapter.New(setupRouter(timestream.NewDataSource()))
	lambda.StartWithContext(context.Background(), handler)
}
