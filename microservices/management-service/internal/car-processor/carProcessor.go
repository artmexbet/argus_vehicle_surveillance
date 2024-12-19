package car_processor

import (
	"Argus/pkg/models"
	"context"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"
)

type INotificationManager interface {
	SendNotification(models.SecurityCarIDType, string) error
}

// Config represents the configuration of the car processor
type Config struct {
	InfoCount          int     `yaml:"info_count" env:"INFO_COUNT" env-default:"100"`
	MovementCheckDelay int     `yaml:"movement_check_delay" env:"MOVEMENT_CHECK_DELAY" env-default:"5"`
	Sensitivity        float64 `yaml:"sensitivity" env:"SENSITIVITY" env-default:"0.1"`
}

// CarProcessor represents the car processor
type CarProcessor struct {
	cfg      *Config
	carInfos map[models.SecurityCarIDType][]models.CarInfo
	survCars map[models.SecurityCarIDType]context.Context
	mu       sync.RWMutex
	nm       INotificationManager
}

// New creates a new car processor
func New(cfg *Config, nm INotificationManager) *CarProcessor {
	return &CarProcessor{
		cfg:      cfg,
		carInfos: make(map[models.SecurityCarIDType][]models.CarInfo),
		survCars: make(map[models.SecurityCarIDType]context.Context),
		nm:       nm,
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

// SetToSecurity sets the specified car to security
func (cp *CarProcessor) SetToSecurity(secCarID models.SecurityCarIDType) {
	slog.Info(fmt.Sprintf("Set car %v to security", secCarID))
	ctx := context.Background()
	cp.survCars[secCarID] = ctx
	go func() {
		for {
			cp.mu.RLock()
			locked := true
			select {
			case <-ctx.Done():
				return
			default:
			}

			data := cp.carInfos[secCarID]
			if len(data) == 0 {
				cp.mu.RUnlock()
				time.Sleep(time.Duration(cp.cfg.MovementCheckDelay) * time.Second)
				continue
			}
			if !data[0].IsCarFound && !data[len(data)-1].IsCarFound {
				slog.Info(fmt.Sprintf("Car %v is not found", secCarID))
				err := cp.nm.SendNotification(secCarID, "Car is not found")
				if err != nil {
					slog.Error(fmt.Sprintf("Error while sending notification: %v", err))
				}

				cp.mu.RUnlock()
				locked = false
				err = cp.StopSecurity(secCarID)
				if err != nil {
					slog.Error(fmt.Sprintf("Error while stopping security: %v", err))
					return
				}
			} else if calculateEuclideanDistance(data[0].Bbox, data[len(data)-1].Bbox) > cp.cfg.Sensitivity {
				slog.Info(fmt.Sprintf("Car %v is moving", secCarID))
				err := cp.nm.SendNotification(secCarID, "Car is moving")
				if err != nil {
					slog.Error(fmt.Sprintf("Error while sending notification: %v", err))
				}

				cp.mu.RUnlock()
				locked = false
				err = cp.StopSecurity(secCarID)
				if err != nil {
					slog.Error(fmt.Sprintf("Error while stopping security: %v", err))
					return
				}
			}
			if locked {
				cp.mu.RUnlock()
			}
			time.Sleep(time.Duration(cp.cfg.MovementCheckDelay) * time.Second)
		}
	}()
}

// StopSecurity stops security for the specified security car ID
func (cp *CarProcessor) StopSecurity(secCarID models.SecurityCarIDType) error {
	slog.Info(fmt.Sprintf("Stop security for car %v", secCarID))
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if ctx, ok := cp.survCars[secCarID]; ok {
		ctx.Done()
		delete(cp.survCars, secCarID)
		delete(cp.carInfos, secCarID)
		return nil
	}
	return fmt.Errorf("security car with ID %v not found", secCarID)
}

func pow(x float64, y float64) float64 {
	return math.Pow(x, y)
}

func calculateEuclideanDistance(bbox1, bbox2 []float32) float64 {
	return math.Sqrt(pow(float64(bbox1[0])-float64(bbox2[0]), 2) + pow(float64(bbox1[1])-float64(bbox2[1]), 2))
}
