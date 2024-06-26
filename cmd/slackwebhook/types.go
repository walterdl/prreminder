package main

type BaseSlackEvent struct {
	// Type is "event_callback" for message events or "url_verification" for the auth event.
	Type string `json:"type"`
}

type AuthSlackEvent struct {
	// Type is "url_verification" for this type of event.
	Type      string `json:"type"`
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
}

type MessageEvent struct {
	// Type is "event_callback" for this type of event.
	Type  string             `json:"type"`
	Event MessageEventDetail `json:"event"`
}

type MessageEventDetail struct {
	// Type is "message" for message-related events.
	Type string `json:"type"`
	// Ts is always present.
	Ts string `json:"ts"`
	// Channel is always present.
	Channel string `json:"channel"`
	// Text is the message content. Present when is a root or a reply message.
	Text string `json:"text"`
	// ThreadTS is present if the message is a reply to another message.
	ThreadTS string `json:"thread_ts"`
	// Subtype is "message_changed" when a message is edited.
	// It is also "message_replied" when the message is a reply.
	// However, due to a Slack bug, this subtype value is not delivered.
	// ThreadTS must be used to determine if the message is a reply.
	Subtype string        `json:"subtype"`
	Message EditedMessage `json:"message"`
}

type EditedMessage struct {
	// Text is the new message content.
	Text string `json:"text"`
	// Subtype is "message" when is a new reply.
	// It is "thread_broadcast" when a reply is broadcasted.
	Subtype string `json:"subtype"`
	// ThreadTS is always present.
	ThreadTS string `json:"thread_ts"`
	// Ts is always present.
	Ts string `json:"ts"`
}
