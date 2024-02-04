package ui

import (
	"fmt"

	"github.com/bahner/go-ma-actor/alias"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) handleEnterCommand(args []string) {

	if len(args) == 2 {

		_did := args[1]
		// If id is not a valid did, then try to find it in the aliases
		if !did.IsValidDID(_did) {
			_did = alias.LookupEntityNick(_did)
		}

		// If it is still not a valid did, then return
		if _did == "" {
			ui.displaySystemMessage("Invalid DID")
			return
		}

		log.Debugf("Trying to find: %s", _did)

		// If the DID is our own identity that is already handled.
		if _did == ui.a.DID.String() {
			ui.displaySystemMessage("You can't enter yourself.")
			return
		}

		// If this is not the same as the last known location, then
		// update the last known location
		if ui.e.DID.String() == _did {
			ui.displaySystemMessage("You are already here.")
			return
		}

		// Update the UI
		err := ui.setEntity(_did)
		if err != nil {
			ui.displaySystemMessage("Error changing entity: " + err.Error())
			return
		}
		ui.msgBox.SetTitle(ui.e.Nick)
		ui.displaySystemMessage("Entered: " + _did)

		// Update the location
		err = ui.a.UpdateLastKnowLocation(_did)
		if err != nil {
			ui.displaySystemMessage("Error updating last known location: " + err.Error())
			return
		}

	} else {
		ui.displaySystemMessage("Usage: /enter <DID>")
	}
}

func (ui *ChatUI) setEntity(did string) error {

	var err error

	log.Debugf("Creating entity for topic %s", did)
	// e, err = getOrCreateEntity(did)
	e, err := entity.GetOrCreate(did)
	if err != nil {
		return fmt.Errorf("error getting or creating entity: %w", err)
	}

	// Loog up the nick for the entity
	e.Nick = alias.LookupEntityDID(did)

	// Now pivot to the new entity
	// and cancel the old.
	old_nick := ui.e
	ui.e = e
	old_nick.Subscription.Cancel()

	log.Infof("Location changed to %s", ui.e.Topic.String())

	// Start handling the new topic
	go ui.subscribeEntityMessages(e)

	return nil

}
func (ui *ChatUI) setActor(a *entity.Entity) error {

	// Loog up the nick for the entity
	a.Nick = alias.LookupEntityDID(a.DID.String())

	// Now pivot to the new entity
	// and cancel the old.
	old_actor := ui.a
	ui.a = a
	old_actor.Subscription.Cancel()

	log.Infof("Actor changed to %s", ui.a.Topic.String())

	// Start handling the new topic
	go ui.subscribeEntityMessages(a)
	go ui.handleIncomingEnvelopes()

	return nil

}
