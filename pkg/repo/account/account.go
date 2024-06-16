package account

import (
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/domain/account"
	interfaces "awesomeProject/pkg/repo/account/interface"
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"log"
	"time"

	"github.com/jackc/pgconn"
)

type accountDataBase struct {
	db *sql.DB
}

func NewAccountDataBase(db *sql.DB) interfaces.AccountRepository {
	return &accountDataBase{
		db: db,
	}
}

func (r *accountDataBase) Migrate(ctx context.Context) error {
	accQuery := `
		CREATE TABLE IF NOT EXISTS account(
			id SERIAL PRIMARY KEY,
			firstName text NOT NULL,
			lastName text NOT NULL,
			birthDate date NOT NULL,
			email text NOT NULL,
			password text NOT NULL,
			subscribersIds integer[]
		);
	`
	_, err := r.db.ExecContext(ctx, accQuery)
	if err != nil {
		message := db.ErrMigrate.Error() + " account"
		log.Printf("%q: %s\n", message, err.Error())
		return db.ErrMigrate
	}

	return err
}

func (r *accountDataBase) Registration(ctx context.Context, newAccount account.Registration) (*account.Info, error) {
	var existingCount int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM account WHERE email = $1", newAccount.Email).Scan(&existingCount)
	if err != nil {
		return nil, err
	}

	if existingCount > 0 {
		return nil, db.ErrDuplicate
	}
	birthDateStr := newAccount.BirthDate.Format("2006-01-02")

	var id int64
	err = r.db.QueryRowContext(ctx,
		"INSERT INTO account(firstName, lastName, email, password, birthDate, subscribersIds) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		newAccount.FirstName, newAccount.LastName, newAccount.Email, newAccount.Password, birthDateStr, pq.Array([]int64{})).Scan(&id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	requestAccount := &account.Info{
		ID:          id,
		FirstName:   newAccount.FirstName,
		LastName:    newAccount.LastName,
		Email:       newAccount.Email,
		BirthDate:   newAccount.BirthDate,
		Subscribers: []int64{},
	}

	return requestAccount, nil
}

func (r *accountDataBase) Login(ctx context.Context, acc account.Login) (int64, error) {
	var id int64
	row := r.db.QueryRowContext(ctx, "SELECT id FROM account WHERE email = $1 and password = $2", acc.Email, acc.Password)

	if err := row.Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, db.ErrNotExist
		}
		return 0, err
	}
	return id, nil
}

func (r *accountDataBase) Put(ctx context.Context, id int64, updateAcc *account.Account) (*account.Info, error) {
	birthDateStr := updateAcc.BirthDate.Format("02.01.2006")
	res, err := r.db.ExecContext(ctx, "UPDATE account SET firstName = $1, lastName = $2, email = $3, password = $4, birthDate = $5, subscribersIds = $6 WHERE id = $7",
		updateAcc.FirstName, updateAcc.LastName, updateAcc.Email, updateAcc.Password, birthDateStr, pq.Array(updateAcc.SubscribersIds), id)
	if err != nil {
		var pgxError *pgconn.PgError
		if errors.As(err, &pgxError) {
			if pgxError.Code == "23505" {
				return nil, db.ErrDuplicate
			}
		}
		return nil, err
	}

	result := &account.Info{
		ID:          id,
		FirstName:   updateAcc.FirstName,
		LastName:    updateAcc.LastName,
		Email:       updateAcc.Email,
		BirthDate:   updateAcc.BirthDate,
		Subscribers: updateAcc.SubscribersIds,
	}
	birthDate, err := time.Parse("02.01.2006", birthDateStr)
	if err != nil {
		return nil, err
	}
	result.BirthDate = birthDate

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, db.ErrUpdateFailed
	}

	return result, nil
}
func (r *accountDataBase) Get(ctx context.Context, id int64) (*account.Account, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, firstName, lastName, birthDate, email, password, subscribersIds FROM account WHERE id = $1", id)

	var result account.Account
	if err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.BirthDate, &result.Email, &result.Password, pq.Array(&result.SubscribersIds)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.ErrNotExist
		}
		return nil, err
	}

	return &result, nil
}

func (r *accountDataBase) GetAll(ctx context.Context) ([]account.Account, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, first_name, last_name, birth_date, email, password, subscribers_ids FROM account")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []account.Account

	for rows.Next() {
		var account account.Account
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.BirthDate,
			&account.Email,
			&account.Password,
			&account.SubscribersIds,
		)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *accountDataBase) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM account WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return db.ErrDeleteFailed
	}

	return err
}
