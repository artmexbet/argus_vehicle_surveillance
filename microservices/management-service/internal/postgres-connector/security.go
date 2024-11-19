package postgres_connector

import (
	"Argus/pkg/models"
	"context"
	"time"
)

func (p *PostgresConnector) GetAllSecuriedCarsByCamera(cameraID models.CameraIDType) ([]models.SecurityCar, error) {
	rows, err := p.conn.Query(context.Background(),
		`SELECT id, camera_id, car_id FROM security_cars
                             WHERE security_date_off IS NULL AND camera_id = $1`, cameraID)
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
	var securityCarID models.SecurityCarIDType
	err := p.conn.QueryRow(context.Background(),
		`INSERT INTO security_cars (car_id, camera_id, account_id, security_date_on) 
					VALUES ($1, $2, $3, $4) RETURNING id`,
		carID, cameraID, accountID, time.Time(securityDateOn)).Scan(&securityCarID)
	return securityCarID, err
}
