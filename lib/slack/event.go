package slack

func IsRootMessage(msg BaseSlackMessageEvent) bool {
	return isMessageEvent(msg) && msg.Event.ThreadTS == "" && msg.Event.Subtype == "" && msg.Event.Text != ""
}

func IsRootMessageEdition(msg BaseSlackMessageEvent) bool {
	if !isMessageEvent(msg) {
		return false
	}

	if msg.Event.Subtype == "message_changed" && msg.Event.Message.ThreadTS == "" && msg.Event.Message.Text != "" {
		// It's a root message edition
		return true
	}

	if msg.Event.Subtype == "message_changed" && msg.Event.Message.Ts == msg.Event.Message.ThreadTS && msg.Event.Message.Text != "" {
		// It's a root message edition which has replies.
		return true
	}

	return false
}

func isMessageEvent(msg BaseSlackMessageEvent) bool {
	return msg.Type == "event_callback" && msg.Event.Type == "message"
}
