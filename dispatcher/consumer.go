package dispatcher

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/go-redis/redis"
)

type redisEvent struct {
	QuizID  uint            `json:"quiz_id"`
	EventID uint            `json:"event_id"`
	Event   json.RawMessage `json:"payload"`
}

type RedisConsumer struct {
	client     *redis.Client
	dispatcher *Dispatcher
}

const redisMaxRetriesPerReq = 2

func NewRedisConsumer(url string, dispatcher *Dispatcher) *RedisConsumer {
	rc := &RedisConsumer{
		client: redis.NewClient(&redis.Options{
			Addr:       url,
			DB:         0,
			MaxRetries: redisMaxRetriesPerReq,
		}),
		dispatcher: dispatcher,
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	rc.Consume(wg, "quiz")
	wg.Wait()
	return rc
}

func (rc *RedisConsumer) Consume(wg *sync.WaitGroup, channel string) {
	go func() {
		pss := rc.client.Subscribe(channel)
		wg.Done()
		receiver := pss.Channel()
		for {
			log.Println("Waiting for message")
			select {
			case msg := <-receiver:
				log.Println(msg)

				rc.dispatcher.broadcast <- []byte(msg.Payload)
			}
		}
	}()

}

func (rc *RedisConsumer) Publish(channel string, quizID, eventID uint, event interface{}) {
	marshaledEvent, err := json.Marshal(event)
	if err != nil {
		log.Println(err)
		return
	}

	msg := &redisEvent{
		QuizID:  quizID,
		EventID: eventID,
		Event:   marshaledEvent,
	}

	marshaledMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	if _, err := rc.client.Publish(channel, marshaledMsg).Result(); err != nil {
		log.Println(err)
	}
}
