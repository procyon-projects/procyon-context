package context

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type testContext struct {
	ApplicationContext
}

func (listener *testContext) GetContextId() ContextId {
	return ""
}

func (listener *testContext) Get(key string) interface{} {
	return nil
}

func (listener *testContext) Put(key string, value interface{}) {

}

var testEventId1 = GetEventId("testEventId1")
var testEventId2 = GetEventId("testEventId2")

func getTestEventId1() ApplicationEventId {
	return testEventId1
}

func getTestEventId2() ApplicationEventId {
	return testEventId2
}

type testEvent1 struct {
}

func (event testEvent1) GetEventId() ApplicationEventId {
	return testEventId1
}

func (event testEvent1) GetParentEventId() ApplicationEventId {
	return 0
}

func (event testEvent1) GetSource() interface{} {
	return nil
}

func (event testEvent1) GetTimestamp() int64 {
	return 0
}

type testEvent2 struct {
}

func (event testEvent2) GetEventId() ApplicationEventId {
	return testEventId1
}

func (event testEvent2) GetParentEventId() ApplicationEventId {
	return 0
}

func (event testEvent2) GetSource() interface{} {
	return nil
}

func (event testEvent2) GetTimestamp() int64 {
	return 0
}

type testApplicationListener struct {
	mock.Mock
}

func (listener testApplicationListener) GetApplicationListenerName() string {
	return "github.com.procyon.projects.testApplicationListener"
}

func (listener testApplicationListener) SubscribeEvents() []ApplicationEventId {
	return []ApplicationEventId{
		getTestEventId1(),
		getTestEventId2(),
	}
}

func (listener testApplicationListener) OnApplicationEvent(context Context, event ApplicationEvent) {
	listener.Called(context, event)
}

func TestSimpleApplicationEventBroadcaster_RegisterApplicationListener(t *testing.T) {
	broadcaster := NewSimpleApplicationEventBroadcaster()
	broadcaster.RegisterApplicationListener(testApplicationListener{})
	assert.Equal(t, 1, len(broadcaster.eventListenerMap[getTestEventId1()]))
	assert.Equal(t, 1, len(broadcaster.eventListenerMap[getTestEventId2()]))
}

func TestSimpleApplicationEventBroadcaster_RemoveAllApplicationListeners(t *testing.T) {
	broadcaster := NewSimpleApplicationEventBroadcaster()
	broadcaster.RegisterApplicationListener(testApplicationListener{})
	broadcaster.RemoveAllApplicationListeners()
	assert.Equal(t, 0, len(broadcaster.eventListenerMap))
}

func TestSimpleApplicationEventBroadcaster_BroadcastEvent(t *testing.T) {
	context := &testContext{}
	broadcaster := NewSimpleApplicationEventBroadcaster()
	listener := testApplicationListener{}

	testEvent1 := testEvent1{}
	listener.On("OnApplicationEvent", context, testEvent1)
	testEvent2 := testEvent2{}
	listener.On("OnApplicationEvent", context, testEvent2)

	broadcaster.RegisterApplicationListener(listener)

	broadcaster.BroadcastEvent(context, testEvent1)
	broadcaster.BroadcastEvent(context, testEvent2)

	listener.AssertExpectations(t)
}

func TestSimpleApplicationEventBroadcaster_UnregisterApplicationListener(t *testing.T) {
	broadcaster := NewSimpleApplicationEventBroadcaster()
	listener := testApplicationListener{}
	broadcaster.RegisterApplicationListener(listener)
	assert.Equal(t, 2, len(broadcaster.eventListenerMap))
	broadcaster.UnregisterApplicationListener(listener)
	assert.Equal(t, 0, len(broadcaster.eventListenerMap))
}
