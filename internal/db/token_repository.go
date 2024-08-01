package db

import (
	"context"
	"fmt"
	"time"
)

func (r *Repository) CleanTokens(ctx context.Context) error {
	currentTime := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, "DELETE FROM tokens WHERE expires_at < $1", currentTime)
	if err != nil {
		return fmt.Errorf("failed cleaning tokens: %w", err)
	}

	return nil
}
