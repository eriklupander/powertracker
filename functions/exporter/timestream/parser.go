package timestream

import "github.com/aws/aws-sdk-go/service/timestreamquery"

func processScalarType(data *timestreamquery.Datum) string {
	return *data.ScalarValue
}
