package repository

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/elysium-challenge/sessions/domain"
	"github.com/jackc/pgx/v5"
)

type PGUserRepository struct {
	conn *pgx.Conn
}

func NewPostgresUserRepository(ctx context.Context, connString string) (*PGUserRepository, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PGUserRepository{conn}, nil
}

func (p *PGUserRepository) Create(user *domain.User) error {
	query := `insert into users(first_name, last_name, email, password) 
	values (@firstName, @lastName, @email, @password, @jwts) returning id`
	args := pgx.NamedArgs{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"password":  user.Password,
		"jwts":      user.JWTs,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&user.ID)
	return err
}

func (p *PGUserRepository) Update(user *domain.User) error {
	query := `update users set first_name=@first_name, last_name=@last_name, 
		email=@email, password=@password, @jwts=@jwts where id=@id`
	args := pgx.NamedArgs{
		"id":         user.ID,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"password":   user.Password,
		"jwts":       user.JWTs,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return err
}

func (p *PGUserRepository) GetByID(id int) (domain.User, error) {
	query := `select * from users where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[domain.User])

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (p *PGUserRepository) GetByEmail(email string) (domain.User, error) {
	query := `select * from users where email=@email`
	args := pgx.NamedArgs{
		"email": email,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	user, err := pgx.CollectOneRow(row, pgx.RowToStructByName[domain.User])

	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
