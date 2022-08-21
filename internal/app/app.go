package app

import (
	"context"

	"github.com/sn5ake6/otus-final-project/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
	Bucket  LeakyBucket
}

func New(logger Logger, storage Storage, bucket LeakyBucket) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
		Bucket:  bucket,
	}
}

func (a *App) Authorize(ctx context.Context, authorize storage.Authorize) (bool, error) {
	res, err := a.Storage.FindIPInBlacklist(authorize.IP)
	if err != nil {
		return false, err
	}

	if res {
		return false, nil
	}

	res, err = a.Storage.FindIPInWhitelist(authorize.IP)
	if err != nil {
		return false, err
	}

	if res {
		return true, nil
	}

	return a.Bucket.Check(authorize), nil
}

func (a *App) Reset(ctx context.Context, authorize storage.Authorize) {
	a.Bucket.Reset(authorize)
}

func (a *App) AddToBlacklist(ctx context.Context, subnet string) error {
	return a.Storage.AddToBlacklist(subnet)
}

func (a *App) DeleteFromBlacklist(ctx context.Context, subnet string) error {
	return a.Storage.DeleteFromBlacklist(subnet)
}

func (a *App) FindIPInBlacklist(ctx context.Context, ip string) (bool, error) {
	return a.Storage.FindIPInBlacklist(ip)
}

func (a *App) AddToWhitelist(ctx context.Context, subnet string) error {
	return a.Storage.AddToWhitelist(subnet)
}

func (a *App) DeleteFromWhitelist(ctx context.Context, subnet string) error {
	return a.Storage.DeleteFromWhitelist(subnet)
}

func (a *App) FindIPInWhitelist(ctx context.Context, ip string) (bool, error) {
	return a.Storage.FindIPInBlacklist(ip)
}
