package state

import (
	"context"
	"sync"

	"github.com/keep-network/keep-core/pkg/net"
)

type messageReceiver struct {
	channel net.BroadcastChannel

	messages      []net.Message
	messagesMutex sync.RWMutex
}

func runMessageReceiver(
	ctx context.Context,
	channel net.BroadcastChannel,
	bufferSize int,
) *messageReceiver {
	messageReceiver := &messageReceiver{
		channel:  channel,
		messages: make([]net.Message, 0),
	}

	go messageReceiver.run(ctx, bufferSize)

	return messageReceiver
}

func (mr *messageReceiver) run(ctx context.Context, bufferSize int) {
	receiverChannel := make(chan net.Message, bufferSize)

	mr.channel.Recv(ctx, func(message net.Message) {
		receiverChannel <- message
	})

	for {
		select {
		case message := <-receiverChannel:
			mr.messagesMutex.Lock()
			mr.messages = append(mr.messages, message)
			mr.messagesMutex.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (mr *messageReceiver) snapshot() []net.Message {
	mr.messagesMutex.RLock()
	defer mr.messagesMutex.RUnlock()

	messagesSnapshot := make([]net.Message, len(mr.messages))
	copy(messagesSnapshot, mr.messages)

	return messagesSnapshot
}
