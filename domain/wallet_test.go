package domain_test

import (
	"broker-dealer/domain"
	"testing"
)

const (
	failTestMSG = "The result is %v. Expected %v."
)

func TestCurrencyWithZeroValue(t *testing.T) {
	c := domain.NewCurrency(0, 0)
	expected := domain.Currency(0)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyWithOneUnit(t *testing.T) {
	c := domain.NewCurrency(1, 0)
	expected := domain.Currency(10000)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyOneUnitWithTwoCents(t *testing.T) {
	c := domain.NewCurrency(1, 200)
	expected := domain.Currency(10200)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyTenUnitsAnd99Cents(t *testing.T) {
	c := domain.NewCurrency(10, 9900)
	expected := domain.Currency(109900)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyWith4321Fractions(t *testing.T) {
	c := domain.NewCurrency(9999999, 4321)
	expected := domain.Currency(99999994321)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyFractionTruncation(t *testing.T) {
	c := domain.NewCurrency(9999999, 987654321)
	expected := domain.Currency(99999999876)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyOneNegativeValue(t *testing.T) {
	c := domain.NewCurrency(-1, 0)
	expected := domain.Currency(-10000)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyHugeNegativeValue(t *testing.T) {
	c := domain.NewCurrency(-9999999, 123456)
	expected := domain.Currency(-99999991234)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyFractionNegativeValue(t *testing.T) {
	c := domain.NewCurrency(9999999, -123456)
	expected := domain.Currency(99999991234)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyBothPartsAsNegativeValue(t *testing.T) {
	c := domain.NewCurrency(-9999999, -123456)
	expected := domain.Currency(-99999991234)
	if c != expected {
		t.Errorf(failTestMSG, c, expected)
	}
}

func TestCurrencyAsFloat64(t *testing.T) {
	c := domain.NewCurrency(123, 6789)
	expected := 123.6789
	result := c.AsFloat64()
	if result != expected {
		t.Errorf(failTestMSG, result, expected)
	}
}

func TestCurrencyAsString(t *testing.T) {
	c := domain.NewCurrency(123, 6789)
	expected := "123.6789"
	result := c.String()
	if result != expected {
		t.Errorf(failTestMSG, result, expected)
	}
}

/*
func TestAddWalletFunds(t *testing.T) {
	wallet := domain.NewWallet("1", 0.0)
	credit := 10.0
	wallet.Credit(credit)
	if wallet.Balance != credit {
		t.Errorf("Wallet balance is not %v, it is %v.", credit, wallet.Balance)
	}
}
*/
