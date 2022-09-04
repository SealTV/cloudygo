package kvstore

import "fmt"

var logger TransactionLogger

func InitializeTransactionLog() error {
	var err error

	// logger, err = NewFileTransactionLogger("transaction.log")
	// logger, err = NewFileTransactionLogger("transaction.log")
	logger, err = NewPostgresStransactionLogger(PostgresDBParams{
		DBName:   "transactions",
		Host:     "127.0.0.1:5432",
		User:     "root",
		Password: "1234",
	})
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()

	e, ok := Event{}, true
	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = Delete(e.Key)
			case EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	logger.Run()

	return err
}
