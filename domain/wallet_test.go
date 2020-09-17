package domain_test

import (
	"home-broker/domain"
	"testing"
)

func TestWalletNewMoneyFromFloatString(t *testing.T) {
	testTable := []struct {
		test     string
		expected domain.Money
	}{
		// MoneyDecimalPlaces = 6 = Float * 1.000.000
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
		v := domain.NewMoneyFromFloatString(table.test)
		if v != table.expected {
			t.Errorf("result[%v] is %v, expected %v", i, v, table.expected)
		}
	}
}
