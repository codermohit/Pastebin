package models

import (
	"database/sql"
	"errors"
	"time"
)

type Paste struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type PasteModel struct {
	DB *sql.DB
}

func (m *PasteModel) Insert(title, content string, expires int) (int, error) {
	stmt := `INSERT INTO pastes(title , content , created, expires)
  VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *PasteModel) Get(id int) (Paste, error) {
	stmt := `SELECT id, title, content, created, expires FROM pastes
  WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	var p Paste

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Paste{}, ErrNoRecord
		} else {
			return Paste{}, err
		}
	}

	return p, nil
}

func (m *PasteModel) Latest() ([]Paste, error) {
  stmt := `SELECT id, title, content, created, expires FROM pastes
  WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

  rows, err := m.DB.Query(stmt)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var pastes []Paste

  for rows.Next() {
    var p Paste

    err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Created, &p.Expires)
    if err != nil {
      return nil, err
    }

    pastes = append(pastes, p)
  }
  
  if err = rows.Err(); err != nil {
    return nil, err
  }

	return pastes, nil
}
