package ui

import (
	"context"

	"github.com/bahner/go-ma-actor/entity/actor"
)

func (ui *ChatUI) startActor() {

	ui.currentActorCtx, ui.currentActorCancel = context.WithCancel(context.Background())

	// Now that the UI is created, we can start the actor and subscribe to its events.
	go ui.a.Subscribe(ui.currentActorCtx, ui.a.Entity)

	// We want to handle envelopes for the actor, then deliver the messages
	// to the UI from the incoming envelopes.
	go ui.a.HandleIncomingEnvelopes(ui.currentActorCtx, ui.chMessages)
	go ui.a.Entity.HandleIncomingMessages(ui.currentActorCtx, ui.chMessages)

	go actor.HelloWorld(ui.currentActorCtx, ui.a) // This wait a bit before sending the message.
}
