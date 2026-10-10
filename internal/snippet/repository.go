package snippet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &sqliteRepository{
		db: db,
	}
}

func (s *sqliteRepository) Create(ctx context.Context, snippet *Snippet) error {
	query := `INSERT INTO snippets (user_id, title, content, is_private, created_at, updated_at) VALUES (?,?,?,?,?,?)`

	result, err := s.db.ExecContext(ctx, query, snippet.UserId, snippet.Title, snippet.Content, snippet.IsPrivate, snippet.CreatedAt, snippet.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert snippet %w", err)
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return fmt.Errorf("failed to last insert id %w", err)
	}

	snippet.ID = lastId

	return nil
}

func (s *sqliteRepository) GetById(ctx context.Context, id int64, userID int64) (*Snippet, error) {
	query := `SELECT id, user_id, title, content, is_private, created_at, updated_at FROM snippets WHERE id = ? AND user_id = ?`

	row := s.db.QueryRowContext(ctx, query, id, userID)

	var sni Snippet

	err := row.Scan(
		&sni.ID,
		&sni.UserId,
		&sni.Title,
		&sni.Content,
		&sni.IsPrivate,
		&sni.CreatedAt,
		&sni.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSnippetNotFound
		}
		return nil, fmt.Errorf("failed to get snippet: %w", err)
	}

	return &sni, nil
}

func (s *sqliteRepository) ListByUserID(ctx context.Context, userID int64) ([]Snippet, error) {
	query := `SELECT id, user_id, title, content, is_private, created_at, updated_at FROM snippets WHERE user_id = ? ORDER BY created_at DESC`

	rows, err := s.db.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to query snippets: %w", err)
	}

	defer rows.Close()

	snippets := make([]Snippet, 0)

	for rows.Next() {
		var sni Snippet
		err := rows.Scan(
			&sni.ID,
			&sni.UserId,
			&sni.Title,
			&sni.Content,
			&sni.IsPrivate,
			&sni.CreatedAt,
			&sni.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan snippet: %w", err)
		}

		snippets = append(snippets, sni)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return snippets, nil
}

func (s *sqliteRepository) Update(ctx context.Context, snippet *Snippet) error {
	query := `UPDATE snippets SET title = ?, content = ?, is_private = ?, updated_at = ? WHERE id = ? AND user_id = ?`

	res, err := s.db.ExecContext(ctx, query, snippet.Title, snippet.Content, snippet.IsPrivate, snippet.UpdatedAt, snippet.ID, snippet.UserId)

	if err != nil {
		return fmt.Errorf("failed to update snippet %w", err)
	}

	rows, err := res.RowsAffected()

	if rows == 0 {
		return ErrSnippetNotFound
	}

	return nil
}

func (s *sqliteRepository) Delete(ctx context.Context, id int64, userID int64) error {
	query := `DELETE FROM snippets WHERE id = ? AND user_id = ?`

	res, err := s.db.ExecContext(ctx, query, id, userID)

	if err != nil {
		return fmt.Errorf("failed to delete snippet %w", err)
	}

	rows, err := res.RowsAffected()

	if rows == 0 {
		return ErrSnippetNotFound
	}

	return nil
}
