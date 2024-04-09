package banner

import (
	"context"
	"errors"
	"slices"

	"avito-backend-trainee-2024/internal/domain/entity"

	entityutils "avito-backend-trainee-2024/internal/pkg/utils/entity"
	sliceutils "avito-backend-trainee-2024/pkg/utils/slice"
)

type BannerRepo interface {
	GetAllBanners(ctx context.Context, offset, limit int) ([]*entity.Banner, error)
	GetBannerByID(ctx context.Context, id int) (*entity.Banner, error)
	GetBannerByFeatureAndTags(ctx context.Context, featureID int, tagIDs []int) (*entity.Banner, error)
	CreateBanner(ctx context.Context, banner entity.Banner) (*entity.Banner, error)
	UpdateBanner(ctx context.Context, id int, updateModel entity.Banner) error
	DeleteBanner(ctx context.Context, id int) (*entity.Banner, error)
}

type FeatureRepo interface {
	GetFeatureByID(ctx context.Context, id int) (*entity.Feature, error)
}

type TagRepo interface {
	GetTagsWithIDs(ctx context.Context, IDs []int) ([]*entity.Tag, error)
}

type Service struct {
	BannerRepo  BannerRepo
	FeatureRepo FeatureRepo
	TagRepo     TagRepo
}

func New(bannerRepo BannerRepo, featureRepo FeatureRepo, tagRepo TagRepo) *Service {
	return &Service{
		BannerRepo:  bannerRepo,
		FeatureRepo: featureRepo,
		TagRepo:     tagRepo,
	}
}

// GetAllBanners todo: here decide data source (postgres or cached in redis)??
func (s *Service) GetAllBanners(ctx context.Context, offset, limit int) ([]*entity.Banner, error) {
	return s.BannerRepo.GetAllBanners(ctx, offset, limit)
}

func (s *Service) GetBannerByFeatureAndTags(ctx context.Context, featureID int, tagIDs []int) (*entity.Banner, error) {
	slices.Sort(tagIDs) // sort slice
	return s.BannerRepo.GetBannerByFeatureAndTags(ctx, featureID, tagIDs)
}

// validateBanner checks if associated with banner tags and feature are presented in db
func (s *Service) validateBanner(ctx context.Context, banner entity.Banner) error {
	feature, err := s.FeatureRepo.GetFeatureByID(ctx, banner.FeatureID)
	if err != nil || feature == nil {
		return ErrNoSuchFeature
	}

	tags, err := s.TagRepo.GetTagsWithIDs(ctx, banner.TagIDs)
	if err != nil {
		return errors.Join(ErrNoSuchTag, err)
	}

	// compare two slices: sorted(tagIDs) and tags by tag.ID field
	slices.Sort(banner.TagIDs)

	if len(banner.TagIDs) != len(tags) {
		return ErrNoSuchTag
	}

	if !sliceutils.Equals(banner.TagIDs, sliceutils.Map(tags, func(tag *entity.Tag) int { return tag.ID })) {
		return ErrNoSuchTag
	}

	return nil
}

func (s *Service) CreateBanner(ctx context.Context, banner entity.Banner) (*entity.Banner, error) {
	// firstly validate that feature and tags associated with banner exists in db
	if err := s.validateBanner(ctx, banner); err != nil {
		return nil, err
	}

	return s.BannerRepo.CreateBanner(ctx, banner)
}

func (s *Service) UpdateBanner(ctx context.Context, id int, updateModel entity.Banner) error {
	// firstly validate that feature and tags associated with banner exists in db
	if err := s.validateBanner(ctx, updateModel); err != nil {
		return err
	}

	// TODO: update only init fields
	banner, err := s.BannerRepo.GetBannerByID(ctx, id)
	if err != nil {
		return ErrNoSuchBanner
	}

	entityutils.InitNilFieldsOfBanner(&updateModel, banner) // TODO: MAYBE CHANGE REPO METHOD FOR OPTIMIZATION?

	return s.BannerRepo.UpdateBanner(ctx, id, updateModel)
}

func (s *Service) DeleteBanner(ctx context.Context, id int) (*entity.Banner, error) {
	return s.BannerRepo.DeleteBanner(ctx, id)
}