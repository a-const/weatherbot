package env

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type User struct {
	ID            int    `db:"id"`
	Username      string `db:"username"`
	Notifications bool   `db:"notifications"`
}

type PostgresDB struct {
	db *sqlx.DB
}

func (pgdb *PostgresDB) NewPostgreDB(cfg Config) error {
	var err error
	pgdb.db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return err
	}
	err = pgdb.db.Ping()

	if err != nil {
		return err
	}

	return nil
}

func (pgdb *PostgresDB) CreateUser(username string) error {
	query := fmt.Sprintf("INSERT INTO users (username, notifications) values ('%s', %t)", username, false)
	row := pgdb.db.QueryRow(query)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (pgdb *PostgresDB) EditNotofications(username string, notifications bool) error {
	query := fmt.Sprintf("UPDATE users SET notifications=%t WHERE username='%s'", notifications, username)
	row := pgdb.db.QueryRow(query)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (pgdb *PostgresDB) GetUsers() (*[]User, error) {
	users := &[]User{}
	err := pgdb.db.Select(users, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	//query := fmt.Sprint("SELECT (username, notifications) FROM users",id)
	//row := pgdb.db.QueryRow(query)
	//err := pgdb.db.Get()
	// if err := row.Err(); err != nil {
	// 	return err
	// }
	return users, nil
}
