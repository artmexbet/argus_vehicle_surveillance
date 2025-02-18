package car_processor

import (
	"Argus/pkg/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"math"
	"sync"
	"time"
)

type INotificationManager interface {
	SendNotification(models.SecurityCarIDType, string) error
}

type IDatabase interface {
	StopVehicleTracking(models.SecurityCarIDType) error
	AddEvent(models.SecurityCarIDType, int64) error
	GetAllSecuriedCarsByCamera(models.CameraIDType) ([]models.SecurityCar, error)
}

// Config represents the configuration of the car processor
type Config struct {
	InfoCount          int     `yaml:"info_count" env:"INFO_COUNT" env-default:"100"`
	MovementCheckDelay int     `yaml:"movement_check_delay" env:"MOVEMENT_CHECK_DELAY" env-default:"5"`
	LostCarTimeout     int     `yaml:"lost_car_timeout" env:"LOST_CAR_TIMEOUT" env-default:"60"`
	Sensitivity        float64 `yaml:"sensitivity" env:"SENSITIVITY" env-default:"0.1"`
}

// CarProcessor represents the car processor
type CarProcessor struct {
	cfg      *Config
	carInfos map[models.SecurityCarIDType][]models.CarInfo
	survCars map[models.SecurityCarIDType]context.Context
	lostCars map[models.SecurityCarIDType]time.Time
	mu       sync.RWMutex
	nm       INotificationManager
	db       IDatabase
}

// New creates a new car processor
func New(cfg *Config, nm INotificationManager, db IDatabase) *CarProcessor {
	cp := &CarProcessor{
		cfg:      cfg,
		carInfos: make(map[models.SecurityCarIDType][]models.CarInfo),
		survCars: make(map[models.SecurityCarIDType]context.Context),
		lostCars: make(map[models.SecurityCarIDType]time.Time),
		nm:       nm,
		db:       db,
	}

	// Add all cars which is under security to the car processor
	tmp, _ := uuid.Parse("25d3e590-9870-11ef-a686-0242ac130002")
	var cameraID = models.CameraIDType(tmp)
	cars, _ := db.GetAllSecuriedCarsByCamera(cameraID)
	for _, car := range cars {
		cp.SetToSecurity(car.ID)
	}

	return cp
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

	// Checks if the car is not found we should set the last bbox to assigned
	// If it is a first carInfo we should mark for us that car is not found (0, 0)
	if !carInfo.IsCarFound {
		if len(cp.carInfos[secCarID]) > 0 {
			carInfo.Bbox = cp.carInfos[secCarID][len(cp.carInfos[secCarID])-1].Bbox
		} else {
			carInfo.Bbox = []float32{0, 0}
		}
	}

	cp.carInfos[secCarID] = append(cp.carInfos[secCarID], carInfo)

	// If the car is not found, and it is not in the lost cars map, add it.
	// It needs to be added to the lost cars map to track the time of the car being lost.
	// Then we check if the car was lost for a long time, we should stop tracking it and send notifications.
	if _, ok := cp.lostCars[secCarID]; !carInfo.IsCarFound && !ok {
		cp.lostCars[secCarID] = time.Now()
	} else if carInfo.IsCarFound {
		delete(cp.lostCars, secCarID)
	}
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

			// There we check if the car is lost for a long time, we should stop tracking it and send notifications.
			if lostTime, ok := cp.lostCars[secCarID]; ok && time.Since(lostTime).Seconds() > float64(cp.cfg.LostCarTimeout) {
				slog.Info(fmt.Sprintf("Car %v is not found", secCarID))

				err := cp.db.AddEvent(secCarID, 1) // Adding events to database. Note: maybe it should be added to the timeseries database in the future
				if err != nil {
					slog.Error(fmt.Sprintf("Error while adding event: %v", err))
				}

				err = cp.nm.SendNotification(secCarID, "Car is not found")
				if err != nil {
					slog.Error(fmt.Sprintf("Error while sending notification: %v", err))
				}

				cp.mu.RUnlock()
				locked = false

				err = cp.db.StopVehicleTracking(secCarID)
				if err != nil {
					slog.Error(fmt.Sprintf("Error while stopping vehicle tracking: %v", err))
				}

				err = cp.StopSecurity(secCarID)
				if err != nil {
					slog.Error(fmt.Sprintf("Error while stopping security: %v", err))
					return
				}
			} else if calculateEuclideanDistance(data[0].Bbox, data[len(data)-1].Bbox) > cp.cfg.Sensitivity {
				slog.Info(fmt.Sprintf("Car %v is moving", secCarID))
				err := cp.db.AddEvent(secCarID, 1)
				if err != nil {
					slog.Error(fmt.Sprintf("Error while adding event: %v", err))
				}

				err = cp.nm.SendNotification(secCarID, "Car is moving")
				if err != nil {
					slog.Error(fmt.Sprintf("Error while sending notification: %v", err))
				}

				cp.mu.RUnlock()
				locked = false

				err = cp.db.StopVehicleTracking(secCarID)
				if err != nil {
					slog.Error(fmt.Sprintf("Error while stopping vehicle tracking: %v", err))
				}

				err = cp.StopSecurity(secCarID)
				if err != nil {
					slog.Error(fmt.Sprintf("Error while stopping security: %v", err))
					return
				}
			}
			if locked {
				cp.mu.RUnlock()
			}
			slog.Debug("distance",
				slog.Any("dist", calculateEuclideanDistance(data[0].Bbox, data[len(data)-1].Bbox)),
				slog.Any("security id", secCarID))
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
