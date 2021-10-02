package notify

import (
	"encoding/json"
	"log"

	"github.com/arriqaaq/kubediff/config"
	"github.com/arriqaaq/kubediff/pkg/event"
	"github.com/slack-go/slack"
)

var _ Notifier = &Slack{}

type Slack struct {
	Token   string
	Channel string
	Title   string
}

func NewSlack(c *config.Config) Notifier {
	return &Slack{
		Token:   c.Notifier.Slack.Token,
		Channel: c.Notifier.Slack.Channel,
		Title:   c.Notifier.Slack.Title,
	}
}

// Handle handles the notification.
func (s *Slack) Handle(e event.Event) error {
	api := slack.New(s.Token)

	message, err := json.Marshal(e)
	if err != nil {
		return err
	}

	attachment := prepareSlackAttachment(string(message), s)

	channelID, timestamp, err := api.PostMessage(s.Channel,
		slack.MsgOptionAttachments(attachment),
		slack.MsgOptionAsUser(true))
	if err != nil {
		log.Printf("%s\n", err)
		return err
	}

	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}

func prepareSlackAttachment(e string, s *Slack) slack.Attachment {

	attachment := slack.Attachment{
		Fields: []slack.AttachmentField{
			{
				Title: s.Title,
				Value: e,
			},
		},
		MarkdownIn: []string{"fields"},
	}

	return attachment
}
