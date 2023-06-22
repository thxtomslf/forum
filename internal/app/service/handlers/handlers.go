package handlers

import (
	"forum/internal/app/httputils"
	serviceUsecase "forum/internal/app/service/usecase"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

type Handlers struct {
	useCase serviceUsecase.UseCase
}

func NewHandler(useCase serviceUsecase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func (h *Handlers) ClearDB(ctx *fasthttp.RequestCtx) {
	err := h.useCase.ClearDB()
	if err != nil {
		httputils.Respond(ctx, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	httputils.Respond(ctx, http.StatusOK, nil)
}

func (h *Handlers) Status(ctx *fasthttp.RequestCtx) {
	numRec, err := h.useCase.Status()
	if err != nil {
		httputils.Respond(ctx, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	httputils.Respond(ctx, http.StatusOK, numRec)
}
