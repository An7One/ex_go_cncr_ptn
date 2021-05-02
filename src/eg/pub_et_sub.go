// Reference:
//	https://chai2010.gitbooks.io/advanced-go-programming-book/content/ch1-basic/ch1-06-goroutine.html
package example

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type (
	subscriber chan interface{}         // the subscriber is a channel
	topicFunc  func(v interface{}) bool // the topic is a filter
)

// Publisher
type Publisher struct {
	m           sync.RWMutex             // read write lock
	buffer      int                      // the size of the buffer of the subscriber
	timeout     time.Duration            // timeout duration of the publishing
	subscribers map[subscriber]topicFunc // info of the subscriber
}

// NewPublisher
func NewPublisher(publishTimeout time.Duration, buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		timeout:     publishTimeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// Subscribe
func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

// SubscribeTopic
func (p *Publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.subscribers[ch] = topic
	p.m.Unlock()
	return ch
}

// Evict cancels the subscription
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// Publish publishes a new topic
func (p *Publisher) Publish(v interface{}) {
	p.m.Lock()
	defer p.m.RUnlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// Close closes the publisher, as well as the channels from the subscribers
func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

// sendTopic
func (p *Publisher) sendTopic(
	sub subscriber,
	topic topicFunc,
	v interface{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}

// DemoPubNSub demos the publisher and subscriber concurrency pattern
func DemoPubNSub() {
	p := NewPublisher(100*time.Millisecond, 10)
	defer p.Close()

	all := p.Subscribe()
	golang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}
		return false
	})

	p.Publish("hello world 01")
	p.Publish("hello world 02")

	go func() {
		for msg := range all {
			fmt.Println("all:", msg)
		}
	}()

	go func() {
		for msg := range golang {
			fmt.Println("golang:", msg)
		}
	}()

	time.Sleep(3 * time.Second)
}
