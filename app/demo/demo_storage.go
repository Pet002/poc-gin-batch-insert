package demo

import (
	"database/sql"
)

type Demo struct {
	ID      int
	Name    string
	Surname string
	Age     int
}

type Detail struct {
	ID     int64
	DemoID int
	Detail string
}

func (s *Storage) InsertToDemo(tx *sql.Tx, demo *Demo) error {
	res, err := tx.Exec(`
		INSERT INTO demo ( `+"`name`,"+"`surname`,"+"`age`)"+`VALUES
		(?, ?, ?);
	`, demo.Name, demo.Surname, demo.Age)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	demo.ID = int(id)
	return nil
}

func (s *Storage) InsertToDetail(tx *sql.Tx, detail *Detail) error {
	res, err := tx.Exec(`
		INSERT INTO detail ( `+"`demo_id`,"+"`detail`)"+`VALUES
		(?, ?);
	`, detail.DemoID, detail.Detail)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	detail.ID = id
	return nil
}
