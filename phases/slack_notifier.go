package phases

import (
	"github.com/Everlane/evan/common"
	"github.com/Everlane/evan/slack"
)

type SlackNotifierPhase struct {
	// URL for the Slack Incoming Webhook
	Webhook string
	Channel string
	Format  func(common.Deployment) (*slack.Payload, error)
}

func (snp *SlackNotifierPhase) CanPreload() bool {
	return false
}

func (snp *SlackNotifierPhase) Execute(deployment common.Deployment, _ interface{}) error {
	payload, err := snp.Format(deployment)
	if err != nil {
		return err
	}

	// Don't send a message to Slack if the format function didn't return
	// a message to send
	if payload == nil {
		return nil
	}

	err = slack.Send(snp.Webhook, "", *payload)
	if err != nil {
		return err
	}

	return nil
}
