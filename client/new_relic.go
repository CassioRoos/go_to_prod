package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

const ContextTransactionKey = "NewRelicTransaction"

func prepareNewRelicExternalSegment(ctx echo.Context, request *http.Request) *newrelic.ExternalSegment {
	if ctx == nil {
		return &newrelic.ExternalSegment{}
	}
	txn, _ := ctx.Get(ContextTransactionKey).(*newrelic.Transaction)
	if txn != nil {
		// Copy struct request to remove sensitive data
		requestCopy := *request
		// Set body empty to send to New Relic
		requestCopy.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
		// Set headers empty to send to New Relic
		requestCopy.Header = http.Header{}
		// Remove GLBID from URL
		if strings.Contains(requestCopy.URL.String(), "/glbid/") {
			requestCopy.URL, _ = url.ParseRequestURI(strings.Split(requestCopy.URL.String(), "/glbid/")[0] + "/glbid/GLBID")
		}

		return newrelic.StartExternalSegment(txn, &requestCopy)
	}
	return &newrelic.ExternalSegment{}
}

func endNewRelicExternalSegment(segment *newrelic.ExternalSegment, response *http.Response) {
	// Copy struct response to remove sensitive data
	responseCopy := http.Response{}
	if response != nil {
		responseCopy = *response
	}
	// Set body empty to send to New Relic
	responseCopy.Body = ioutil.NopCloser(bytes.NewReader([]byte{}))
	// Set headers empty to send to New Relic
	responseCopy.Header = http.Header{}

	segment.Response = &responseCopy
	segment.End()
}

func prepareNewRelicMessageProducerSegment(ctx echo.Context, destinationType newrelic.MessageDestinationType, destinationName string) *newrelic.MessageProducerSegment {
	if ctx == nil {
		return &newrelic.MessageProducerSegment{}
	}
	txn, _ := ctx.Get(ContextTransactionKey).(*newrelic.Transaction)
	if txn != nil {
		s := newrelic.MessageProducerSegment{
			Library:              "ActiveMQ",
			DestinationType:      destinationType,
			DestinationName:      destinationName,
			DestinationTemporary: false,
		}
		s.StartTime = newrelic.StartSegmentNow(txn)
		return &s
	}
	return &newrelic.MessageProducerSegment{}
}
