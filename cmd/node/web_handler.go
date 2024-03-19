package main

import (
	"fmt"
	"net/http"

	"github.com/bahner/go-ma"
	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// Assuming you have initialized variables like `h` and `rendezvous` somewhere in your main function or globally

type WebEntity struct {
	P2P    *p2p.P2P
	Entity *entity.Entity
}

type WebEntityDocument struct {
	Title               string
	H1                  string
	H2                  string
	Addrs               []multiaddr.Multiaddr
	PeersWithSameRendez peer.IDSlice
	AllConnectedPeers   peer.IDSlice
	Topics              []string
}

func NewWebEntityDocument() *WebEntityDocument {
	return &WebEntityDocument{}
}

func (data *WebEntity) WebHandler(w http.ResponseWriter, r *http.Request) {
	webHandler(w, r, data.P2P, data.Entity)
}

func webHandler(w http.ResponseWriter, _ *http.Request, p *p2p.P2P, e *entity.Entity) {

	doc := NewWebEntityDocument()

	titleStr := fmt.Sprintf("Entity: %s", e.DID.Id)
	h1str := titleStr
	doc.Title = titleStr
	doc.H1 = h1str
	doc.H2 = fmt.Sprintf("%s@%s", ma.RENDEZVOUS, (p.Host.ID().String()))
	doc.Addrs = p.Host.Addrs()
	doc.AllConnectedPeers = p.AllConnectedPeers()
	doc.PeersWithSameRendez = p.ConnectedProtectedPeers()
	doc.Topics = p.PubSub.GetTopics()

	fmt.Fprint(w, doc.String())
}

func (d *WebEntityDocument) String() string {

	html := "<!DOCTYPE html>\n<html>\n<head>\n"
	html += "<style>table, th, td {border: 1px solid black;}</style>"
	if d.Title != "" {
		html += "<title>" + d.Title + "</title>\n"
	}
	html += fmt.Sprintf(`<meta http-equiv="refresh" content="%d">`, config.HttpRefresh())
	html += "</head>\n<body>\n"
	if d.H1 != "" {
		html += "<h1>" + d.H1 + "</h1>\n"
	}
	html += "<hr>"
	if d.H2 != "" {
		html += "<h2>" + d.H2 + "</h2>\n"
	}

	// Subscribed topics
	if len(d.Topics) > 0 {
		html += fmt.Sprintf("<h2>Topics (%d):</h2>\n", len(d.Topics))
		html += UnorderedListFromTopicsSlice(d.Topics)
	}

	// Peers with Same Rendezvous
	if len(d.PeersWithSameRendez) > 0 {
		html += fmt.Sprintf("<h2>Discovered peers (%d):</h2>\n", len(d.PeersWithSameRendez))
		html += UnorderedListFromPeerIDSlice(d.PeersWithSameRendez)
	}

	// Info leak? Not really important anyways.
	// // Addresses
	if len(d.Addrs) > 0 {
		html += "<h2>Addresses</h2>\n"
		html += "<table>\n"
		for _, addr := range d.Addrs {
			html += "<tr><td>" + addr.String() + "</td></tr>\n"
		}
		html += "</table>\n"
	}

	// All Connected Peers
	if len(d.AllConnectedPeers) > 0 {
		html += fmt.Sprintf("<h2>libp2p Network Peers (%d):</h2>\n", len(d.AllConnectedPeers))
		html += UnorderedListFromPeerIDSlice(d.AllConnectedPeers)
	}

	html += "</body>\n</html>"
	return html
}
