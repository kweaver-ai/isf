// Package outboxer is an implementation of the outbox pattern.
// The producer of messages can durably store those messages in a local outbox before sending to a Message Endpoint.
// The durable local storage may be implemented in the Message Channel directly, especially when combined
// with Idempotent Messages.
package outboxer

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

const (
	messageBatchSize = 1
	cleanUpBatchSize = 1
)

var (
	// ErrMissingEventStream is used when no event stream is provided.
	ErrMissingEventStream = errors.New("an event stream is required for the outboxer to work")

	// ErrMissingDataStore is used when no data store is provided.
	ErrMissingDataStore = errors.New("a data store is required for the outboxer to work")

	ErrNoEventsLeft = errors.New("no events left")
)

//// ExecerContext defines the exec context method that is used within a transaction.
//type ExecerContext interface {
//	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
//}

// DataStore defines the data store methods.
type DataStore interface {
	GetEvent(tx *gorm.DB) (*OutboxMessage, bool, error)

	Add(m *OutboxMessage) error
	//AddWithinTx(ctx context.Context, m *OutboxMessage, fn func(ExecerContext) error) error

	// SetAsDispatched show be in get event and commit if the tx commit
	SetAsDispatched(tx *gorm.DB, id uint64) error

	Remove(retentionTime time.Duration, batchSize int32) error

	GetTx() *gorm.DB
}

// EventStream defines the event stream methods.
type EventStream interface {
	Send(*OutboxMessage) error
}

// Outboxer implements the outbox pattern.
type Outboxer struct {
	ds                   DataStore
	es                   EventStream
	checkInterval        time.Duration
	cleanUpInterval      time.Duration
	cleanUpRetentionTime time.Duration
	cleanUpBatchSize     int32
	messageBatchSize     int32

	errChan           chan error
	okChan            chan struct{}
	manualTriggerChan chan struct{}
}

// New creates a new instance of Outboxer.
func New(opts ...Option) (*Outboxer, error) {
	o := Outboxer{
		errChan:           make(chan error),
		okChan:            make(chan struct{}),
		manualTriggerChan: make(chan struct{}),
		messageBatchSize:  messageBatchSize,
		cleanUpBatchSize:  cleanUpBatchSize,
	}

	for _, opt := range opts {
		opt(&o)
	}

	if o.ds == nil {
		return nil, ErrMissingDataStore
	}

	if o.es == nil {
		return nil, ErrMissingEventStream
	}

	return &o, nil
}

// ErrChan returns the error channel.
func (o *Outboxer) ErrChan() <-chan error {
	return o.errChan
}

// OkChan returns the ok channel that is used when each message is successfully delivered.
func (o *Outboxer) OkChan() <-chan struct{} {
	return o.okChan
}

func (o *Outboxer) manualTrigger() {
	go func() { o.manualTriggerChan <- struct{}{} }()
}

// Send sends a message.
func (o *Outboxer) Send(m *OutboxMessage) error {
	err := o.ds.Add(m)
	if err == nil {
		o.manualTrigger()
	}
	return err
}

//// TODO SendWithinTx encapsulate any database call within a transaction.
//func (o *Outboxer) SendWithinTx(ctx context.Context, evt *OutboxMessage, fn func(ExecerContext) error) error {
//	return o.ds.AddWithinTx(ctx, evt, fn)
//}

// Start encapsulates two go routines. Starts the dispatcher, which is responsible for getting the messages
// from the data store and sending to the event stream.
// Starts the cleanup process, that makes sure old messages are removed from the data store.
func (o *Outboxer) Start(ctx context.Context) {
	go o.StartDispatcher(ctx)
	go o.StartCleanup(ctx)
}

// StartDispatcher starts the dispatcher, which is responsible for getting the messages
// from the data store and sending to the event stream.
func (o *Outboxer) StartDispatcher(ctx context.Context) {
	// TODO 更详细的error
	ticker := time.NewTicker(o.checkInterval)
	dispatch := func() {
		for {
			err := o.ds.GetTx().Transaction(func(tx *gorm.DB) error {
				evt, hasEvent, err := o.ds.GetEvent(tx)
				if err != nil {
					return err
				}

				if !hasEvent {
					return ErrNoEventsLeft
				}

				if err := o.ds.SetAsDispatched(tx, evt.ID); err != nil {
					return err
				}

				if err := o.es.Send(evt); err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				if !errors.Is(err, ErrNoEventsLeft) {
					o.errChan <- err
				}
				break
			}

			o.okChan <- struct{}{}
		}
	}

	for {
		select {
		// TODO Send SetAsDispatched order rollback
		//TODO single getEvents and SetAsDisPatched in a tx or  gorm tx lock row lock

		case <-ticker.C:
			dispatch()
		case <-o.manualTriggerChan:
			dispatch()
		case <-ctx.Done():
			// TODO some log?
			return
		}
	}
}

// StartCleanup starts the cleanup process, that makes sure old messages are removed from the data store.
func (o *Outboxer) StartCleanup(ctx context.Context) {
	// TODO add comme log
	ticker := time.NewTicker(o.cleanUpInterval)
	for {
		select {
		case <-ticker.C:
			if err := o.ds.Remove(o.cleanUpRetentionTime, o.cleanUpBatchSize); err != nil {
				o.errChan <- err
			}
		case <-ctx.Done():
			// TODO some log?
			return
		}
	}
}

// Stop closes all channels.
func (o *Outboxer) Stop() {
	close(o.errChan)
	close(o.okChan)
}
