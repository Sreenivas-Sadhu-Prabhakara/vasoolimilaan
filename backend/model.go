package backend

import (
	"encoding/json"
	"fmt"
)

// DayClose separates the three things a day's till conflates: cash from today's
// sales, cash collected against OLD dues, and goods given on NEW credit.
type DayClose struct {
	Date          string  `json:"date"`
	CashSales     float64 `json:"cashSales"`
	DuesCollected float64 `json:"duesCollected"`
	NewCredit     float64 `json:"newCredit"`
}

// TrueSales is what you actually sold today (cash sales + goods on new credit).
func (d DayClose) TrueSales() float64 { return d.CashSales + d.NewCredit }

// CashIn is the money that hit the till (cash sales + dues recovered).
func (d DayClose) CashIn() float64 { return d.CashSales + d.DuesCollected }

// NetCreditChange is how the outstanding udhaar moved (new credit − dues collected).
func (d DayClose) NetCreditChange() float64 { return d.NewCredit - d.DuesCollected }

// Validate reports whether the DayClose is well formed.
func (d DayClose) Validate() error {
	if d.CashSales < 0 || d.DuesCollected < 0 || d.NewCredit < 0 {
		return fmt.Errorf("amounts cannot be negative")
	}
	return nil
}

// Summary aggregates the day-close ledger.
type Summary struct {
	Days             int     `json:"days"`
	TotalCashSales   float64 `json:"totalCashSales"`
	TotalDues        float64 `json:"totalDuesCollected"`
	TotalNewCredit   float64 `json:"totalNewCredit"`
	TotalTrueSales   float64 `json:"totalTrueSales"`
	NetCreditChange  float64 `json:"netCreditChange"`
}

// Summarize totals the components so recovery and fresh credit don't masquerade
// as good or bad sales days.
func Summarize(records []Record) Summary {
	var s Summary
	for _, r := range records {
		var d DayClose
		if json.Unmarshal(r.Input, &d) != nil {
			continue
		}
		s.Days++
		s.TotalCashSales += d.CashSales
		s.TotalDues += d.DuesCollected
		s.TotalNewCredit += d.NewCredit
	}
	s.TotalTrueSales = s.TotalCashSales + s.TotalNewCredit
	s.NetCreditChange = s.TotalNewCredit - s.TotalDues
	return s
}

// parseEntry decodes+validates a day-close; headline is true sales, label the date.
func parseEntry(raw []byte) (float64, string, error) {
	var d DayClose
	if err := json.Unmarshal(raw, &d); err != nil {
		return 0, "", fmt.Errorf("invalid json")
	}
	if err := d.Validate(); err != nil {
		return 0, "", err
	}
	return d.TrueSales(), d.Date, nil
}
