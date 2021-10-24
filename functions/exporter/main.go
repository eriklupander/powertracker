package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/eriklupander/powertracker/functions/exporter/timestream"
)

var chiLambda *chiadapter.ChiLambda

// handler is invoked whenever this lambda executes.
func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := chiLambda.ProxyWithContext(ctx, req)

	fmt.Printf("Type of resp: %T\n", resp)


	// hack: copy single valued headers to multi-val
	for k, v := range resp.MultiValueHeaders {
		fmt.Printf("headers: %v %v\n", k, v)
		if _, ok := resp.MultiValueHeaders[k]; !ok {
			resp.Headers[k] = v[0]
		}
	}
	return resp, err
}

// main is called when a new lambda is constructed, when this happens is up to the underlying AWS machinery so
// don't rely on it happening on every invocation.
func main() {
	chiLambda = chiadapter.New(setupRouter(timestream.NewDataSource()))
	lambda.StartWithContext(context.Background(), handler)
}
