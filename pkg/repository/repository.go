package repository

import (
	"banner"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user banner.User) (int, error)
	CheckNickNameAndEmail(nickname, email string) (int, error)
	GetPasswordHash(nickname string) (string, error)
	GetUser(nickname, password string) (banner.User, error)
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

type Repository struct {
	Authorization
	Banner
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Banner:        NewBannerPostgres(db),
	}
}
