package domain

import (
	"time"
)

// Location 폐의약품 수거 위치 엔티티
type Location struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:200;not null" json:"name"`
	Type           string    `gorm:"size:50;not null;index" json:"type"`
	Address        string    `gorm:"size:500;not null" json:"address"`
	Latitude       float64   `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude      float64   `gorm:"type:decimal(11,8);not null" json:"longitude"`
	OperatingHours string    `gorm:"size:200" json:"operating_hours,omitempty"`
	Contact        string    `gorm:"size:50" json:"contact,omitempty"`
	Location       string    `gorm:"type:geography(POINT,4326);index:,type:gist" json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 테이블 이름 지정
func (Location) TableName() string {
	return "locations"
}

// LocationType 상수
const (
	LocationTypePostbox      = "postbox"
	LocationTypePharmacy     = "pharmacy"
	LocationTypeHealthCenter = "health_center"
)
