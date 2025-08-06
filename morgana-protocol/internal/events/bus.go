package events

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Subscriber represents a function that handles events
type Subscriber func(event Event)

// SubscriberWithFilter represents a subscriber with an optional filter function
type SubscriberWithFilter struct {
	Handler Subscriber
	Filter  func(Event) bool // Optional filter function, nil means accept all events
}

// EventBus is a thread-safe pub/sub system for events
type EventBus interface {
	// Publish publishes an event to all subscribers
	Publish(event Event)

	// PublishAsync publishes an event asynchronously without blocking
	PublishAsync(event Event) bool

	// Subscribe adds a subscriber for a specific event type
	Subscribe(eventType EventType, subscriber Subscriber) string

	// SubscribeWithFilter adds a subscriber with a filter function
	SubscribeWithFilter(eventType EventType, subscriber SubscriberWithFilter) string

	// SubscribeAll adds a subscriber for all event types
	SubscribeAll(subscriber Subscriber) string

	// Unsubscribe removes a subscriber by its ID
	Unsubscribe(subscriptionID string) bool

	// Close shuts down the event bus and cleans up resources
	Close() error

	// Stats returns statistics about the event bus
	Stats() BusStats
}

// BusStats contains statistics about the event bus
type BusStats struct {
	TotalPublished    int64             `json:"total_published"`
	TotalDropped      int64             `json:"total_dropped"`
	ActiveSubscribers int               `json:"active_subscribers"`
	QueueSize         int               `json:"queue_size"`
	QueueCapacity     int               `json:"queue_capacity"`
	SubscribersByType map[EventType]int `json:"subscribers_by_type"`
}

// bus implements the EventBus interface
type bus struct {
	// Subscriber management
	subscribers    map[EventType]map[string]SubscriberWithFilter
	allSubscribers map[string]SubscriberWithFilter
	subscriberMu   sync.RWMutex
	nextSubID      int64

	// Asynchronous event processing
	eventQueue *CircularBuffer
	workerPool chan struct{}
	workers    int
	stopped    int32
	stopCh     chan struct{}
	workerWg   sync.WaitGroup

	// Statistics
	totalPublished int64
	totalDropped   int64

	// Configuration
	config BusConfig
}

// BusConfig contains configuration for the event bus
type BusConfig struct {
	// Buffer size for asynchronous event processing (default: 10000)
	BufferSize int

	// Number of worker goroutines for processing events (default: 4)
	Workers int

	// Enable debug logging (default: false)
	Debug bool

	// Panic recovery for subscriber errors (default: true)
	RecoverPanics bool
}

// DefaultBusConfig returns a default configuration for the event bus
func DefaultBusConfig() BusConfig {
	return BusConfig{
		BufferSize:    10000,
		Workers:       4,
		Debug:         false,
		RecoverPanics: true,
	}
}

// NewEventBus creates a new event bus with the given configuration
func NewEventBus(config BusConfig) EventBus {
	if config.BufferSize <= 0 {
		config.BufferSize = 10000
	}
	if config.Workers <= 0 {
		config.Workers = 4
	}

	b := &bus{
		subscribers:    make(map[EventType]map[string]SubscriberWithFilter),
		allSubscribers: make(map[string]SubscriberWithFilter),
		eventQueue:     NewCircularBuffer(config.BufferSize),
		workerPool:     make(chan struct{}, config.Workers),
		workers:        config.Workers,
		stopCh:         make(chan struct{}),
		config:         config,
	}

	// Initialize worker pool
	for i := 0; i < config.Workers; i++ {
		b.workerPool <- struct{}{}
	}

	// Start worker goroutines
	b.startWorkers()

	return b
}

// startWorkers starts the worker goroutines for processing async events
func (b *bus) startWorkers() {
	for i := 0; i < b.workers; i++ {
		b.workerWg.Add(1)
		go b.worker()
	}
}

// worker processes events from the queue
func (b *bus) worker() {
	defer b.workerWg.Done()

	events := make([]Event, 100) // Batch size for processing events

	for {
		select {
		case <-b.stopCh:
			return
		default:
			// Try to get events from the queue
			count := b.eventQueue.PopBatch(len(events), events)
			if count == 0 {
				time.Sleep(time.Millisecond) // Brief sleep if no events
				continue
			}

			// Process the batch of events
			for i := 0; i < count; i++ {
				b.processEvent(events[i])
			}
		}
	}
}

// processEvent handles a single event by calling all appropriate subscribers
func (b *bus) processEvent(event Event) {
	if event == nil {
		return
	}

	if b.config.Debug {
		log.Printf("[EventBus] Processing event: %s (Task: %s)", event.Type(), event.TaskID())
	}

	b.subscriberMu.RLock()

	// Call type-specific subscribers
	if typeSubscribers, exists := b.subscribers[event.Type()]; exists {
		for _, subscriber := range typeSubscribers {
			b.callSubscriber(subscriber, event)
		}
	}

	// Call all-event subscribers
	for _, subscriber := range b.allSubscribers {
		b.callSubscriber(subscriber, event)
	}

	b.subscriberMu.RUnlock()
}

// callSubscriber calls a subscriber with proper error handling and recovery
func (b *bus) callSubscriber(subscriber SubscriberWithFilter, event Event) {
	// Check filter if present
	if subscriber.Filter != nil && !subscriber.Filter(event) {
		return
	}

	if b.config.RecoverPanics {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[EventBus] Subscriber panic recovered: %v", r)
			}
		}()
	}

	subscriber.Handler(event)
}

// Publish publishes an event to all subscribers synchronously
func (b *bus) Publish(event Event) {
	if atomic.LoadInt32(&b.stopped) == 1 {
		return
	}

	atomic.AddInt64(&b.totalPublished, 1)
	b.processEvent(event)
}

// PublishAsync publishes an event asynchronously without blocking
func (b *bus) PublishAsync(event Event) bool {
	if atomic.LoadInt32(&b.stopped) == 1 {
		return false
	}

	if b.eventQueue.Push(event) {
		atomic.AddInt64(&b.totalPublished, 1)
		return true
	}

	// Queue is full, drop the event
	atomic.AddInt64(&b.totalDropped, 1)
	if b.config.Debug {
		log.Printf("[EventBus] Event dropped due to full queue: %s", event.Type())
	}
	return false
}

// Subscribe adds a subscriber for a specific event type
func (b *bus) Subscribe(eventType EventType, subscriber Subscriber) string {
	return b.SubscribeWithFilter(eventType, SubscriberWithFilter{
		Handler: subscriber,
		Filter:  nil,
	})
}

// SubscribeWithFilter adds a subscriber with a filter function
func (b *bus) SubscribeWithFilter(eventType EventType, subscriber SubscriberWithFilter) string {
	subID := fmt.Sprintf("sub_%d", atomic.AddInt64(&b.nextSubID, 1))

	b.subscriberMu.Lock()
	defer b.subscriberMu.Unlock()

	if b.subscribers[eventType] == nil {
		b.subscribers[eventType] = make(map[string]SubscriberWithFilter)
	}
	b.subscribers[eventType][subID] = subscriber

	return subID
}

// SubscribeAll adds a subscriber for all event types
func (b *bus) SubscribeAll(subscriber Subscriber) string {
	subID := fmt.Sprintf("all_sub_%d", atomic.AddInt64(&b.nextSubID, 1))

	b.subscriberMu.Lock()
	defer b.subscriberMu.Unlock()

	b.allSubscribers[subID] = SubscriberWithFilter{
		Handler: subscriber,
		Filter:  nil,
	}

	return subID
}

// Unsubscribe removes a subscriber by its ID
func (b *bus) Unsubscribe(subscriptionID string) bool {
	b.subscriberMu.Lock()
	defer b.subscriberMu.Unlock()

	// Check all-event subscribers first
	if _, exists := b.allSubscribers[subscriptionID]; exists {
		delete(b.allSubscribers, subscriptionID)
		return true
	}

	// Check type-specific subscribers
	for eventType, subs := range b.subscribers {
		if _, exists := subs[subscriptionID]; exists {
			delete(subs, subscriptionID)
			if len(subs) == 0 {
				delete(b.subscribers, eventType)
			}
			return true
		}
	}

	return false
}

// Close shuts down the event bus and cleans up resources
func (b *bus) Close() error {
	if !atomic.CompareAndSwapInt32(&b.stopped, 0, 1) {
		return fmt.Errorf("event bus already closed")
	}

	close(b.stopCh)
	b.workerWg.Wait()

	// Process any remaining events in the queue
	for !b.eventQueue.IsEmpty() {
		if event := b.eventQueue.Pop(); event != nil {
			b.processEvent(event)
		}
	}

	return nil
}

// Stats returns statistics about the event bus
func (b *bus) Stats() BusStats {
	b.subscriberMu.RLock()
	defer b.subscriberMu.RUnlock()

	subsByType := make(map[EventType]int)
	totalSubs := len(b.allSubscribers)

	for eventType, subs := range b.subscribers {
		count := len(subs)
		subsByType[eventType] = count
		totalSubs += count
	}

	return BusStats{
		TotalPublished:    atomic.LoadInt64(&b.totalPublished),
		TotalDropped:      atomic.LoadInt64(&b.totalDropped),
		ActiveSubscribers: totalSubs,
		QueueSize:         b.eventQueue.Size(),
		QueueCapacity:     b.eventQueue.Capacity(),
		SubscribersByType: subsByType,
	}
}
