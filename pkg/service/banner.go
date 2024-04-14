package service

import (
	"banner"
	"banner/pkg/repository"
)

type BannerService struct {
	repo repository.Banner
}

func NewBannerService(repo repository.Banner) *BannerService {
	return &BannerService{repo: repo}
}

func (s *BannerService) CheckBanner(tagIds []int, featureId int) (bool, error) {
	return s.repo.CheckBanner(tagIds, featureId)
}

func (s *BannerService) CreateBanner(banner banner.Banner) (int, error) {
	return s.repo.CreateBanner(banner)
}

func (s *BannerService) GetBannerById(id int) (banner.Banner, error) {
	return s.repo.GetBannerById(id)
}

func (s *BannerService) UpdateBannerById(id int, banner banner.Banner) error {
	return s.repo.UpdateBannerById(id, banner)
}

func (s *BannerService) DeleteBannerById(id int) error {
	return s.repo.DeleteBannerById(id)
}

func (s *BannerService) GetUserBanner(input banner.UserBannerInput, role string) (banner.Content, error) {
	return s.repo.GetUserBanner(input, role)
}

func (s *BannerService) GetAllBanners(input banner.FilterInput) ([]banner.Banner, error) {
	return s.repo.GetAllBanners(input)
}
