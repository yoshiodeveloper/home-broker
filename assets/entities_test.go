package assets_test

import (
	"home-broker/assets"
	testassets "home-broker/tests/assets"
	"testing"
	"time"
)

func TestNewAsset(t *testing.T) {
	assetID := assets.AssetID("PETR4")
	exID := assets.ExchangeID("EXPETR4")
	asset := assets.NewAsset(assetID, exID)
	expected := assets.Asset{
		ID:         assetID,
		ExchangeID: exID,
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
		DeletedAt:  time.Time{},
	}
	err := testassets.CheckAssets(asset, expected)
	if err != nil {
		t.Error(err)
	}
}
