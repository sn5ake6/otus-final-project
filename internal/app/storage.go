package app

import (
	"context"
	"log"

	"github.com/sn5ake6/otus-final-project/internal/config"
	memorystorage "github.com/sn5ake6/otus-final-project/internal/storage/memory"
	sqlstorage "github.com/sn5ake6/otus-final-project/internal/storage/sql"
)

type Storage interface {
	Connect(ctx context.Context) error
	AddToBlacklist(subnet string) error
	DeleteFromBlacklist(subnet string) error
	FindIPInBlacklist(ip string) (bool, error)
	AddToWhitelist(subnet string) error
	DeleteFromWhitelist(subnet string) error
	FindIPInWhitelist(ip string) (bool, error)
}

func NewStorage(storageConfig config.StorageConf) Storage {
	var storage Storage

	switch storageConfig.Type {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		storage = sqlstorage.New(storageConfig.Dsn)
	default:
		log.Fatal("Unknown storage type: " + storageConfig.Type)
	}

	return storage
}
