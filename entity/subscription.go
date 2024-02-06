package entity

import (
	"context"
	"fmt"

	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

type Subscription struct {
	Cancel   context.CancelFunc
	Messages chan *p2ppubsub.Message
}

// Subscribe to the entity's topic
// Returns a subscription which contains the cancel function and a channel for messages
// The entity's own Subscription field is set to the returned subscription
func (e *Entity) Subscribe() (*Subscription, error) {
	ctx, cancel := context.WithCancel(context.Background())

	sub, err := e.Topic.Subscribe()
	if err != nil {
		cancel() // Ensure resources are cleaned up in case of early return
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	e.Subscription = &Subscription{
		Cancel:   cancel,
		Messages: make(chan *p2ppubsub.Message),
	}

	go func() {
		defer sub.Cancel()
		defer close(e.Subscription.Messages)
		for {
			select {
			case <-ctx.Done():
				log.Infof("entity/subscribe: Context done. Closing the topic subscription")
				return
			default:
				message, err := sub.Next(ctx)
				if err != nil {
					log.Errorf("entity/subscribe: error getting next message: %v", err)
					continue
				}
				e.Subscription.Messages <- message
			}
		}
	}()

	return e.Subscription, nil
}