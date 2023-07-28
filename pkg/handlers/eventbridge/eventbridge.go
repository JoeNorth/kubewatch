package eventbridge

import (
	"context"
	"encoding/json"
	"os"
	"time"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	eb "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eb_types "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/bitnami-labs/kubewatch/pkg/event"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

var eventbridgeWarnMsg = `
%s

The option arguments for the EventBridge endpointId and eventBusName
and EKS cluster ARN can be set using the following flags or environment variables:

--endpointId/-e
--clusterArn/-c
--eventBusName/-b

export KW_EVENTBRIDGE_ENDPOINT_ID=endpointId
export KW_EVENTBRIDGE_CLUSTER_ARN=clusterArn
export KW_EVENTBRIDGE_EVENT_BUS_NAME=eventBusName

Command line flags will override environment variables

`

// EventBridge handler implements the handler.Handler interface
type EventBridge struct {
	EndpointId   string
	ClusterArn   string
	EventBusName string
}

type EventBridgeEntryDetail struct {
	Message    string         `json:"message,omitempty"`
	Operation  string         `json:"operation,omitempty"`
	ClusterId  string         `json:"clusterId,omitempty"`
	Namespace  string         `json:"namespace,omitempty"`
	Kind       string         `json:"kind,omitempty"`
	ApiVersion string         `json:"apiVersion,omitempty"`
	Component  string         `json:"component,omitempty"`
	Host       string         `json:"host,omitempty"`
	Reason     string         `json:"reason,omitempty"`
	Status     string         `json:"status,omitempty"`
	Name       string         `json:"name,omitempty"`
	Obj        runtime.Object `json:"obj,omitempty"`
	OldObj     runtime.Object `json:"oldObj,omitempty"`
}

func (m *EventBridge) Init(c *config.Config) error {

	m.EndpointId = c.Handler.EventBridge.EndpointId
	m.ClusterArn = c.Handler.EventBridge.ClusterArn
	m.EventBusName = c.Handler.EventBridge.EventBusName

	if m.EndpointId == "" {
		m.EndpointId = os.Getenv("KW_EVENTBRIDGE_ENDPOINT_ID")
	}
	if m.EndpointId == "" {
		logrus.Warnf(eventbridgeWarnMsg, "Missing EventBridge endpointId, using default endpoint.")
	}
	if m.ClusterArn == "" {
		m.ClusterArn = os.Getenv("KW_EVENTBRIDGE_CLUSTER_ARN")
	}
	if m.ClusterArn == "" {
		logrus.Warnf(eventbridgeWarnMsg, "Missing EKS Cluster ARN. Events will not include cluster information.")
	}
	if m.EventBusName == "" {
		m.EventBusName = os.Getenv("KW_EVENTBRIDGE_EVENT_BUS_NAME")
	}
	if m.EventBusName == "" {
		logrus.Warnf(eventbridgeWarnMsg, "Missing EventBridge event bus name, using default event bus.")
	}

	// TODO Validate ClusterArn is valid
	// TODO Validate EventBusName is valid
	return nil
}

func (m *EventBridge) Handle(e event.Event) {

	eventEntry, err := preparePutEventsEntry(e, m)
	if err != nil {
		logrus.Errorf("Failed to marshal EventBridgeEntryDetail: %v", err)
		return
	}

	//TODO Calculate the byte size of the PutEventsRequestEntry to ensure it is not over 256KB
	// https://docs.aws.amazon.com/eventbridge/latest/userguide/eb-putevent-size.html

	ctx := context.TODO()

	logrus.Info("Loading AWS config")
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		logrus.Fatalf("Failed to load AWS config: %v", err)
		return
	}

	eb_client := eb.NewFromConfig(cfg)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	entryInput := eb.PutEventsInput{
		Entries: []eb_types.PutEventsRequestEntry{*eventEntry},
	}

	if m.EndpointId != "" {
		entryInput.EndpointId = &m.EndpointId
	}

	_, err = eb_client.PutEvents(ctx, &entryInput)

	if err != nil {
		logrus.Errorf("Failed to send event to EventBridge: %v", err)
		return
	}

	// TODO Add logging for response FailedEntryCount and Entries.
	// Specifically, the EventID should be logged as INFO.

}

func preparePutEventsEntry(e event.Event, m *EventBridge) (*eb_types.PutEventsRequestEntry, error) {

	details := EventBridgeEntryDetail{
		Message:    e.Message(),
		Operation:  m.formatReason(e),
		ClusterId:  m.ClusterArn,
		Namespace:  e.Namespace,
		Kind:       e.Kind,
		ApiVersion: e.ApiVersion,
		Component:  e.Component,
		Host:       e.Host,
		Reason:     e.Reason,
		Status:     e.Status,
		Name:       e.Name,
		Obj:        e.Obj,
		OldObj:     e.OldObj,
	}

	eventEntry, err := createEventEntry(e, m, details)

	return eventEntry, err
}

func (m *EventBridge) formatReason(e event.Event) string {
	switch e.Reason {
	case "Created":
		return "create"
	case "Updated":
		return "update"
	case "Deleted":
		return "delete"
	default:
		return "unknown"
	}
}

// Creates a PutEventsRequestEntry from an event
func createEventEntry(e event.Event, m *EventBridge, d EventBridgeEntryDetail) (*eb_types.PutEventsRequestEntry, error) {

	var detailType string = "Kubewatch Event"
	var source string

	if m.ClusterArn != "" {
		source = "kubewatch/" + m.ClusterArn
	} else {
		source = "kubewatch"
	}

	// Convert EventBridgeEntryDetail to JSON
	b, err := json.Marshal(d)
	if err != nil {
		return &eb_types.PutEventsRequestEntry{}, err
	}

	detail := string(b)

	eventEntry := eb_types.PutEventsRequestEntry{
		Detail:     &detail,
		DetailType: &detailType,
		Resources: []string{
			m.ClusterArn,
		},
		Source: &source,
	}

	if m.EventBusName != "" {
		eventEntry.EventBusName = &m.EventBusName
	}

	return &eventEntry, nil
}

//Validates that EndpointId is valid
// func validateEndpointId(c *config.Config, endpointId string) error {
// 	Regex to match ^[A-Za-z0-9\-]+[\.][A-Za-z0-9\-]+$
// }
// Validates that ClusterArn is valid
// func validateClusterArn(c *config.Config, clusterArn string) error {
// Regex to match (arn:aws[\w-]*:eks:[a-z]{2}-[a-z]+-[\w-]+:[0-9]{12}:cluster\/)?[0-9A-Za-z][A-Za-z0-9\-_]*
// }

// Validates that EventBusName is valid
// func validateEventBusName(c *config.Config, eventBusName string) error {
// 	Regex to match (arn:aws[\w-]*:events:[a-z]{2}-[a-z]+-[\w-]+:[0-9]{12}:event-bus\/)?[\.\-_A-Za-z0-9]+
// }
