package usecase

import "github.com/RakibulBh/AI-pr-reviewer/internal/repository"

type GithubUsecase struct {
	repository *repository.GithubRepository
}

func NewGithubUsecase(repository *repository.GithubRepository) *GithubUsecase {
	return &GithubUsecase{
		repository: repository,
	}
}
