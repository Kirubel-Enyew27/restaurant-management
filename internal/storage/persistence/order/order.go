package order

import (
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"

	"go.uber.org/zap"
)

type order struct {
	db  persistencedb.PersistenceDB
	log *zap.Logger
}

func Init(db persistencedb.PersistenceDB, log *zap.Logger) storage.Order {
	return &order{
		db:  db,
		log: log,
	}
}
