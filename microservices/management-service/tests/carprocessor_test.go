package carprocessor_test

import (
	carProcessor "Argus/microservices/management-service/internal/car-processor"
	"Argus/pkg/models"
	"github.com/google/uuid"
	"math/rand/v2"
	"testing"
)

func addRandomVarience(bboxes [][]float32, multiplier float32) {
	for i := range bboxes {
		for j := range bboxes[i] {
			bboxes[i][j] += rand.Float32() * multiplier
			if bboxes[i][j] < 0 {
				bboxes[i][j] = 0
			} else if bboxes[i][j] > 1 {
				bboxes[i][j] = 1
			}
		}
	}
	for i := range bboxes[0] {
		// guarantee at least one move this far
		bboxes[0][i] = bboxes[0][i] + multiplier
	}
}

func TestCarProcessor(t *testing.T) {
	config := carProcessor.Config{InfoCount: 100, MovementCheckDelay: 5, MovementThreshhold: 0.02}
	var uid uuid.UUID
	uid, _ = uuid.NewRandom()
	id := models.SecurityCarIDType(uid)

	cp := carProcessor.New(&config)
	bboxes := [][]float32{
		{0.5, 0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5, 0.5},
		{0.5, 0.5, 0.5, 0.5},
	}
    // TODO: this should be derived from the real data
	varience := float32(0.005)
	addRandomVarience(bboxes, varience)

	for _, bbox := range bboxes {
		cp.AppendCarInfo(id, models.CarInfo{
			ID:         1,
			Bbox:       bbox,
			IsCarFound: true,
		})
	}

	shouldMove := cp.ShouldNotify(id)
	if shouldMove {
		t.Fatal("This should not notify")
	}
}
