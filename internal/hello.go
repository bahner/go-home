package internal

import (
	"context"

	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma/msg"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
)

func HelloWorld(ctx context.Context, a *actor.Actor, b *p2ppubsub.Topic) {

	if a == nil {
		return
	}

	if b == nil {
		return
	}

	if a.Entity == nil {
		return
	}

	me := a.Entity.DID.Id
	greeting := []byte("Hello, world! " + me + " is here.")

	if b != nil {
		mesg, _ := msg.NewBroadcast(me, greeting, "text/plain", a.Keyset.SigningKey.PrivKey)
		mesg.Broadcast(ctx, b)
	}
}
