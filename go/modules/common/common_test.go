package bddcommon_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"example.com/shop/bddcommon"
)

func TestMoneyParsesDollars(t *testing.T) {
	for text, want := range map[string]bddcommon.Money{"$12.50": 1250, "$12": 1200, "$0.00": 0} {
		got, err := bddcommon.ParseMoney(text)
		if err != nil || got != want {
			t.Errorf("ParseMoney(%q) = %v, %v; want %v", text, got, err, want)
		}
	}
	if _, err := bddcommon.ParseMoney("12.50"); err == nil {
		t.Error("money without $ should be refused")
	}
}

func TestToMoneyRoundsHalfUpToCents(t *testing.T) {
	for text, want := range map[string]bddcommon.Money{"12.5": 1250, "1.005": 101, "2": 200, "0.004": 0, "-1.005": -101} {
		got, err := bddcommon.ToMoney(text)
		if err != nil || got != want {
			t.Errorf("ToMoney(%q) = %v, %v; want %v", text, got, err, want)
		}
	}
}

func TestMoneyString(t *testing.T) {
	if got := bddcommon.Money(350).String(); got != "$3.50" {
		t.Errorf("got %s", got)
	}
	if got := bddcommon.Money(-330).String(); got != "-$3.30" {
		t.Errorf("got %s", got)
	}
}

func TestOutcomeKeepsTheLastError(t *testing.T) {
	var o bddcommon.Outcome
	o.Attempt(func() error { return errors.New("no") })
	if o.LastError == nil || o.LastError.Error() != "no" {
		t.Fatalf("got %v", o.LastError)
	}
	o.Attempt(func() error { return nil })
	if o.LastError != nil {
		t.Fatalf("a successful attempt should clear the error, got %v", o.LastError)
	}
}

func TestDeprecatedStepWarns(t *testing.T) {
	var out bytes.Buffer
	saved := bddcommon.DeprecationOutput
	bddcommon.DeprecationOutput = &out
	defer func() { bddcommon.DeprecationOutput = saved }()

	bddcommon.Deprecated("old words", "new words", "1.1.0")

	if !strings.Contains(out.String(), `Use "new words" instead`) {
		t.Fatalf("got %q", out.String())
	}
}
