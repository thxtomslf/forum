package usecase

import (
	"forum/internal/app/models"
	threadRepo "forum/internal/app/thread/repository"
	"strconv"
)

type UseCase struct {
	threadRepo threadRepo.Repository
}

func NewUseCase(threadRepo threadRepo.Repository) *UseCase {
	return &UseCase{
		threadRepo: threadRepo,
	}
}

func (u *UseCase) ThreadInfo(idOrSlug string) (*models.Thread, error) {
	var id uint64
	var err error
	if id, err = strconv.ParseUint(idOrSlug, 10, 64); err != nil {
		thread, err := u.threadRepo.FindThreadBySlug(idOrSlug)
		if err != nil {
			return nil, err
		}
		return thread, nil
	}
	thread, err := u.threadRepo.FindThreadByID(id)
	if err != nil {
		return nil, err
	}
	return thread, nil
}

func (u *UseCase) ChangeThread(idOrSlug string, thread models.Thread) (models.Thread, error) {
	var id uint64
	var err error
	if id, err = strconv.ParseUint(idOrSlug, 10, 64); err != nil {
		thread, err = u.threadRepo.UpdateThreadBySlug(idOrSlug, thread)
		if err != nil {
			return models.Thread{}, err
		}
		return thread, nil
	}
	thread, err = u.threadRepo.UpdateThreadByID(id, thread)
	if err != nil {
		return models.Thread{}, err
	}
	return thread, nil
}

func (u *UseCase) VoteThread(idOrSlug string, vote models.Vote) (models.Thread, error) {
	thread, err := u.threadRepo.VoteThreadByID(idOrSlug, vote)
	if err != nil {
		return models.Thread{}, err
	}
	return thread, nil
}
