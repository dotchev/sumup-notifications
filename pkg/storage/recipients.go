package storage

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sumup-notifications/pkg/model"
)

type Recipients struct {
	DB *pgxpool.Pool
}

func (s Recipients) Load(ctx context.Context, recipient string) (model.RecipientContact, error) {
	var contact model.RecipientContact

	var contactJSON []byte
	err := s.DB.QueryRow(ctx,
		"SELECT contact FROM recipients WHERE recipient = $1", recipient).Scan(&contactJSON)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return contact, ErrNotFound
		}
		return contact, err
	}

	err = json.Unmarshal(contactJSON, &contact)
	return contact, err
}

// Store stores a new recipient or updates an existing one
func (s Recipients) Store(ctx context.Context, recipient string, contact model.RecipientContact) error {
	contactJSON, err := json.Marshal(contact)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(ctx,
		"INSERT INTO recipients (recipient, contact) VALUES ($1, $2) ON CONFLICT (recipient) DO UPDATE SET contact = $2",
		recipient, contactJSON)
	return err
}

func (s Recipients) Delete(ctx context.Context, recipient string) error {
	_, err := s.DB.Exec(ctx,
		"DELETE FROM recipients WHERE recipient = $1", recipient)
	return err
}
