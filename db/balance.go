package db

import (
	"database/sql"

	"github.com/volssan/balances/models"
)

func (db Database) AddBalance(balance *models.Balance) error {
	query := `INSERT INTO user_balances (user_id, balance) VALUES ($1, $2) RETURNING created_at, updated_at`
	err := db.Conn.QueryRow(query, balance.UserId, balance.Balance).Scan(&balance.CreatedAt, &balance.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) UpdateBalance(balance *models.Balance) error {
	query := `UPDATE user_balances SET balance=$2 WHERE user_id=$1 RETURNING balance, updated_at;`
	err := db.Conn.QueryRow(query, balance.UserId, balance.Balance).Scan(&balance.Balance, &balance.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoMatch
		}
		return err
	}

	return nil
}

func (db Database) GetBalanceByUserId(userId int) (models.Balance, error) {
	balance := models.Balance{}

	query := `SELECT * FROM user_balances WHERE user_id = $1;`
	err := db.Conn.QueryRow(query, userId).Scan(&balance.UserId, &balance.Balance, &balance.CreatedAt, &balance.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return balance, ErrNoMatch
		}
		return balance, err
	}

	return balance, nil
}