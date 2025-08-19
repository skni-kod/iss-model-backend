package models

import (
	"time"
)

type ISSPosition struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"size:50;not null"`
	Latitude   float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude  float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	Altitude   float64   `json:"altitude" gorm:"type:decimal(10,5);not null"`
	Velocity   float64   `json:"velocity" gorm:"type:decimal(12,6);not null"`
	Visibility string    `json:"visibility" gorm:"size:20;not null"`
	Footprint  float64   `json:"footprint" gorm:"type:decimal(10,4);not null"`
	Timestamp  int64     `json:"timestamp" gorm:"uniqueIndex;not null"`
	Daynum     float64   `json:"daynum" gorm:"type:decimal(15,7);not null"`
	SolarLat   float64   `json:"solar_lat" gorm:"type:decimal(10,8);not null"`
	SolarLon   float64   `json:"solar_lon" gorm:"type:decimal(11,8);not null"`
	Units      string    `json:"units" gorm:"size:20;not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (ISSPosition) TableName() string {
	return "iss_positions"
}

type ISSPositionResponse struct {
	Name       string  `json:"name"`
	ID         int     `json:"id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Altitude   float64 `json:"altitude"`
	Velocity   float64 `json:"velocity"`
	Visibility string  `json:"visibility"`
	Footprint  float64 `json:"footprint"`
	Timestamp  int64   `json:"timestamp"`
	Daynum     float64 `json:"daynum"`
	SolarLat   float64 `json:"solar_lat"`
	SolarLon   float64 `json:"solar_lon"`
	Units      string  `json:"units"`
}

type HistoricalRequest struct {
	Timestamp int64  `json:"timestamp" validate:"required"`
	Units     string `json:"units,omitempty"`
}

func (r *ISSPositionResponse) ToISSPosition() *ISSPosition {
	return &ISSPosition{
		Name:       r.Name,
		Latitude:   r.Latitude,
		Longitude:  r.Longitude,
		Altitude:   r.Altitude,
		Velocity:   r.Velocity,
		Visibility: r.Visibility,
		Footprint:  r.Footprint,
		Timestamp:  r.Timestamp,
		Daynum:     r.Daynum,
		SolarLat:   r.SolarLat,
		SolarLon:   r.SolarLon,
		Units:      r.Units,
	}
}
