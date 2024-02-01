package handler

import "fio-expander/internal/app/model"

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type IDResponse struct {
	Response
	ID string `json:"id"`
}

type FindResponse struct {
	Response
	Finded []model.ExtendedFIO
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

type FIOExp interface {
	Create(e *model.ExtendedFIO) error
	Update(e *model.ExtendedFIO) error
	Delete(e *model.ExtendedFIO) error
	FindBy(e *model.FindFIO) (error, *[]model.ExtendedFIO)
}
