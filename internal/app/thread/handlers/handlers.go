package handlers

import (
	customErr "forum/internal/app/errors"
	"forum/internal/app/httputils"
	"forum/internal/app/models"
	threadUsecase "forum/internal/app/thread/usecase"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
)

type Handlers struct {
	useCase threadUsecase.UseCase
}

func NewHandler(useCase threadUsecase.UseCase) *Handlers {
	return &Handlers{
		useCase: useCase,
	}
}

func (h *Handlers) ThreadInfo(ctx *fasthttp.RequestCtx) {
	idOrSlug := ctx.UserValue("slug_or_id").(string)
	thread, err := h.useCase.ThreadInfo(idOrSlug)
	if errors.Is(err, customErr.ErrForumNotFound) {
		resp := map[string]string{
			"message": "Can't find thread by slug or id: " + idOrSlug,
		}
		httputils.RespondErr(ctx, http.StatusNotFound, resp)
		return
	}
	if err != nil {
		httputils.Respond(ctx, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	httputils.Respond(ctx, http.StatusOK, thread)
}

func (h *Handlers) ChangeThread(ctx *fasthttp.RequestCtx) {
	var thread models.Thread
	if err := easyjson.Unmarshal(ctx.PostBody(), &thread); err != nil {
		httputils.Respond(ctx, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}

	idOrSlug := ctx.UserValue("slug_or_id").(string)
	thread, err := h.useCase.ChangeThread(idOrSlug, thread)

	if errors.Is(err, customErr.ErrThreadNotFound) {
		resp := map[string]string{
			"message": "Can't find thread by slug or id: " + idOrSlug,
		}
		httputils.RespondErr(ctx, http.StatusNotFound, resp)
		return
	}
	if err != nil {
		httputils.Respond(ctx, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	httputils.Respond(ctx, http.StatusOK, thread)
}

func (h *Handlers) VoteThread(ctx *fasthttp.RequestCtx) {
	var vote models.Vote
	if err := easyjson.Unmarshal(ctx.PostBody(), &vote); err != nil {
		httputils.Respond(ctx, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}

	idOrSlug := ctx.UserValue("slug_or_id").(string)
	nickname := vote.Nickname

	thread, err := h.useCase.VoteThread(idOrSlug, vote)

	if errors.Is(err, customErr.ErrThreadNotFound) {
		resp := map[string]string{
			"message": "Can't find thread by slug or id: " + idOrSlug,
		}
		httputils.RespondErr(ctx, http.StatusNotFound, resp)
		return
	}
	if errors.Is(err, customErr.ErrUserNotFound) {
		resp := map[string]string{
			"message": "Can't find user by nickname: " + nickname,
		}
		httputils.RespondErr(ctx, http.StatusNotFound, resp)
		return
	}
	if err != nil {
		httputils.Respond(ctx, http.StatusInternalServerError, nil)
		log.Println(err)
		return
	}
	httputils.Respond(ctx, http.StatusOK, thread)
}
