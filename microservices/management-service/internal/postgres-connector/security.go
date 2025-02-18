package postgres_connector

import (
	"Argus/pkg/models"
	"context"
	"log/slog"
	"time"
)

func (p *PostgresConnector) GetAllSecuriedCarsByCamera(cameraID models.CameraIDType) ([]models.SecurityCar, error) {
	rows, err := p.conn.Query(context.Background(),
		`SELECT id, camera_id, car_id FROM security_cars
                             WHERE security_date_off IS NULL AND camera_id = $1;`, cameraID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var securityCars []models.SecurityCar
	for rows.Next() {
		var securityCar models.SecurityCar
		err = rows.Scan(&securityCar.ID, &securityCar.CameraID, &securityCar.CarID)
		if err != nil {
			return nil, err
		}
		securityCars = append(securityCars, securityCar)
	}

	return securityCars, nil
}

func (p *PostgresConnector) SetCarToSecurity(
	carID models.CarIDType,
	cameraID models.CameraIDType,
	accountID models.AccountIDType,
	securityDateOn models.TimestampType,
) (models.SecurityCarIDType, error) {
	tx, err := p.conn.Begin(context.Background())
	if err != nil {
		return models.SecurityCarIDType{}, err
	}

	var securityCarID models.SecurityCarIDType
	err = tx.QueryRow(context.Background(),
		`INSERT INTO security_cars (car_id, camera_id, account_id, security_date_on) 
					VALUES ($1, $2, $3, $4) RETURNING id`,
		carID, cameraID, accountID, time.Time(securityDateOn)).Scan(&securityCarID)
	if err == nil {
		tx.Commit(context.Background())
		slog.Info("Car is set to security", slog.Any("securityCarID", securityCarID))
	} else {
		tx.Rollback(context.Background())
	}
	return securityCarID, err
}

// StopVehicleTracking stops tracking of the specified security car by secId
func (p *PostgresConnector) StopVehicleTracking(secId models.SecurityCarIDType) error {
	_ = p.conn.QueryRow(context.Background(),
		`UPDATE security_cars SET security_date_off = NOW() WHERE id = $1`, secId)
	return nil
}

// AddEvent adds an event of the specified type to the security car by secId
func (p *PostgresConnector) AddEvent(secId models.SecurityCarIDType, eventTypeId int64) error {
	_ = p.conn.QueryRow(context.Background(),
		`INSERT INTO events (sc_id, et_id, time)
					VALUES ($1, $2, NOW())`, secId, eventTypeId)
	return nil
}

func (p *PostgresConnector) GetCarsByUserLogin(login string) ([]models.SecurityCar, error) {
	rows, err := p.conn.Query(context.Background(),
		`SELECT sc.id, sc.camera_id, sc.car_id FROM security_cars sc
							 JOIN accounts a ON sc.account_id = a.id
							 WHERE a.login = $1 AND sc.security_date_off IS NULL;`, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	securityCars := make([]models.SecurityCar, 0)
	for rows.Next() {
		var securityCar models.SecurityCar
		err = rows.Scan(&securityCar.ID, &securityCar.CameraID, &securityCar.CarID)
		if err != nil {
			return nil, err
		}
		securityCars = append(securityCars, securityCar)
	}

	return securityCars, nil
}
