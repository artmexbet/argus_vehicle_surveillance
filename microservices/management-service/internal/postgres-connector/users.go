package postgres_connector

import (
	"Argus/pkg/models"
	"context"
)

// GetAccountIdByLogin returns id of account with given login string
func (p *PostgresConnector) GetAccountIdByLogin(login string) (models.AccountIDType, error) {
	var id models.AccountIDType
	err := p.conn.QueryRow(context.Background(), `SELECT id FROM public.accounts 
          WHERE login=$1;`, login).Scan(&id)
	return id, err
}
