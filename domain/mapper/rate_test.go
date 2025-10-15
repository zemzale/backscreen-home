package mapper

import (
	"cmp"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/zemzale/backscreen-home/domain/entity"
)

func TestRateFromXML(t *testing.T) {
	xmlFile, err := os.Open("testdata/ecb.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer xmlFile.Close()

	rates, err := RatesFromXML(xmlFile)
	if err != nil {
		t.Fatal(err)
	}

	expectedRates := []entity.Rate{
		{PublishedAt: time.Date(2025, time.October, 10, 3, 0, 0, 0, time.FixedZone("EEST", 3*60*60)), Code: "AUD", Value: "1.76500000"},
		{PublishedAt: time.Date(2025, time.October, 10, 3, 0, 0, 0, time.FixedZone("EEST", 3*60*60)), Code: "BGN", Value: "1.95580000"},
		{PublishedAt: time.Date(2025, time.October, 10, 3, 0, 0, 0, time.FixedZone("EEST", 3*60*60)), Code: "BRL", Value: "6.20820000"},
		{PublishedAt: time.Date(2025, time.October, 13, 3, 0, 0, 0, time.FixedZone("EEST", 3*60*60)), Code: "AUD", Value: "1.77750000"},
		{PublishedAt: time.Date(2025, time.October, 13, 3, 0, 0, 0, time.FixedZone("EEST", 3*60*60)), Code: "BGN", Value: "1.95580000"},
		{PublishedAt: time.Date(2025, time.October, 13, 3, 0, 0, 0, time.FixedZone("EEST", 3*60*60)), Code: "BRL", Value: "6.33440000"},
	}

	if len(rates) != len(expectedRates) {
		t.Fatalf("Expected 9 rates, got %d", len(rates))
	}

	sortFunc := func(a, b entity.Rate) int {
		codeDiff := cmp.Compare(a.Code, b.Code)
		if codeDiff != 0 {
			return codeDiff
		}
		return a.PublishedAt.Compare(b.PublishedAt)
	}

	slices.SortFunc(rates, sortFunc)
	slices.SortFunc(expectedRates, sortFunc)

	for i := range rates {
		want := expectedRates[i]
		got := rates[i]

		if got.Code != want.Code {
			t.Errorf("Expected rate %d to have code %s, got %s", i, want.Code, got.Code)
		}
		if got.Value != want.Value {
			t.Errorf("Expected rate %d to have value %s, got %s", i, want.Value, got.Value)
		}
		if got.PublishedAt.Compare(want.PublishedAt) != 0 {
			t.Errorf("Expected rate %d to have published at %s, got %s", i, want.PublishedAt, got.PublishedAt)
		}
	}
}
