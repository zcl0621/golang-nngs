package handlerHistory

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"higo-game-node/mongodb"
	"sync"
	"time"
)

var Queue *DelayedQueue
var collection *mongo.Collection

func InitQueue() {
	Queue = &DelayedQueue{}

	mdb := mongodb.GetMgoCli().Database("higo-game-node")
	collection = mdb.Collection("handler-history")
}

type Message struct {
	Content *HandlerHistory
}

type DelayedQueue struct {
	mu       sync.Mutex
	messages []*Message
	timer    *time.Timer
}

func (dq *DelayedQueue) ScheduleMessage(msg *Message, delay time.Duration) {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	dq.messages = append(dq.messages, msg)

	if dq.timer == nil {
		dq.timer = time.AfterFunc(delay, func() {
			dq.Flush()
		})
	} else {
		if dq.timer.Stop() {
			dq.timer.Reset(delay)
		} else {
			<-dq.timer.C
			dq.Flush()
		}
	}
}

func (dq *DelayedQueue) Flush() {
	dq.mu.Lock()
	defer dq.mu.Unlock()

	for _, msg := range dq.messages {
		collection.InsertOne(context.Background(), msg)
	}

	dq.messages = nil
	dq.timer = nil
}
