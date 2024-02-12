package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma/msg"

	log "github.com/sirupsen/logrus"
)

const (
	defaultMsg               = "yo"
	defaultBroadcast         = "Hello, world!"
	pubsubMessagesBuffersize = 32
)

func init() {
	pflag.String("msg", defaultMsg, "Message to send as a pong. For fun and identification.")
	viper.BindPFlag("pong.msg", pflag.Lookup("msg"))
	viper.SetDefault("pong.msg", defaultMsg)

}

func main() {

	pflag.Parse()

	config.Init("pong")
	config.InitLogging()
	config.InitP2P()
	config.InitActor()

	p, err := p2p.Init(nil)
	if err != nil {
		log.Errorf("Error initializing p2p node: %v", err)
		os.Exit(69) // EX_UNAVAILABLE
	}

	if err != nil {
		log.Errorf("Error initializing p2p node: %v", err)
		os.Exit(69) // EX_UNAVAILABLE
	}

	// We need to discover peers before we can do anything else.
	p.DiscoverPeers()

	k := config.GetKeyset()
	a, err := actor.NewFromKeyset(k)
	if err != nil {
		log.Errorf("Error initializing actor: %v", err)
		os.Exit(70) // EX_SOFTWARE
	}

	me := a.Entity.DID.String()

	fmt.Printf("I am : %s\n", me)
	fmt.Printf("My public key is: %s\n", p.Node.ID().String())

	// Now we can start continuous discovery in the background.
	ctx, cancel := config.GetDiscoveryContext()
	defer cancel()

	go p.DiscoveryLoop(ctx)

	ctxBackground := context.Background()
	// Just me the entity here
	go a.Subscribe(ctxBackground, a.Entity)
	go handleEnvelopeEvents(ctxBackground, a)
	go handleMessageEvents(ctxBackground, a)

	b, err := msg.NewBroadcast(me, me, []byte(defaultBroadcast), "text/plain", k.SigningKey.PrivKey)
	if err != nil {
		log.Fatalf("Error creating broadcast: %v", err)
	}

	err = b.Broadcast(ctx, a.Entity.Topic)
	if err != nil {
		log.Fatalf("Error sending broadcast: %v", err)
	}

	// This is defined in web.go. It makes it possible to add extra parameters to the handler.
	h := &entity.WebHandlerData{
		P2P:    p,
		Entity: a.Entity,
	}
	http.HandleFunc("/", h.WebHandler)
	log.Infof("Listening on %s", config.GetHttpSocket())
	err = http.ListenAndServe(config.GetHttpSocket(), nil)
	if err != nil {
		log.Fatal(err)
	}
}
