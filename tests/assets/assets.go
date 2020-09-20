package tests

import (
	"fmt"
	"home-broker/assets"
	"time"
)

var (
	// BaseTime is a base time for assets.
	BaseTime = time.Date(2020, time.Month(1), 10, 11, 12, 13, 14, time.UTC)
)

// GetAsset returns an asset entity.
func GetAsset() assets.Asset {
	asset := assets.NewAsset("PETR4", "EXPETR4")
	asset.CreatedAt = BaseTime
	asset.UpdatedAt = BaseTime.Add(time.Hour * 2)
	return asset
}

// CheckAssets compares if two assets are equals.
func CheckAssets(a assets.Asset, b assets.Asset) error {
	if a.ID != b.ID {
		return fmt.Errorf("asset.ID is %v, expected %v", a.ID, b.ID)
	}
	if a.Name != b.Name {
		return fmt.Errorf("asset.Name is %v, expected %v", a.Name, b.Name)
	}
	if a.ExchangeID != b.ExchangeID {
		return fmt.Errorf("asset.ExchangeID is %v, expected %v", a.ExchangeID, b.ExchangeID)
	}
	if !a.CreatedAt.Equal(b.CreatedAt) {
		return fmt.Errorf("asset.CreatedAt is %v, expected %v", a.CreatedAt, b.CreatedAt)
	}
	if !a.UpdatedAt.Equal(b.UpdatedAt) {
		return fmt.Errorf("asset.UpdatedAt is %v, expected %v", a.UpdatedAt, b.UpdatedAt)
	}
	if !a.DeletedAt.Equal(b.DeletedAt) {
		return fmt.Errorf("asset.DeletedAt is %v, expected %v", a.DeletedAt, b.DeletedAt)
	}
	return nil
}
