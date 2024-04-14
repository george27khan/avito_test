package user

import (
	db "avito_test/pkg/postgres_db/connection"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User тип для представления записи из таблицы employee
type User struct {
	UserId    int64
	UserName  string
	Password  string
	IsAdmin   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Insert функция для добавление записи в таблицу
func (u *User) Insert(ctx context.Context, conn *pgxpool.Conn) error {
	query := "INSERT INTO avito_banner.user(user_name, password, is_admin) VALUES (@user_name, @password, @is_admin)"
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	u.Password = string(hashPwd)
	if err != nil {
		panic(err)
	}
	args := pgx.NamedArgs{
		"user_name": u.UserName,
		"password":  u.Password,
		"is_admin":  u.IsAdmin,
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
	if err != nil {
		return err
	}
	defer conn.Release()
	query := "delete from avito_banner.user t where t.user_id = $1"
	_, err = conn.Exec(ctx, query, u.UserId)
	if err := conn.Ping(ctx); err != nil {
		return err
	}
	return nil
}

// Get функция для поиска пользователя
func Get(ctx context.Context, userName string) (User, error) {
	var u User
	conn, err := db.PGPool.Acquire(ctx)
	if err != nil {
		return u, err
	}
	defer conn.Release()

	query := "select user_id, user_name, password, is_admin, created_at, updated_at from avito_banner.user t where t.user_name = $1"
	if err := conn.QueryRow(ctx, query, userName).Scan(&u.UserId, &u.UserName, &u.Password, &u.IsAdmin, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return u, err
	}
	return u, nil
}

// VerifyPassword функция для проверки пароля
func (u *User) VerifyPassword(password string) (ok bool) {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return
	}
	ok = true
	return
}
