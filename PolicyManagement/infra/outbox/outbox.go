package outbox

import (
	"context"
	"time"

	clog "policy_mgnt/utils/gocommon/v2/log"
	mariadb "policy_mgnt/utils/gocommon/v2/outbox/storage/mariadb_mysql"

	cdb "policy_mgnt/utils/gocommon/v2/db"
	outboxer "policy_mgnt/utils/gocommon/v2/outbox"
)

var outboxerSingleton *outboxer.Outboxer

func InitOutBoxer() error {
	ds, err := mariadb.WithInstance(cdb.NewDB())
	if err != nil {
		return err
	}

	outboxerSingleton, err = outboxer.New(
		outboxer.WithDataStore(ds),
		outboxer.WithEventStream(NewRedisES()),
		outboxer.WithCheckInterval(10*time.Second),
		outboxer.WithMessageBatchSize(1),
		outboxer.WithCleanupInterval(5*time.Second),
		outboxer.WithCleanupBatchSize(10), // TODO
		outboxer.WithCleanupRetentionTime(2*24*time.Hour),
	)
	if err != nil {
		return err
	}

	return nil
}

func NewOutBoxer() *outboxer.Outboxer {
	return outboxerSingleton
}

func StartOutBoxer(ctx context.Context) error {
	go func() {
		// TODO 绑定个logger?
		logger := clog.NewLogger()
		outboxerSingleton.Start(ctx)
		defer outboxerSingleton.Stop()
		// we can also listen for errors and ok messages that were send
		for {
			select {
			case err := <-outboxerSingleton.ErrChan():
				logger.Errorf("could not send message: %s", err)
			case <-outboxerSingleton.OkChan():
				logger.Info("message received")
			}
		}
	}()
	return nil
}
