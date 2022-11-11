package services

import (
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestGetTemplater_AwsSqs(t *testing.T) {
	n := Notification{
		Message: "{{.message}}",
		AwsSqs: &AwsSqsNotification{
			MessageAttributes: map[string]string{
				"attributeKey": "{{.messageAttributeValue}}",
			},
		},
	}

	templater, err := n.GetTemplater("", template.FuncMap{})
	if !assert.NoError(t, err) {
		return
	}

	var notification Notification

	err = templater(&notification, map[string]interface{}{
		"message":               "abcdef",
		"messageAttributeValue": "123456",
	})

	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "abcdef", notification.Message)
	assert.Equal(t, map[string]string{
		"attributeKey": "123456",
	}, notification.AwsSqs.MessageAttributes)
}

func TestSendLocalEndpoint_AwsSqs(t *testing.T) {
	opt := AwsSqsOptions{
		EndpointUrl: "http://127.0.0.1",
	}

	notification := Notification{
		AwsSqs: &AwsSqsNotification{},
	}

	s := NewAwsSqsService(opt)
	if err := s.Send(notification, Destination{}); err != nil {
		assert.EqualError(t, err, "operation error SQS: GetQueueUrl, failed to sign request: failed to retrieve credentials: failed to refresh cached credentials, no EC2 IMDS role found, operation error ec2imds: GetMetadata, request canceled, context deadline exceeded")
	}
}

func TestSend_WithoutCredential_AwsSqs(t *testing.T) {
	opt := AwsSqsOptions{}
	notification := Notification{
		AwsSqs: &AwsSqsNotification{},
	}

	s := NewAwsSqsService(opt)
	if err := s.Send(notification, Destination{}); err != nil {
		assert.EqualError(t, err, "operation error SQS: GetQueueUrl, failed to resolve service endpoint, an AWS region is required, but was not found")
	}
}
