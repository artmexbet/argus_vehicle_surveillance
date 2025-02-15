package car_processor

import (
	"Argus/pkg/models"
	"github.com/google/uuid"
	"math/rand"
	"slices"
	"testing"
	"time"
)

type NotificationManager struct {
	Notifications map[models.SecurityCarIDType][]string
}

func NewNotificationManager() *NotificationManager {
	return &NotificationManager{
		Notifications: make(map[models.SecurityCarIDType][]string),
	}
}

func (nm *NotificationManager) SendNotification(secCarID models.SecurityCarIDType, text string) error {
	if _, ok := nm.Notifications[secCarID]; !ok {
		nm.Notifications[secCarID] = make([]string, 0)
	}
	nm.Notifications[secCarID] = append(nm.Notifications[secCarID], text)
	return nil
}

type Database struct {
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) StopVehicleTracking(secId models.SecurityCarIDType) error {
	return nil
}

func (db *Database) AddEvent(secId models.SecurityCarIDType, eventTypeId int64) error {
	return nil
}

func (db *Database) GetAllSecuriedCarsByCamera(models.CameraIDType) ([]models.SecurityCar, error) {
	return nil, nil
}

func initCP(nm INotificationManager, db IDatabase) *CarProcessor {
	cfg := &Config{
		InfoCount:          100,
		MovementCheckDelay: 5,
		Sensitivity:        0.01,
	}
	return New(cfg, nm, db)
}

func TestCarProcessor_AppendCarInfo(t *testing.T) {
	nm := NewNotificationManager()
	db := NewDatabase()
	cp := initCP(nm, db)
	secCarID := models.SecurityCarIDType(uuid.New())
	carInfo := models.CarInfo{
		ID:         models.CarIDType(1),
		Bbox:       []float32{1, 1, 1, 1},
		IsCarFound: true,
	}

	cp.AppendCarInfo(secCarID, carInfo)
	if len(cp.carInfos[secCarID]) != 1 {
		t.Fatalf("expected 1, got %d", len(cp.carInfos[secCarID]))
	}
	if cp.carInfos[secCarID][0].ID != carInfo.ID || slices.Compare(cp.carInfos[secCarID][0].Bbox, carInfo.Bbox) != 0 {
		t.Fatalf("expected %v, got %v", carInfo, cp.carInfos[secCarID][0])
	}
}

// generateRandomFloat generates a random float32 number from 0 to 1
func generateRandomFloat() float32 {
	return rand.Float32()
}

// generateStayingCarInfo generates n car info with the near-staying bbox
func generateStayingCarInfo(n int, carId models.CarIDType) []models.CarInfo {
	res := make([]models.CarInfo, 0, n)
	startBbox := []float32{generateRandomFloat(), generateRandomFloat(),
		generateRandomFloat(), generateRandomFloat()}
	for i := 0; i < n; i++ {
		for j := 0; j < 4; j++ {
			tmp := startBbox[j] + generateRandomFloat()/500*(rand.Float32()-0.5)
			if tmp < 0 {
				tmp = 0
			} else if tmp > 1 {
				tmp = 1
			}
			startBbox[j] = tmp
		}
		newBbox := make([]float32, 4)
		copy(newBbox, startBbox)
		res = append(res, models.CarInfo{
			ID:         carId,
			Bbox:       newBbox,
			IsCarFound: true,
		})
	}
	return res
}

func generateMovingCarInfo(n int, carId models.CarIDType) []models.CarInfo {
	res := make([]models.CarInfo, 0, n)
	startBbox := []float32{generateRandomFloat(), generateRandomFloat(),
		generateRandomFloat(), generateRandomFloat()}
	for i := 0; i < n; i++ {
		for j := 0; j < 4; j++ {
			tmp := startBbox[j] + generateRandomFloat()/10*(rand.Float32()-0.5)
			if tmp < 0 {
				tmp = 0
			} else if tmp > 1 {
				tmp = 1
			}
			startBbox[j] = tmp
		}

		newBbox := make([]float32, 4)
		copy(newBbox, startBbox)

		res = append(res, models.CarInfo{
			ID:         carId,
			Bbox:       newBbox,
			IsCarFound: true,
		})
	}
	return res
}

func generateNotFoundCarInfo(n int, carId models.CarIDType) []models.CarInfo {
	res := make([]models.CarInfo, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, models.CarInfo{
			ID:         carId,
			Bbox:       []float32{0, 0, 0, 0},
			IsCarFound: false,
		})
	}
	return res
}

func TestCarProcessor_SetToSecurity(t *testing.T) {
	FPS := 15

	nm := NewNotificationManager()
	db := NewDatabase()
	cp := initCP(nm, db)
	secCarID := models.SecurityCarIDType(uuid.New())
	carInfos := generateStayingCarInfo(200, models.CarIDType(1))
	cp.SetToSecurity(secCarID)
	if _, ok := cp.survCars[secCarID]; !ok {
		t.Fatalf("expected true, got false")
		return
	}

	// Test staying car
	for _, carInfo := range carInfos {
		cp.AppendCarInfo(secCarID, carInfo)
		if nm.Notifications[secCarID] != nil {
			t.Fatalf("expected nil, got %v. system sent notification", nm.Notifications[secCarID])
			return
		}
		time.Sleep(time.Second / time.Duration(FPS))
	}

	// Test stop security
	err := cp.StopSecurity(secCarID)
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
		return
	}

	if _, ok := cp.survCars[secCarID]; ok {
		t.Fatalf("expected false, got true")
		return
	}

	if len(cp.carInfos[secCarID]) != 0 {
		t.Fatalf("expected 0, got %d", len(cp.carInfos[secCarID]))
		return
	}
}

func TestCarProcessor_SetToSecurityMovingCar(t *testing.T) {
	FPS := 15

	nm := NewNotificationManager()
	db := NewDatabase()
	cp := initCP(nm, db)
	secCarID := models.SecurityCarIDType(uuid.New())
	carInfos := generateMovingCarInfo(200, models.CarIDType(1))
	cp.SetToSecurity(secCarID)

	isNotified := false
	for _, carInfo := range carInfos {
		cp.AppendCarInfo(secCarID, carInfo)
		if nm.Notifications[secCarID] != nil {
			isNotified = true
		}
		time.Sleep(time.Second / time.Duration(FPS))
	}

	if !isNotified {
		t.Fatalf("expected true, got false. system didn't send notification")
		return
	}

	if notifications := nm.Notifications[secCarID]; notifications[0] != "Car is moving" {
		t.Fatalf("expected Car is moving, got %s", notifications[0])
		return
	}

	if _, ok := cp.survCars[secCarID]; ok {
		t.Fatalf("expected false, got true")
		return
	}
}

func TestCarProcessor_SetToSecurityNotFoundCar(t *testing.T) {
	FPS := 15

	nm := NewNotificationManager()
	cp := initCP(nm, NewDatabase())
	secCarID := models.SecurityCarIDType(uuid.New())
	carInfos := generateNotFoundCarInfo(200, models.CarIDType(1))
	cp.SetToSecurity(secCarID)

	isNotified := false
	for _, carInfo := range carInfos {
		cp.AppendCarInfo(secCarID, carInfo)
		if nm.Notifications[secCarID] != nil {
			isNotified = true
		}
		time.Sleep(time.Second / time.Duration(FPS))
	}

	if !isNotified {
		t.Fatalf("expected true, got false. system didn't send notification")
		return
	}

	if notifications := nm.Notifications[secCarID]; notifications[0] != "Car is not found" {
		t.Fatalf("expected Car is moving, got %s", notifications[0])
		return
	}

	if _, ok := cp.survCars[secCarID]; ok {
		t.Fatalf("expected false, got true")
		return
	}
}
