package repository

import (
	"github.com/neisii/yakburigo/internal/domain"
	"gorm.io/gorm"
)

type LocationRepository interface {
	FindNearby(lat, lon float64, radiusMeters int, locationType string) ([]LocationWithDistance, error)
	FindByID(id uint) (*domain.Location, error)
	Create(location *domain.Location) error
	Count() (int64, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

type LocationWithDistance struct {
	domain.Location
	Distance float64 `json:"distance"`
}

func (r *locationRepository) FindNearby(lat, lon float64, radiusMeters int, locationType string) ([]LocationWithDistance, error) {
	var results []LocationWithDistance

	query := `
		SELECT 
			id, name, type, address, latitude, longitude, 
			operating_hours, contact, created_at, updated_at,
			ST_Distance(
				location,
				ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography
			) as distance
		FROM locations
		WHERE ST_DWithin(
			location,
			ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography,
			$3
		)
	`

	args := []interface{}{lon, lat, radiusMeters}

	if locationType != "" && locationType != "all" {
		query += " AND type = $4"
		args = append(args, locationType)
	}

	query += " ORDER BY distance ASC"

	err := r.db.Raw(query, args...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *locationRepository) FindByID(id uint) (*domain.Location, error) {
	var location domain.Location
	err := r.db.First(&location, id).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) Create(location *domain.Location) error {
	return r.db.Create(location).Error
}

func (r *locationRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.Location{}).Count(&count).Error
	return count, err
}
