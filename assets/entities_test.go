package assets_test

import (
	"home-broker/assets"
	"home-broker/money"
	testassets "home-broker/tests/assets"
	"testing"
	"time"
)

func TestNewAssetUnitZero(t *testing.T) {
	m := assets.NewAssetUnitZero()
	if m != 0 {
		t.Errorf("received %v, expected 0", m)
	}
}
func TestNewAssetUnitFromFloatString(t *testing.T) {
	testTable := []struct {
		test     string
		expected assets.AssetUnit
	}{
		// AssetUnitDecimalPlaces = 6 = Float * 1.000.000
		{test: "0", expected: 0},
		{test: "0.0000001", expected: 0}, // truncates
		{test: "0.000001", expected: 1},
		{test: "0.100001", expected: 100001},
		{test: "1", expected: 1000000},
		{test: "1.999999", expected: 1999999},
		{test: "9", expected: 9000000},
		{test: "1.999999", expected: 1999999},
		{test: "1.999999999", expected: 1999999},
		{test: "999999999999.999999999", expected: 999999999999999999},
	}
	for i, table := range testTable {
		t.Run(table.test, func(t *testing.T) {
			v, err := assets.NewAssetUnitFromFloatString(table.test)
			if err != nil {
				t.Errorf("result[%v] returned an error: %v", i, err)
			}
			if v != table.expected {
				t.Errorf("result[%v] is %v, expected %v", i, v, table.expected)
			}
		})
	}
}

func TestNewMoneyFromFloatString_InvalidStrings(t *testing.T) {
	testTable := []string{"", "a", ".", ",", "9,9", "aa.bb", "99.a"}
	for i, vStr := range testTable {
		t.Run(vStr, func(t *testing.T) {
			_, err := money.NewMoneyFromFloatString(vStr)
			if err == nil {
				t.Errorf("result[%v] did not returned an error: %v", i, err)
			}
		})
	}
}

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
