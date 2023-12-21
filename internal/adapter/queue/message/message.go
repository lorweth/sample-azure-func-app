package message

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	pkgerrors "github.com/pkg/errors"
)

type Queue struct {
	client *azqueue.QueueClient
}

func New(client *azqueue.QueueClient) Queue {
	return Queue{
		client: client,
	}
}

func (p Queue) Publish(ctx context.Context, msg string) error {
	ctx, cancel := context.WithTimeout(ctx, sendMessageTimeout)
	defer cancel()

	_, err := p.client.EnqueueMessage(ctx, msg, nil)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}
