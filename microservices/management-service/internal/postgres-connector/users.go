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

// GetTelegramId returns telegram id of user attached to security car
func (p *PostgresConnector) GetTelegramId(secId models.SecurityCarIDType) (int64, error) {
	var telegramId int64
	err := p.conn.QueryRow(
		context.Background(),
		`SELECT telegram_id FROM public.profile p 
    JOIN public.accounts a on p.account_id = a.id
    JOIN public.security_cars sc on a.id = sc.account_id
    WHERE sc.id = $1`,
		secId,
	).Scan(&telegramId)
	return telegramId, err
}

// CheckHasUserTelegramId returns true if id specified to user, else returns false
func (p *PostgresConnector) CheckHasUserTelegramId(accountId models.AccountIDType) (bool, error) {
	var hasTelegramId bool
	err := p.conn.QueryRow(
		context.Background(),
		`SELECT telegram_id IS NOT NULL FROM public.profile p 
		WHERE p.account_id = $1`,
		accountId,
	).Scan(&hasTelegramId)
	return hasTelegramId, err
}
