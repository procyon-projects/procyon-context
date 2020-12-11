package context

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testApplicationContextEvent(t *testing.T, event ApplicationContextEvent, eventId ApplicationEventId, parentEventId ApplicationEventId) {
	assert.Equal(t, eventId, event.GetEventId())
	assert.Equal(t, parentEventId, event.GetParentEventId())
	assert.NotEqual(t, 0, event.GetTimestamp())
	assert.NotNil(t, event.GetSource())
	assert.NotNil(t, event.GetApplicationContext())
}

func TestApplicationContextStartedEvent(t *testing.T) {
	context := &testContext{}
	event := NewApplicationContextStartedEvent(context)
	testApplicationContextEvent(t, event, ApplicationContextStartedEventId(), ApplicationContextEventId())
}

func TestApplicationContextRefreshedEvent(t *testing.T) {
	context := &testContext{}
	event := NewApplicationContextRefreshedEvent(context)
	testApplicationContextEvent(t, event, ApplicationContextRefreshedEventId(), ApplicationContextEventId())
}

func TestApplicationContextStoppedEvent(t *testing.T) {
	context := &testContext{}
	event := NewApplicationContextStoppedEvent(context)
	testApplicationContextEvent(t, event, ApplicationContextStoppedEventId(), ApplicationContextEventId())
}

func TestApplicationContextClosedEvent(t *testing.T) {
	context := &testContext{}
	event := NewApplicationContextClosedEvent(context)
	testApplicationContextEvent(t, event, ApplicationContextClosedEventId(), ApplicationContextEventId())
}
