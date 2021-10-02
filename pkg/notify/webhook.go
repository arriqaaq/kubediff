package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arriqaaq/kubediff/config"
	"github.com/arriqaaq/kubediff/pkg/event"
	"github.com/arriqaaq/kubediff/pkg/log"
)

var _ Notifier = &Webhook{}

type Webhook struct {
	URL string
}

type payload struct {
	Event     event.Event `json:"meta"`
	TimeStamp time.Time   `json:"timestamp"`
}

func NewWebhook(c *config.Config) Notifier {
	return &Webhook{
		URL: c.Notifier.Webhook.Url,
	}
}

func (w *Webhook) Handle(event event.Event) (err error) {
	pl := &payload{
		Event:     event,
		TimeStamp: time.Now(),
	}

	err = w.post(pl)
	if err != nil {
		log.Error(err)
		log.Debugf("error sending event to webhook %v", event)
	}

	log.Debugf("Event successfully sent to Webhook %v", event)
	return nil
}

func (w *Webhook) post(pl *payload) error {

	message, err := json.Marshal(pl)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", w.URL, bytes.NewBuffer(message))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error response from webhook: %s", fmt.Sprint(resp.StatusCode))
	}

	return nil
}
