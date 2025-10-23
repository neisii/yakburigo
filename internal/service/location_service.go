package service

import (
	"errors"
	"github.com/neisii/yakburigo/internal/repository"
)

type LocationService interface {
	GetNearbyLocations(lat, lon float64, radiusMeters int, locationType string) ([]repository.LocationWithDistance, error)
	GetLocationByID(id uint) (*repository.LocationWithDistance, error)
}

type locationService struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationService{repo: repo}
}

func (s *locationService) GetNearbyLocations(lat, lon float64, radiusMeters int, locationType string) ([]repository.LocationWithDistance, error) {
	if lat < -90 || lat > 90 {
		return nil, errors.New("latitude must be between -90 and 90")
	}
	if lon < -180 || lon > 180 {
		return nil, errors.New("longitude must be between -180 and 180")
	}
	if radiusMeters <= 0 {
		return nil, errors.New("radius must be positive")
	}
	if radiusMeters > 50000 {
		return nil, errors.New("radius must be less than 50km")
	}

	locations, err := s.repo.FindNearby(lat, lon, radiusMeters, locationType)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func (s *locationService) GetLocationByID(id uint) (*repository.LocationWithDistance, error) {
	if id == 0 {
		return nil, errors.New("invalid location id")
	}

	location, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := &repository.LocationWithDistance{
		Location: *location,
		Distance: 0,
	}

	return result, nil
}
