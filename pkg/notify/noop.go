package notify

import (
	"log"

	"github.com/arriqaaq/kubediff/config"
	"github.com/arriqaaq/kubediff/pkg/event"
)

var _ Notifier = &NoOp{}

type NoOp struct {
}

func NewNoOp(c *config.Config) Notifier {
	return &NoOp{}
}

func (s *NoOp) Handle(e event.Event) error {
	log.Printf("Woah! Message successfully sent %+v\n ", e)
	return nil
}
