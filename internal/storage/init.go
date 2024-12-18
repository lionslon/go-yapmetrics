package storage

type StorageWorker interface {
	Restore() error
	Dump() error
	IntervalDump()
	Check() error
	Close() error
}

type StorageProvider int

const (
	FileProvider StorageProvider = iota + 1
	DBProvider
)
