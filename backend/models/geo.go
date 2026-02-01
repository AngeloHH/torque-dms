package models

import (
	"encoding/json"
	"time"
)

type Country struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	ISOCode      string `gorm:"unique;size:2" json:"iso_code"`
	ISOCode3     string `gorm:"unique;size:3" json:"iso_code_3"`
	Name         string `json:"name"`
	PhoneCode    string `json:"phone_code"`
	CurrencyCode string `json:"currency_code"`
	FlagEmoji    string `json:"flag_emoji"`
	Active       bool   `gorm:"default:true" json:"active"`
}

type Location struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `json:"name"`
	Type      LocationType `json:"type"`
	Address   string       `json:"address"`
	City      string       `json:"city"`
	State     string       `json:"state"`
	Zip       string       `json:"zip"`
	CountryID uint         `json:"country_id"`
	Country   Country      `gorm:"foreignKey:CountryID" json:"country"`
	Latitude  float64      `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude float64      `gorm:"type:decimal(11,8)" json:"longitude"`
	Capacity  int          `json:"capacity"`
	Active    bool         `gorm:"default:true" json:"active"`
	CreatedAt time.Time    `json:"created_at"`
}

type Route struct {
	ID               uint            `gorm:"primaryKey" json:"id"`
	Name             string          `json:"name"`
	FromLocationID   uint            `json:"from_location_id"`
	FromLocation     Location        `gorm:"foreignKey:FromLocationID" json:"from_location"`
	ToLocationID     uint            `json:"to_location_id"`
	ToLocation       Location        `gorm:"foreignKey:ToLocationID" json:"to_location"`
	DistanceKM       float64         `json:"distance_km"`
	EstimatedMinutes int             `json:"estimated_minutes"`
	Waypoints        json.RawMessage `gorm:"type:json" json:"waypoints"`
	Active           bool            `gorm:"default:true" json:"active"`
	CreatedAt        time.Time       `json:"created_at"`
}