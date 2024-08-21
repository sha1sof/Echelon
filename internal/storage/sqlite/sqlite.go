package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sha1sof/Echelon-/internal/storage"
	"time"
)

type Storage struct {
	db *sql.DB
}

// New конструктор БД.
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite3.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// URLSaver сохраняет ссылки на видео и хеш.
func (s *Storage) URLSaver(ctx context.Context, videoID string, hash []byte, lifetime time.Duration) (uid int64, err error) {
	const op = "storage.sqlite3.URLSaver"

	stmt, err := s.db.Prepare("INSERT INTO video (video_id, hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, videoID, hash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrVideoIDExist)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

// Url для получения данных
func (s *Storage) Url(ctx context.Context, videoID string) (cache []byte, err error) {
	const op = "storage.sqlite3.Url"

	stmt, err := s.db.Prepare("SELECT hash FROM video WHERE video_id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, videoID)

	var hash []byte
	err = row.Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return hash, fmt.Errorf("%s: %w", op, storage.ErrVideoIDNotFound)
		}

		return hash, fmt.Errorf("%s: %w", op, err)
	}

	return hash, nil
}
