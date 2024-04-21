package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID
	URL         string
	BackHalf    string
	Engagements int64
	UserID      sql.Null[uuid.UUID]
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type LinkModel struct {
	DB *sql.DB
}

func (m *LinkModel) Insert(url, backHalf string, userID string) (*Link, error) {
	var u uuid.NullUUID
	var err error

	if userID == "" {
		err = u.Scan(nil)
	} else {
		err = u.Scan(userID)
	}
	if err != nil {
		return nil, err
	}

	uID, err := u.Value()
	if err != nil {
		return nil, err
	}

	stmt := `INSERT INTO links (id, url, back_half, user_id, created_at, updated_at)
					 VALUES (?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`
	_, err = m.DB.Exec(stmt, uuid.New(), url, backHalf, uID)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "links_uc_back_half") {
				return nil, ErrDuplicateBackHalf
			}
		}

		return nil, err
	}

	link := &Link{}
	stmt = "SELECT id, url, back_half, user_id, created_at, updated_at FROM links ORDER BY created_at DESC LIMIT 1"
	err = m.DB.QueryRow(stmt).Scan(&link.ID, &link.URL, &link.BackHalf, &link.UserID, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return link, nil
}

func (m *LinkModel) GetByID(id string) (*Link, error) {
	link := &Link{}

	stmt := "SELECT * FROM links WHERE id = ?"
	err := m.DB.QueryRow(stmt, id).Scan(&link.ID, &link.URL, &link.BackHalf, &link.Engagements, &link.UserID, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return link, nil
}

func (m *LinkModel) GetByBackHalf(backHalf string) (*Link, error) {
	link := &Link{}

	stmt := "SELECT * FROM links WHERE back_half = ?"
	err := m.DB.QueryRow(stmt, backHalf).Scan(&link.ID, &link.URL, &link.BackHalf, &link.Engagements, &link.UserID, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return link, nil
}

func (m *LinkModel) GetByUserID(userID string) ([]*Link, error) {
	stmt := "SELECT * FROM links WHERE user_id = ? ORDER BY created_at DESC"
	rows, err := m.DB.Query(stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := []*Link{}

	for rows.Next() {
		link := &Link{}
		err := rows.Scan(&link.ID, &link.URL, &link.BackHalf, &link.Engagements, &link.UserID, &link.CreatedAt, &link.UpdatedAt)
		if err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (m *LinkModel) BackHalfExists(backHalf string) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM links WHERE back_half = ?)"
	err := m.DB.QueryRow(stmt, backHalf).Scan(&exists)

	return exists, err
}

func (m *LinkModel) Update(id, url, backHalf string) error {
	stmt := "UPDATE links SET url = ?, back_half = ?, updated_at = UTC_TIMESTAMP() WHERE id = ?"
	_, err := m.DB.Exec(stmt, url, backHalf, id)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "links_uc_back_half") {
				return ErrDuplicateBackHalf
			}
		}

		return err
	}

	return nil
}

func (m *LinkModel) UpdateEngagements(id string, engagements int) error {
	stmt := "UPDATE links SET engagements = ? WHERE id = ?"
	_, err := m.DB.Exec(stmt, engagements, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *LinkModel) Delete(id string) error {
	stmt := "DELETE FROM links WHERE id = ?"
	_, err := m.DB.Exec(stmt, id)
	return err
}
