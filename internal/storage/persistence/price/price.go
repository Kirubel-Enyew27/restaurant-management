package price

import (
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"

	"go.uber.org/zap"
)

type price struct {
	db  persistencedb.PersistenceDB
	log *zap.Logger
}

func Init(db persistencedb.PersistenceDB, log *zap.Logger) storage.Price {
	return &price{
		db:  db,
		log: log,
	}
}
