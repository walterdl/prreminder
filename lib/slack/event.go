package slack

func IsRootMessage(msg BaseSlackMessageEvent) bool {
	isCorrectSubtype := msg.Event.Subtype == "" || msg.Event.Subtype == "file_share"
	return isMessageEvent(msg) &&
		msg.Event.ThreadTS == "" &&
		isCorrectSubtype &&
		msg.Event.Text != ""
}

func IsRootMessageEdition(msg BaseSlackMessageEvent) bool {
	if !isMessageEvent(msg) {
		return false
	}

	if editionTS(msg) == "" {
		return false
	}

	return true
}

// editionTS returns the TS of the root message edited.
// If the message is not an edition, it returns an empty string
func editionTS(msg BaseSlackMessageEvent) string {
	if msg.Event.Subtype == "message_changed" &&
		msg.Event.Message.ThreadTS == "" &&
		msg.Event.Message.Text != "" {
		// It's a root message edition
		return msg.Event.Message.TS
	}

	if msg.Event.Subtype == "message_changed" &&
		msg.Event.Message.TS == msg.Event.Message.ThreadTS &&
		msg.Event.Message.Text != "" &&
		msg.Event.Message.Text != msg.Event.PreviousMessage.Text {
		// It's a root message edition which has replies.
		return msg.Event.Message.TS
	}

	return ""
}

func isMessageEvent(msg BaseSlackMessageEvent) bool {
	return msg.Type == "event_callback" && msg.Event.Type == "message"
}

func Channel(msg BaseSlackMessageEvent) string {
	return msg.Event.Channel
}

// TS returns the TS of the root message from a edition or deletion event.
func TS(msg BaseSlackMessageEvent) string {
	if IsRootMessageEdition(msg) {
		return editionTS(msg)
	}

	if IsRootDeletion(msg) {
		return msg.Event.PreviousMessage.TS
	}

	return msg.Event.TS
}

func IsRootDeletion(msg BaseSlackMessageEvent) bool {
	if msg.Event.Type == "message" &&
		msg.Event.Subtype == "message_deleted" &&
		msg.Event.ThreadTS == "" {
		return true
	}

	return false
}
