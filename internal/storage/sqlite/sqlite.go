package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

// New создает новый Storage
func New(storagePath string) (*Storage, error) {
	const repo = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", repo, err)
	}

	// Создание таблицы
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
	`)
	if err != nil {
		db.Close() // Закрываем базу при ошибке
		return nil, fmt.Errorf("%s: %w", repo, err)
	}

	// Создание индекса
	_, err = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		db.Close() // Закрываем базу при ошибке
		return nil, fmt.Errorf("%s: %w", repo, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) RedirectByAlias(alias string) (string, error) {
	const repo = "storage.sqlite.Redirect"

	url, err := s.GetUrl(alias)

	if err != nil {
		return "", fmt.Errorf("%s: %w", repo, err)
	}

	return url, nil
}

func (s *Storage) SaveUrl(urlToSave string, alias string) (int64, error) {
	const repo = "storage.sqlite.SaveUrl"

	stmt, err := s.db.Prepare(`
		INSERT INTO url(url,alias) VALUES(?,?)
	`)

	if err != nil {
		return -1, fmt.Errorf("%s: %w", repo, err)
	}

	res, err := stmt.Exec(urlToSave, alias)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %w", repo, storage.ErrUrlExists)
		}

		return -1, fmt.Errorf("%s: %w", repo, err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return -1, fmt.Errorf("%s: failed to get last insert id: %w", repo, err)
	}

	return id, nil
}

func (s *Storage) GetUrl(alias string) (string, error) {
	const repo = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare(`
		SELECT url FROM url WHERE alias = ?
	`)

	if err != nil {
		return "", fmt.Errorf("%s: %w", repo, err)
	}

	var resUrl string

	err = stmt.QueryRow(alias).Scan(&resUrl)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", repo, storage.ErrUrlNotFound)
		}

		return "", fmt.Errorf("%s: %w", repo, err)
	}

	return resUrl, nil
}

func (s *Storage) DeleteUrl(alias string) (int64, error) {
	const repo = "storage.sqlite.DeleteUrl"

	// Сначала получаем id
	var id int64
	err := s.db.QueryRow(`
		SELECT id FROM url WHERE alias = ?
	`, alias).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return -1, fmt.Errorf("%s: %w", repo, storage.ErrUrlNotFound)
		}
		return -1, fmt.Errorf("%s: %w", repo, err)
	}

	// Удаляем запись
	res, err := s.db.Exec(`
		DELETE FROM url WHERE alias = ?
	`, alias)

	if err != nil {
		return -1, fmt.Errorf("%s: %w", repo, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return -1, fmt.Errorf("%s: failed to get rows affected: %w", repo, err)
	}

	if rowsAffected == 0 {
		return -1, fmt.Errorf("%s: %w", repo, storage.ErrUrlNotFound)
	}

	return id, nil
}
