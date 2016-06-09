package slack

import (
	"testing"

	"github.com/blendlabs/go-assert"
)

func TestClientAddEventListener(t *testing.T) {
	assert := assert.New(t)
	c := NewClient(UUIDv4().ToShortString())
	c.AddEventListener(EventBotAdded, func(c *Client, m *Message) {})
	assert.NotEmpty(c.EventListeners[EventBotAdded])
}

func TestClientRemoveEventListener(t *testing.T) {
	assert := assert.New(t)
	c := NewClient(UUIDv4().ToShortString())
	c.AddEventListener(EventBotAdded, func(c *Client, m *Message) {})
	assert.NotEmpty(c.EventListeners[EventBotAdded])

	c.RemoveEventListeners(EventBotAdded)
	assert.Empty(c.EventListeners[EventBotAdded])
}
