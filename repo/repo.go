package repo

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Thevtok/auth/db"
	"github.com/Thevtok/auth/model"
	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	dbHost := db.DotEnv("DB_HOST")
	dbPort := db.DotEnv("DB_PORT")
	dbUser := db.DotEnv("DB_USER")
	dbPassword := db.DotEnv("DB_PASSWORD")
	dbName := db.DotEnv("DB_NAME")
	sslMode := db.DotEnv("SSL_MODE")
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)
	db, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	} else {
		log.Println("Database Successfully Connected")
	}
	return db
}

type LoginRepo interface {
	GetByUsernameAndPassword(username string, password string) (*model.User, error)
}

type loginRepo struct {
	db *sql.DB
}

func (r *loginRepo) GetByUsernameAndPassword(username string, password string) (*model.User, error) {
	query := "SELECT c.username, c.password, s.c_username FROM credentials c JOIN students s ON c.username = s.c_username WHERE c.username = $1 AND c.password = $2"
	row := r.db.QueryRow(query, username, password)

	user := &model.User{}
	err := row.Scan(&user.Username, &user.Password, &user.C_Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}

	return user, nil
}

func NewStudentRepo(db *sql.DB) LoginRepo {
	repo := new(loginRepo)
	repo.db = db

	return repo
}
