package main

import "database/sql"

func PrepareDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`create table exchanges(
		id text not null,
		code text not null,
		codein text not null,
		name text not null,
		high text not null, 
		low text not null,
		varbid text not null,
		pctchange text not null,
		bid text not null, 
		ask text not null,
		timestamp text not null,
		createdate text not null
	)`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
