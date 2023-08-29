package handlers

import (
	"bwg2/config"
	"bwg2/internal/tsfl"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"os"
	"time"
)

type service interface {
	ApplicationInputMoney(ctx context.Context, userId uuid.UUID, count uint64) (string, error)
	InputMoney(ctx context.Context, transactionId uuid.UUID) error
	OutputMoney(ctx context.Context, userId uuid.UUID, count uint64) error
	ApplicationOutputMoney(ctx context.Context, userId uuid.UUID, count uint64, url string) error
}

type FiberServer struct {
	app *fiber.App
	log *zerolog.Logger
	service
}

// NewFiberServer fabric build FiberServer
func NewFiberServer(conf *config.Config, serv service) *FiberServer {
	logg := zerolog.New(os.Stderr)

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Duration(conf.Server.Timeout * 1000000),
		WriteTimeout: time.Duration(conf.Server.Timeout * 1000000),
	})
	return &FiberServer{
		app:     app,
		log:     &logg,
		service: serv,
	}
}

// Run init handlers and start server on conf.Postgres.Port
func (fs *FiberServer) Run(conf *config.Config) error {
	fs.app.Put("input", fs.applicationInputMoney)
	fs.app.Put("output", fs.applicationOutputMoney)
	fs.app.Get(":transactionId", fs.inputMoney)

	err := fs.app.Listen(conf.Server.Port)
	if err != nil {
		fs.log.Error().Timestamp().Err(err).Send()
		return &tsfl.RespondError{
			Status:  500,
			Message: "",
			Err:     err,
		}
	}
	return err
}

// try input money on balance
func (fs *FiberServer) applicationInputMoney(ctx *fiber.Ctx) error {
	request := &tsfl.InputMoneyRequest{}

	err := ctx.BodyParser(&request)
	if err != nil {
		ctx.Status(400)
		respond := tsfl.RespondError{
			Status:  400,
			Message: "",
			Err:     err,
		}
		return ctx.SendString(respond.Error())
	}

	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		ctx.Status(400)
		respond := tsfl.RespondError{
			Status:  400,
			Message: "uncorrected uuid",
			Err:     err,
		}
		return ctx.SendString(respond.Error())
	}

	url, err := fs.service.ApplicationInputMoney(context.Background(), userId, request.Count)
	if err != nil {
		ctx.Status(err.(*tsfl.RespondError).Status)
		return ctx.SendString(err.Error())
	}

	return ctx.SendString(url)
}

// called when clicking on the link sent when calling applicationInputMoney
func (fs *FiberServer) inputMoney(ctx *fiber.Ctx) error {
	tempTransactionId := ctx.Params("transactionId", "")
	transactionId, err := uuid.Parse(tempTransactionId)
	if err != nil {
		ctx.Status(404)
		return ctx.Next()
	}
	err = fs.service.InputMoney(context.Background(), transactionId)
	if err != nil {
		ctx.Status(err.(*tsfl.RespondError).Status)
		return ctx.SendString(err.Error())
	}
	ctx.Status(200)
	return ctx.SendString("")
}

// try output money
func (fs *FiberServer) applicationOutputMoney(ctx *fiber.Ctx) error {
	request := &tsfl.OutputMoneyRequest{}

	err := ctx.BodyParser(&request)
	if err != nil {
		ctx.Status(400)
		respond := tsfl.RespondError{
			Status:  400,
			Message: "",
			Err:     err,
		}
		return ctx.SendString(respond.Error())
	}

	userId, err := uuid.Parse(request.UserId)
	if err != nil {
		ctx.Status(400)
		respond := tsfl.RespondError{
			Status:  400,
			Message: "uncorrected uuid",
			Err:     err,
		}
		return ctx.SendString(respond.Error())
	}
	err = fs.ApplicationOutputMoney(context.Background(), userId, request.Count, request.Url)
	if err != nil {
		ctx.Status(err.(*tsfl.RespondError).Status)
		return ctx.SendString(err.Error())
	}
	ctx.Status(200)
	return ctx.SendString("We have already started processing your request")
}
