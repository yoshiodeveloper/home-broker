package domain

import (
	"testing"
)

func TestCurrencyCreation(t *testing.T) {
	tests := []struct {
		currency    Currency
		expectedInt int64
	}{
		// $0.0000
		{NewCurrency(0, 0), 0},
		// $0.0009
		{NewCurrency(0, 9), 9},
		// $0.0099
		{NewCurrency(0, 99), 99},
		// $0.0999
		{NewCurrency(0, 999), 999},
		// $0.9999
		{NewCurrency(0, 9999), 9999},

		// $1.0
		{NewCurrency(1, 0), 10000},
		// $1.0200 ($1.02)
		{NewCurrency(1, 200), 10200},
		// $10.9900 ($10.99)
		{NewCurrency(10, 9900), 109900},
		// $9999999.4321
		{NewCurrency(9999999, 4321), 99999994321},
		// $9999999.987654321 to $9999999.9876
		{NewCurrency(9999999, 987654321), 99999999876},

		// -$0.0009
		{NewCurrency(0, -9), -9},
		// -$0.0099
		{NewCurrency(0, -99), -99},
		// -$0.0999
		{NewCurrency(0, -999), -999},
		// -$0.9999
		{NewCurrency(0, -9999), -9999},

		// Special cases using double negatives.
		{NewCurrency(-9, -9), -90009},
		{NewCurrency(-9, -9999), -99999},

		// -$1.0
		{NewCurrency(-1, 0), -10000},
		// -$1.0200 (-$1.02)
		{NewCurrency(-1, 200), -10200},
		// -$10.9900 (-$10.99)
		{NewCurrency(-10, 9900), -109900},
		// -$9,999,999.4321
		{NewCurrency(-9999999, 4321), -99999994321},
		// -$9,999,999.987654321 to -$9,999,999.9876
		{NewCurrency(-9999999, 987654321), -99999999876},
	}
	for _, test := range tests {
		if int64(test.currency) != test.expectedInt {
			t.Errorf(failTestMSG, int64(test.currency), test.expectedInt)
		}
	}
}

func TestCurrencyFromGoodStrings(t *testing.T) {
	tests := []struct {
		floatStr    string
		expectedInt int64
	}{
		{"0", 0},
		{"0.0", 0},
		{"0.00000", 0},
		{"0.0009", 9},
		{"0.0009999", 9},
		{"0.009", 90},
		{"0.0099", 99},
		{"0.0099999", 99},
		{"0.09", 900},
		{"0.099", 990},
		{"0.0999", 999},
		{"0.0999999", 999},
		{"0.9999", 9999},
		{"0.9999999", 9999},
		{"1.0", 10000},
		{"1.2", 12000},
		{"1.24", 12400},
		{"1.02", 10200},
		{"1.029", 10290},
		{"1.0299", 10299},
		{"10.9999999", 109999},
		{"9999999.9999999", 99999999999},
		{"9999999.987654321", 99999999876},

		{"-0", 0},
		{"-0.0", 0},
		{"-0.00000", 0},
		{"-0.0009", -9},
		{"-0.0009999", -9},
		{"-0.009", -90},
		{"-0.0099", -99},
		{"-0.0099999", -99},
		{"-0.09", -900},
		{"-0.099", -990},
		{"-0.0999", -999},
		{"-0.0999999", -999},
		{"-0.9999", -9999},
		{"-0.9999999", -9999},
		{"-1.0", -10000},
		{"-1.2", -12000},
		{"-1.24", -12400},
		{"-1.02", -10200},
		{"-1.029", -10290},
		{"-1.0299", -10299},
		{"-10.9999999", -109999},
		{"-9999999.9999999", -99999999999},
		{"-9999999.987654321", -99999999876},
	}
	for _, test := range tests {
		c, err := NewCurrencyFromString(test.floatStr)
		if err != nil {
			t.Errorf("the test of %v raised an error: %v", test.floatStr, err)
			continue
		}
		if int64(c) != test.expectedInt {
			t.Errorf("the result of \"%s\" is %v. Expected %v.", test.floatStr, int64(c), test.expectedInt)
		}
	}
}

func TestCurrencyFromBadStrings(t *testing.T) {
	tests := []string{
		"",
		"0.0.0",
		"a.0",
		"9.a",
		"9.-9",
		"a",
		"9,99",
	}
	for _, floatStr := range tests {
		_, err := NewCurrencyFromString(floatStr)
		if err == nil {
			t.Errorf("the result of \"%s\" didn't raise an error", floatStr)
			continue
		}
	}
}

func TestCurrencyAsFloat64(t *testing.T) {
	tests := []struct {
		currency      Currency
		expectedFloat float64
	}{
		// $0.0
		{NewCurrency(0, 0), 0.0},
		// $0.0009
		{NewCurrency(0, 9), 0.0009},
		// $0.9000
		{NewCurrency(0, 9000), 0.9000},
		// $0.0099
		{NewCurrency(0, 99), 0.0099},
		// $0.0999
		{NewCurrency(0, 999), 0.0999},
		// $0.9999
		{NewCurrency(0, 9999), 0.9999},
		// $9.9999
		{NewCurrency(9, 9999), 9.9999},
		// $9999999.999999 to $9999999.9999
		{NewCurrency(9999999, 999999), 9999999.9999},

		// -$0.0009
		{NewCurrency(0, -9), -0.0009},
		// -$0.9000
		{NewCurrency(0, -9000), -0.9000},
		// -$0.0099
		{NewCurrency(0, -99), -0.0099},
		// -$0.0999
		{NewCurrency(0, -999), -0.0999},
		// -$0.9999
		{NewCurrency(0, -9999), -0.9999},
		// -$9.9999
		{NewCurrency(-9, 9999), -9.9999},
		// -$9999999.999999 to -$9999999.9999
		{NewCurrency(-9999999, 999999), -9999999.9999},
	}
	for _, test := range tests {
		result := test.currency.AsFloat64()
		if result != test.expectedFloat {
			t.Errorf(failTestMSG, result, test.expectedFloat)
		}
	}
}

func TestCurrencyAsString(t *testing.T) {
	tests := []struct {
		currency       Currency
		expectedString string
	}{
		// $0.0
		{NewCurrency(0, 0), "0.0000"},
		// $0.0009
		{NewCurrency(0, 9), "0.0009"},
		// $0.9000
		{NewCurrency(0, 9000), "0.9000"},
		// $0.0099
		{NewCurrency(0, 99), "0.0099"},
		// $0.0999
		{NewCurrency(0, 999), "0.0999"},
		// $0.9999
		{NewCurrency(0, 9999), "0.9999"},
		// $9.9999
		{NewCurrency(9, 9999), "9.9999"},
		// $9999999.999999 to $9999999.9999
		{NewCurrency(9999999, 999999), "9999999.9999"},

		// -$0.0009
		{NewCurrency(0, -9), "-0.0009"},
		// -$0.9000
		{NewCurrency(0, -9000), "-0.9000"},
		// -$0.0099
		{NewCurrency(0, -99), "-0.0099"},
		// -$0.0999
		{NewCurrency(0, -999), "-0.0999"},
		// -$0.9999
		{NewCurrency(0, -9999), "-0.9999"},
		// -$9.9999
		{NewCurrency(-9, 9999), "-9.9999"},
		// -$9999999.999999 to -$9999999.9999
		{NewCurrency(-9999999, 999999), "-9999999.9999"},
	}
	for _, test := range tests {
		result := test.currency.String()
		if result != test.expectedString {
			t.Errorf(failTestMSG, result, test.expectedString)
		}
	}
}
