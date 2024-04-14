package service

import (
	"banner"
	"banner/pkg/repository"
)

type Authorization interface {
	CreateUser(user banner.User) (int, error)
	CheckNickNameAndEmail(nickname, email string) (int, error)
	GetPasswordHash(nickname string) (string, error)
	GenerateToken(nickname, passwordHash string) (string, error)
	ParseToken(accessToken string) (int, string, error)
}

type Banner interface {
	CheckBanner(tagIds []int, featureId int) (bool, error)
	CreateBanner(banner banner.Banner) (int, error)
	GetBannerById(id int) (banner.Banner, error)
	UpdateBannerById(id int, banner banner.Banner) error
	DeleteBannerById(id int) error
	GetUserBanner(input banner.UserBannerInput, role string) (banner.Content, error)
	GetAllBanners(input banner.FilterInput) ([]banner.Banner, error)
}

type Service struct {
	Authorization
	Banner
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Banner:        NewBannerService(repos.Banner),
	}
}
