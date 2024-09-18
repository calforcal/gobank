package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface{
	CreateAccount(*Account) error
	DeleteAccount(int) (bool, error)
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountById(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=michaelcallahan dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts (
		id serial primary key,
		first_name varchar(55),
		last_name varchar(55),
		number bigint,
		balance decimal,
		created_at timestamp default NOW()
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := 
		`
		INSERT INTO accounts
		(first_name, last_name, number, balance, created_at)
		VALUES ($1, $2, $3, $4, $5)
		`

	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id int) (bool, error) {
	query := 
	`
	DELETE FROM accounts 
	WHERE id = $1
	`
	
	res, err := s.db.Exec(query, id)

	if err != nil {
		return false, err
	}

	rowsAffected, err := res.RowsAffected()

	if rowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (s *PostgresStore) GetAccountById(id int) (*Account, error) {
	query:= 
	`
	SELECT * FROM accounts
	WHERE id = $1;
	`
	row := s.db.QueryRow(query, id)

	if row != nil {
		account := &Account{}
		err := row.Scan(
			&account.ID, 
			&account.FirstName, 
			&account.LastName, 
			&account.Number,
			&account.Balance, 
			&account.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		return account, nil
	}
	return nil, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID, 
			&account.FirstName, 
			&account.LastName, 
			&account.Number,
			&account.Balance, 
			&account.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
