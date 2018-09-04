package quiz

import (
	"fmt"
	"log"
	"strings"

	"github.com/nlopes/slack"
)

type SlackListener struct {
	Client    *slack.Client
	BotID     string
	ChannelID string
}

// ListenAndResponse listens slack events and response particular messages. It replies by slack message button.
func (s *SlackListener) ListenAndResponse() {
	rtm := s.Client.NewRTM()

	// Start listening slack events
	go rtm.ManageConnection()

	// Handle slack events
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			if err := s.handleMessageEvent(ev); err != nil {
				log.Printf("[ERROR] Failed to handle message: %s", err)
			}
		}
	}
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
	// Parse message
	m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
	if len(m) == 0 || strings.ToLower(m[0]) != "go" {
		return fmt.Errorf("invalid message")
	}

	// Display user quiz settings
	ats, cid := ComposeQuizConfig()
	UsersData.Init(ev.Msg.User)
	UsersData.SetCallbackID(ev.Msg.User, cid)

	params := slack.PostMessageParameters{
		Attachments: ats,
	}

	if _, _, err := s.Client.PostMessage(ev.Channel, "", params); err != nil {
		return fmt.Errorf("failed to post message: %s", err)
	}

	return nil
}
