package notify

import (
	"github.com/arriqaaq/kubediff/config"
	"github.com/arriqaaq/kubediff/pkg/event"
	"k8s.io/apimachinery/pkg/util/errors"
)

type Notifier interface {
	Handle(e event.Event) error
}

type NotifierList struct {
	notifiers []Notifier
}

func (n *NotifierList) Handle(e event.Event) error {

	errs := make([]error, 0, len(n.notifiers))
	for _, n := range n.notifiers {
		if err := n.Handle(e); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.NewAggregate(errs)
	}
	return nil
}

func NewNotifierList(conf *config.Config) Notifier {
	var notifiers []Notifier
	if conf.Notifier.Slack.Enabled {
		notifiers = append(notifiers, NewSlack(conf))
	}
	if conf.Notifier.Webhook.Enabled {
		notifiers = append(notifiers, NewWebhook(conf))
	}
	if conf.Notifier.NoOp != "" {
		notifiers = append(notifiers, NewNoOp(conf))
	}

	return &NotifierList{notifiers}
}
