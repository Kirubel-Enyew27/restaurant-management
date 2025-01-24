package food

import (
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"

	"go.uber.org/zap"
)

type food struct {
	db  persistencedb.PersistenceDB
	log *zap.Logger
}

func Init(db persistencedb.PersistenceDB, log *zap.Logger) storage.Food {
	return &food{
		db:  db,
		log: log,
	}
}
