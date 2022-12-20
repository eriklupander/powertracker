package main

import (
	"encoding/json"
	"github.com/hasura/go-graphql-client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

const tibberGQLSubscriptionUrl = "wss://api.tibber.com/v1-beta/gql/subscriptions"

func recordPowerUsageFromWatty(accessToken, homeId string) error {

	subscriptionClient := graphql.NewSubscriptionClient(tibberGQLSubscriptionUrl).
		WithConnectionParams(map[string]interface{}{
			"token": accessToken,
		})

	defer subscriptionClient.Close()

	// GraphQL variable
	variables := map[string]interface{}{
		"homeId": graphql.ID(homeId),
	}

	// Channel to pass data from subscription callback to "main" goroutine
	dataChan := make(chan *subscription)

	// Subscribe to real-time power usage
	id, err := subscriptionClient.Subscribe(&subscription{}, variables, func(dataValue []byte, errValue error) error {
		if errValue != nil {
			return errValue
		}
		if dataValue == nil {
			return errors.New("got nil data")
		}
		m := &subscription{}
		if err := json.Unmarshal(dataValue, m); err != nil {
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
		err = subscriptionClient.Run()
		if err != nil {
			logrus.WithError(err).Error("error calling Run()")
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
	if err := subscriptionClient.Unsubscribe(id); err != nil {
		logrus.WithError(err).Error("error occurred trying to unsubscribe from subscription")
	}
	return nil
}

// subscription forms the root of our GraphQL query having a homeId parameter.
type subscription struct {
	LiveMeasurement liveMeasurement `graphql:"liveMeasurement(homeId: $homeId)"`
}

// liveMeasurement forms the timestamp + accumulated usage part of the GraphQL query
type liveMeasurement struct {
	Timestamp              graphql.String `graphql:"timestamp"`
	AccumulatedConsumption graphql.Float  `graphql:"accumulatedConsumption"`
}
