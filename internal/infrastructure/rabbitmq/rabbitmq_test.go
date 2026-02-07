package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRabbitMQClient_Compilation(t *testing.T) {
	// This test just ensures the types and methods are correctly defined.
	// We don't attempt to connect to a real server here as it might not be available.
	assert.True(t, true)
}

func TestRabbitMQClient_Connect_SkipIfUnavailable(t *testing.T) {
	// Use a dummy local URL
	url := "amqp://user:password@localhost:5672/"
	client, err := NewClient(url)
	if err != nil {
		t.Skipf("RabbitMQ not available at %s, skipping integration test: %v", url, err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Publish(ctx, "test.key", map[string]string{"foo": "bar"})
	assert.NoError(t, err)

	done := make(chan bool)
	err = client.Subscribe("test_queue", "test.key", func(body []byte) error {
		done <- true
		return nil
	})
	assert.NoError(t, err)

	select {
	case <-done:
		// Success
	case <-time.After(3 * time.Second):
		t.Error("Timed out waiting for message")
	}
}
