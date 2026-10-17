package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/yonesbeheshti/ai-gateway/entity"
)

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	DBName   string `koanf:"db_name"`
	Sslmode  string `koanf:"sslmode"`
}

type PostgresDB struct {
	cfg Config
	db  *sql.DB
}

func New(config Config) *PostgresDB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", config.Username, config.Password, config.Host, config.Port, config.DBName, config.Sslmode)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("open db error:", err)
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	db.SetConnMaxIdleTime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &PostgresDB{
		db:  db,
		cfg: config,
	}
}

func (p *PostgresDB) Conn() *sql.DB {
	return p.db
}

func (p *PostgresDB) AddUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	const op = "AddUser"
	err := p.db.QueryRowContext(ctx, "INSERT INTO users(email, api_key, token_balance) values($1, $2, $3) returning id;", user.Email, user.ApiKey, user.TokenBalance).Scan(&user.ID)
	if err != nil {
		log.Printf("%v: failed to insert user into database, err %v\n", op, err)
		return nil, err
	}

	log.Printf("%v: inserted user into database\n", op)

	return user, nil
}

func (p *PostgresDB) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const op = "GetUserByEmail"
	var user entity.User

	err := p.db.QueryRowContext(ctx, "SELECT * FROM users WHERE email=$1", email).Scan(&user.ID, &user.Email, &user.ApiKey, &user.TokenBalance, &user.CreatedAt)
	if err != nil {

		if err == sql.ErrNoRows {
			log.Printf("%v: user not found\n", op)
			return nil, nil
		}
		log.Printf("%v: failed to get user by email, err: %v\n", op, err)
		return nil, err
	}

	return &user, nil
}

func (p *PostgresDB) GetUserByApiKey(ctx context.Context, apiKey string) (*entity.User, error) {

	const op = "GetUserByApiKey"
	var user entity.User
	err := p.db.QueryRowContext(ctx,
		"SELECT id, email, api_key, token_balance, created_at FROM users WHERE api_key=$1",
		apiKey,
	).Scan(&user.ID, &user.Email, &user.ApiKey, &user.TokenBalance, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (p *PostgresDB) DeductUserTokens(ctx context.Context, userID uint64, tokens int64) (int64, error) {
	const query = `
		UPDATE users
		SET token_balance = GREATEST(0, token_balance - $1)
		WHERE id = $2
		RETURNING token_balance;
	`
	var remainingBalance int64
	err := p.db.QueryRowContext(ctx, query, tokens, userID).Scan(&remainingBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to deduct tokens: %w", err)
	}
	return remainingBalance, nil
}
