package ui

func (ui *ChatUI) handleSetCommand(args []string) {

	if len(args) == 3 {
		switch args[1] {
		case "broadcast":
			ui.handleSetBroadcastCommand(args)
		}
	} else {
		ui.handleHelpSetCommand(args)
	}

}

func (ui *ChatUI) handleHelpSetCommand(args []string) {
	ui.displaySystemMessage("Usage: /set broadcast on|off")
	ui.displaySystemMessage("For now toggles broadcast messages on and off")
}
