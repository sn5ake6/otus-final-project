package app

import "github.com/sn5ake6/otus-final-project/internal/storage"

type LeakyBucket interface {
	Check(authorize storage.Authorize) bool
	Reset(authorize storage.Authorize)
}
