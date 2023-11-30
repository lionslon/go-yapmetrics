package storing

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"path"
	"time"

	"github.com/lionslon/go-yapmetrics/internal/storage"
)

func Restore(s *storage.MemStorage, filePath string) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		zap.S().Error(err)
	}

	var data storage.AllMetrics
	if err := json.Unmarshal(file, &data); err != nil {
		zap.S().Error(err)
	}

	if len(data.Counter) != 0 {
		s.UpdateCounterData(data.Counter)
	}
	if len(data.Gauge) != 0 {
		s.UpdateGaugeData(data.Gauge)
	}
}

func IntervalDump(s *storage.MemStorage, filePath string, storeInterval int) {
	dir, _ := path.Split(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			zap.S().Error(err)
		}
	}
	pollTicker := time.NewTicker(time.Duration(storeInterval) * time.Second)
	defer pollTicker.Stop()
	for range pollTicker.C {
		err := dump(s, filePath)
		if err != nil {
			zap.S().Error(err)
		}
	}
}

func dump(s *storage.MemStorage, filePath string) error {

	var metrics storage.AllMetrics

	metrics.Counter = s.GetCounterData()
	metrics.Gauge = s.GetGaugeData()

	data, err := json.MarshalIndent(metrics, "", "   ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0666)
}
