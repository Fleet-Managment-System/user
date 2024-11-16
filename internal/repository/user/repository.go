package repository

import (
	"context"
	"time"
	"user/internal/domain/models"
	"user/internal/repository"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	tableName = "\"users\""

	idColumn           = "id"
	firstnameColumn    = "firstname"
	lastnameColumn     = "lastname"
	emailColumn        = "email"
	passwordHashColumn = "passwordHash"
	createdAtColumn    = "createdAt"
	updatedAtColumn    = "updatedAt"
)

type UserRepository struct{}

const dbDSN = "host=localhost port=1234 dbname=project user=spuser password=SPuser96 sslmode=disable"

func (repo *UserRepository) Create(ctx context.Context, user *models.CreateUserModel) (int64, error) {
	pool, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		return -1, status.Errorf(codes.Internal, "Failed to connect to database %s", err)
	}

	insertBuilder := squirrel.Insert(tableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(firstnameColumn, lastnameColumn, emailColumn, passwordHashColumn, createdAtColumn, updatedAtColumn).
		Values(user.FirstName, user.LastName, user.Email, user.PasswordHash, user.CreatedAt, user.UpdateAt).
		Suffix("RETURNING id")

	query, args, err := insertBuilder.ToSql()
	if err != nil {
		return 0, status.Errorf(codes.Internal, "Failed to build query: %s", err)
	}

	var id int64
	err = pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return -1, status.Errorf(codes.Internal, "Failed to proceed query: %s", err)
	}

	return id, nil
}

func (repo *UserRepository) Update(ctx context.Context, userId int64, user *models.UpdateUserModel) error {
	pool, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to connect to database %s", err)
	}

	updateBuilder := squirrel.Update(tableName).
		PlaceholderFormat(squirrel.Dollar).
		Set(updatedAtColumn, user.UpdatedAt)

	if user.FirstName != "" {
		updateBuilder = updateBuilder.Set(firstnameColumn, user.FirstName)
	}
	if user.LastName != "" {
		updateBuilder = updateBuilder.Set(lastnameColumn, user.LastName)
	}

	updateBuilder = updateBuilder.Where(squirrel.Eq{idColumn: userId})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to build query: %s", err)
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to proceed query: %s", err)
	}

	return nil
}
func (repo *UserRepository) GetOne(ctx context.Context, userId int64) (*models.UserModel, error) {
	pool, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to connect to database %s", err)
	}

	selectBuilder := squirrel.Select(idColumn, emailColumn, firstnameColumn, lastnameColumn, passwordHashColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(squirrel.Dollar).
		From(tableName).
		Where(squirrel.Eq{idColumn: userId})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to build query: %s", err)
	}

	var id int64
	var firstname string
	var lastname string
	var email string
	var passwordHash string
	var createdAt time.Time
	var updatedAt *time.Time
	err = pool.QueryRow(ctx, query, args...).
		Scan(&id, &email, &firstname, &lastname, &passwordHash, &createdAt, &updatedAt)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to proceed query: %s", err)
	}

	return &models.UserModel{
		ID:           id,
		FirstName:    firstname,
		LastName:     lastname,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    createdAt,
		UpdateAt:     updatedAt,
	}, nil
}

func (repo *UserRepository) GetAll(ctx context.Context, pagination *repository.Pagination) ([]*models.UserModel, error) {
	pool, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to connect to database %s", err)
	}

	var limit uint64
	var offset uint64
	if pagination != nil {
		limit = pagination.Limit
		offset = pagination.Offset
	} else {
		limit = 10
		offset = 0
	}

	selectBuilder := squirrel.
		Select(
			idColumn,
			firstnameColumn,
			lastnameColumn,
			emailColumn,
			passwordHashColumn,
			createdAtColumn,
			updatedAtColumn).
		From(tableName).
		Limit(limit).
		Offset(offset)
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to build query: %s", err)
	}

	rows, err := pool.Query(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to proceed query: %s", err)
	}

	defer rows.Close()

	var users []*models.UserModel
	for rows.Next() {
		var user models.UserModel

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PasswordHash,
			&user.CreatedAt,
			&user.UpdateAt,
		); err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to proceed query: %s", err)
		}

		users = append(users, &user)
	}

	return users, nil
}

func (repo *UserRepository) Delete(ctx context.Context, userId int64) error {
	pool, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to connect to database %s", err)
	}

	selectBuilder := squirrel.Delete(tableName).Where("id = ?", userId)
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to build query: %s", err)
	}

	_, err = pool.Exec(ctx, query, args...)
	if err != nil {
		return status.Errorf(codes.Internal, "Failed to proceed query: %s", err)
	}

	return nil
}
