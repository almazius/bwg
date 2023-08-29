package repository

import (
	"bwg2/config"
	"bwg2/internal/tsfl"
	"bwg2/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"os"
)

type Postgres struct {
	pool *pgxpool.Pool
	log  *zerolog.Logger
}

func NewRepository(ctx context.Context, conf *config.Config) (*Postgres, error) {
	logg := zerolog.New(os.Stderr)

	pool, err := utils.GetPool(ctx, conf)
	if err != nil {
		logg.Error().Timestamp().Err(err).Send()
		return nil, &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}
	p := &Postgres{
		pool: pool,
		log:  &logg,
	}
	return p, nil
}

func (p *Postgres) ApplicationInputMoney(ctx context.Context, userId uuid.UUID, count uint64) (string, error) {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		return "", &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}

	var transactionId uuid.UUID

	err = conn.QueryRow(ctx, `insert into transactions
    (transactionId, userId, count, input, confirmed, info) values  ($1, $2, $3, $4,$5, $6) returning transactionId`,
		uuid.New(), userId, count, true, false, fmt.Sprintf("Refill balance on %d penny", count)).Scan(&transactionId)
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return "", &r
	}

	return transactionId.String(), nil
}

func (p *Postgres) InputMoney(ctx context.Context, userId uuid.UUID, count uint64) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		return &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}
	_, err = conn.Exec(ctx, `update users set balance = balance + $2  where userid = $1 `, userId, count)
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &r
	}

	return nil
}

func (p *Postgres) GetInfoFromTransactions(ctx context.Context, transactionId uuid.UUID) (*tsfl.TransactionInfo, error) {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		return nil, &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}

	// I don't know, how scan in struct with pgxpool :(
	var (
		TransactionId uuid.UUID
		UserId        uuid.UUID
		Count         uint64
		Input         bool
		Confirmed     bool
		Info          string
	)

	err = conn.QueryRow(ctx, `select * from transactions where transactionid = $1`, transactionId).Scan(
		&TransactionId, &UserId, &Count, &Input, &Confirmed, &Info)
	if err != nil {
		p.log.Error().Timestamp().Err(err).Msg("AAA")
		return nil, &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}

	tInfo := &tsfl.TransactionInfo{TransactionId: TransactionId, UserId: UserId,
		Count: Count, Input: Input, Confirmed: Confirmed, Info: Info}

	return tInfo, err
}

func (p *Postgres) IsExistUser(ctx context.Context, userId uuid.UUID) (bool, error) {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		return false, &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}
	var result bool

	err = conn.QueryRow(ctx, `select exists(select * from users where userid=$1)`, userId).Scan(&result)
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return false, &r
	}

	return result, nil
}

func (p *Postgres) CreateUser(ctx context.Context, userId uuid.UUID) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &r
	}

	_, err = conn.Exec(ctx, `insert into users (userid, balance) VALUES ($1, 0)`, userId)
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &r
	}

	return nil
}

func (p *Postgres) GetBalance(ctx context.Context, userId uuid.UUID) (uint64, error) {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return 0, &r
	}

	var balance uint64
	err = conn.QueryRow(ctx, `select balance from users where userid=$1`, userId).Scan(&balance)

	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return 0, &r
	}

	return balance, nil
}

func (p *Postgres) UpdateTransactionsConfirmed(ctx context.Context, transactionId uuid.UUID, condition bool) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &r
	}

	_, err = conn.Exec(ctx, `update transactions set confirmed = $2 where transactionid = $1`, transactionId, condition)
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &r
	}
	return nil
}

func (p *Postgres) ApplicationOutputMoney(ctx context.Context, userId uuid.UUID, count uint64) (string, error) {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		return "", &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}

	var transactionId uuid.UUID

	err = conn.QueryRow(ctx, `insert into transactions
    (transactionId, userId, count, input, confirmed, info) values  ($1, $2, $3, $4,$5, $6) returning transactionId`,
		uuid.New(), userId, count, false, false, fmt.Sprintf("Relief balance on %d penny", count)).Scan(&transactionId)
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return "", &r
	}

	return transactionId.String(), nil
}

func (p *Postgres) OutputMoney(ctx context.Context, userId uuid.UUID, count uint64) error {
	conn, err := p.pool.Acquire(ctx)
	defer conn.Release()
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		return &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}
	_, err = conn.Exec(ctx, `update users set balance = balance - $2  where userid = $1 `, userId, count)
	if err != nil {
		p.log.Error().Timestamp().Err(err).Send()
		r := tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
		return &r
	}

	return nil
}
