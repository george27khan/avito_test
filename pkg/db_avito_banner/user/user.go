package user

import (
	db "avito_test/pkg/db_avito_banner"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

// User тип для представления записи из таблицы employee
type User struct {
	Id        int64
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Insert функция для добавление записи в таблицу
func (u *User) Insert(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "INSERT INTO avito_banner.user(user_name, created_at, updated_at) VALUES (@user_name, @created_at, @updated_at)"

	args := pgx.NamedArgs{
		"user_name":  u.UserName,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
	if res, err := conn.Exec(ctx, query, args); err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println(res)
	}

	return nil
}

// Delete функция для удаления записи из таблицы
func (u *User) Delete(ctx context.Context) error {
	conn, err := db.PGPool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		return err
	}
	query := "delete from avito_banner.user t where t.user_id = $1"
	_, err = conn.Exec(ctx, query, u.Id)
	if err := conn.Ping(ctx); err != nil {
		return err
	}
	return nil
}
