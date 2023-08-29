package usecase

import (
	"bwg2/internal/tsfl"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"os"
	"time"
)

// url for send form for input money
var url = "127.0.0.1:8080"

type repository interface {
	InputMoney(ctx context.Context, userId uuid.UUID, count uint64) error
	CreateUser(ctx context.Context, userId uuid.UUID) error
	IsExistUser(ctx context.Context, userId uuid.UUID) (bool, error)
	GetBalance(ctx context.Context, userId uuid.UUID) (uint64, error)
	ApplicationInputMoney(ctx context.Context, userId uuid.UUID, count uint64) (string, error)
	GetInfoFromTransactions(ctx context.Context, transactionId uuid.UUID) (*tsfl.TransactionInfo, error)
	UpdateTransactionsConfirmed(ctx context.Context, transactionId uuid.UUID, condition bool) error
	ApplicationOutputMoney(ctx context.Context, userId uuid.UUID, count uint64) (string, error)
	OutputMoney(ctx context.Context, userId uuid.UUID, count uint64) error
}

type Service struct {
	log *zerolog.Logger
	repository
}

func NewService(rep repository) *Service {
	logg := zerolog.New(os.Stderr)
	return &Service{
		log:        &logg,
		repository: rep,
	}
}

func (s *Service) ApplicationInputMoney(ctx context.Context, userId uuid.UUID, count uint64) (string, error) {
	// check what user is existed
	isExist, err := s.repository.IsExistUser(ctx, userId)
	if err != nil {
		return "", err
	}
	if !isExist {
		err = s.repository.CreateUser(ctx, userId)
	}

	// create transaction for input money
	orderId, err := s.repository.ApplicationInputMoney(ctx, userId, count)
	if err != nil {
		s.log.Error().Timestamp().Err(err).Send()
		return "", err
	}
	s.log.Info().Timestamp().Msg(fmt.Sprintf("User [%s] want input money [%d]\n", userId.String(), count))

	// Any status for application for replenishment of the balance
	//return "", &tsfl.RespondError{
	//	Status:  403,
	//	Message: "You can't refill money",
	//	Err:     errors.New("unable to top up your balance"),
	//}

	// return link with input form
	return fmt.Sprintf("%s/%s", url, orderId), nil
}

func (s *Service) InputMoney(ctx context.Context, transactionId uuid.UUID) error {
	// check transaction
	info, err := s.repository.GetInfoFromTransactions(ctx, transactionId)
	if err != nil {
		return err
	}
	// input money in balance
	err = s.repository.InputMoney(ctx, info.UserId, info.Count)
	if err != nil {
		return err
	}

	// update transaction(confirmed) to true
	err = s.repository.UpdateTransactionsConfirmed(ctx, transactionId, true)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ApplicationOutputMoney(ctx context.Context, userId uuid.UUID, count uint64, url string) error {
	// I would like to send a withdrawal request to Rabbitmq, but I don't have the strength to connect it, sorry:(
	go s.tryOutMoney(ctx, userId, count, url)
	s.log.Info().Timestamp().Msg(fmt.Sprintf("User [%s] want output money [%d]\n", userId.String(), count))
	return nil
}

func (s *Service) tryOutMoney(ctx context.Context, userId uuid.UUID, count uint64, url string) {
	// get transactionId
	tempTransactionId, err := s.repository.ApplicationOutputMoney(ctx, userId, count)
	if err != nil {
		s.log.Error().Timestamp().Err(err).Msg("can't output money")
		return
	}
	transactionId, err := uuid.Parse(tempTransactionId)
	if err != nil {
		s.log.Error().Timestamp().Err(err).Msg("can't output money")
		return
	}

	// work...
	time.Sleep(5 * time.Second)

	// check balance
	balance, err := s.repository.GetBalance(ctx, userId)
	if err != nil {
		s.log.Error().Timestamp().Err(err).Msg("can't output money")
		return
	}
	if balance < count {
		s.log.Error().Timestamp().Err(errors.New("insufficient funds")).Msg("insufficient funds")
		return
	}

	// output money
	err = s.repository.OutputMoney(ctx, userId, count)
	if err != nil {
		s.log.Error().Timestamp().Err(errors.New("insufficient funds")).Msg("insufficient funds")
		return
	}

	// update transaction(confirmed) to true
	err = s.repository.UpdateTransactionsConfirmed(ctx, transactionId, true)
	if err != nil {
		s.log.Error().Timestamp().Err(errors.New("insufficient funds")).Msg("insufficient funds")
		return
	}

	// send data no rabbitmq
	//conn, err := net.Dial("tcp", url)
	//_, err = conn.Write([]byte(fmt.Sprintf("you got it %d", count)))
	//if err != nil {
	//	s.log.Error().Timestamp().Err(errors.New("insufficient funds")).Msg("insufficient funds")
	//	return
	//}
	_ = url
}
