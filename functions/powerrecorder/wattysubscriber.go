package main

import (
	"encoding/json"
	"fmt"
	"github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
	"time"
)

const tibberGQLSubscriptionUrl = "wss://api.tibber.com/v1-beta/gql/subscriptions"

func recordPowerUsageFromWatty(accessToken, homeId string) error {

	sclient := graphql.NewSubscriptionClient(tibberGQLSubscriptionUrl).
		WithConnectionParams(map[string]interface{}{
			"token": accessToken,
		})

	defer sclient.Close()

	variables := map[string]interface{}{
		"homeId": graphql.ID(homeId),
	}

	dataChan := make(chan *subscription)

	// Subscribe to real-time power usage
	id, err := sclient.Subscribe(&subscription{}, variables, func(dataValue *json.RawMessage, errValue error) error {
		if errValue != nil {
			return errValue
		}
		if dataValue == nil {
			return errors.New("got nil data")
		}
		m := &subscription{}
		if err := json.Unmarshal(*dataValue, m); err != nil {
			return errors.Wrap(err, "unmarshalling measurement")
		}

		// pass data to channel
		dataChan <- m
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "starting subscription")
	}

	// finally run the subscription in a goroutine. If start fails, we'll pass nil to the dataChan.
	go func() {
		err = sclient.Run()
		if err != nil {
			fmt.Println("Error calling Run(): " + err.Error())
			dataChan <- nil // pass nil in order to cancel select below
		}
	}()

	// block here until we have data. Once we get data or time out, unsubscribe and exit.
	select {
	case sub := <-dataChan:
		if sub != nil {
			ingest(record{HomeId: homeId, AccumulatedConsumption: float64(sub.LiveMeasurement.AccumulatedConsumption)})
		}
	case <-time.NewTimer(time.Second * 10).C:

	}
	if err := sclient.Unsubscribe(id); err != nil {
		fmt.Println("an error occurred trying to unsubscribe from subscription: " + err.Error())
	}
	return nil
}

type subscription struct {
	LiveMeasurement liveMeasurement `graphql:"liveMeasurement(homeId: $homeId)"`
}
type liveMeasurement struct {
	Timestamp              graphql.String `graphql:"timestamp"`
	AccumulatedConsumption graphql.Float  `graphql:"accumulatedConsumption"`
}
