package storage

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"os"
	"path"
	"time"
)

// fileProvider структура для работы со структурой данных в файле
type fileProvider struct {
	filePath      string
	storeInterval int
	st            *MemStorage
}

// Check структура для работы со структурой данных в файле
func (f *fileProvider) Check() error {
	return errors.New("not provided for this storage type")
}

// NewFileProvider конструктор для работы со структурой данных в файле
func NewFileProvider(filePath string, storeInterval int, m *MemStorage) StorageWorker {
	return &fileProvider{
		filePath:      filePath,
		storeInterval: storeInterval,
		st:            m,
	}
}

// Dump подчищает и записывает в файл
func (f *fileProvider) Dump() error {
	dir, _ := path.Split(f.filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			zap.S().Error(err)
		}
	}

	data, err := json.MarshalIndent(f.st, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(f.filePath, data, 0666)
}

// IntervalDump обертка над записью в файле с указанным интервалом
func (f *fileProvider) IntervalDump() {
	pollTicker := time.NewTicker(time.Duration(f.storeInterval) * time.Second)
	defer pollTicker.Stop()
	for range pollTicker.C {
		err := f.Dump()
		if err != nil {
			zap.S().Error(err)
		}
	}
}

// Restore восстанавливает данные из файла
func (f *fileProvider) Restore() error {
	file, err := os.ReadFile(f.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(file, f.st)
}
