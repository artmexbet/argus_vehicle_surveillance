package car_processor

import (
	"Argus/pkg/models"
	"sync"
	"time"
)

// Config represents the configuration of the car processor
type Config struct {
	InfoCount          int `yaml:"info_count" env:"INFO_COUNT" env-default:"100"`
	MovementCheckDelay int `yaml:"movement_check_delay" env:"MOVEMENT_CHECK_DELAY" env-default:"5"`
}

// CarProcessor represents the car processor
type CarProcessor struct {
	cfg      *Config
	carInfos map[models.SecurityCarIDType][]models.CarInfo
	mu       sync.RWMutex
}

// New creates a new car processor
func New(cfg *Config) *CarProcessor {
	return &CarProcessor{
		cfg:      cfg,
		carInfos: make(map[models.SecurityCarIDType][]models.CarInfo),
	}
}

// AppendCarInfo appends car info to the car processor
func (cp *CarProcessor) AppendCarInfo(secCarID models.SecurityCarIDType, carInfo models.CarInfo) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Initialize car info slice if it doesn't exist
	if _, ok := cp.carInfos[secCarID]; !ok {
		cp.carInfos[secCarID] = make([]models.CarInfo, 0, cp.cfg.InfoCount)
	}

	// Remove the first element if the slice is full
	if len(cp.carInfos[secCarID]) >= cp.cfg.InfoCount {
		cp.carInfos[secCarID] = cp.carInfos[secCarID][1:]
	}

	cp.carInfos[secCarID] = append(cp.carInfos[secCarID], carInfo)
}

// GetCarInfos returns car infos for the specified security car ID
func (cp *CarProcessor) GetCarInfos(secCarID models.SecurityCarIDType) []models.CarInfo {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	return cp.carInfos[secCarID]
}

func (cp *CarProcessor) SetToSecurity(secCarID models.SecurityCarIDType, event chan struct{}) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	// Do something with the message
	// If the car is moving for a long time, send a message to the security service

	time.Sleep(time.Duration(cp.cfg.MovementCheckDelay) * time.Second)
}
