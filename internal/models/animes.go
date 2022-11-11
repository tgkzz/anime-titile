package models

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Anime struct {
	ID     int
	Title  string
	Status string
	Owner  string
}

type AnimeModel struct {
	DB *pgxpool.Pool
}

func (m *AnimeModel) Insert(tittle string, status string, owner string) (int, error) {
	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return 0, err
	}
	row := conn.QueryRow(context.Background(), "insert into animes (tittle, status, owner) values ($1, $2, $3)", tittle, status, owner)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *AnimeModel) Get(id int, owner string) (*Anime, error) {
	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	row := conn.QueryRow(context.Background(), "select id, title, status from animes where id = $1 and owner = $2", id, owner)
	s := &Anime{}
	err = row.Scan(&s.ID, &s.Title, &s.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *AnimeModel) All(owner string) ([]*Anime, error) {
	conn, err := m.DB.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	rows, err2 := conn.Query(context.Background(), "select id, title, status from animes where owner = $1", owner)
	defer conn.Release()

	if err2 != nil {
		return nil, err2
	}

	animes := []*Anime{}

	for rows.Next() {
		s := &Anime{}
		err = rows.Scan(&s.ID, &s.Title, &s.Status)
		if err != nil {
			return nil, err
		}
		animes = append(animes, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return animes, err
}
