package backend

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type memStore struct{ items []Record }

func (m *memStore) Save(r Record) (Record, error) {
	r.ID = int64(len(m.items) + 1)
	m.items = append([]Record{r}, m.items...)
	return r, nil
}
func (m *memStore) List(limit int) ([]Record, error) { return m.items, nil }

func TestDerived(t *testing.T) {
	d := DayClose{CashSales: 3000, DuesCollected: 2000, NewCredit: 500}
	if math.Abs(d.TrueSales()-3500) > 1e-9 {
		t.Fatalf("trueSales=%v want 3500", d.TrueSales())
	}
	if math.Abs(d.CashIn()-5000) > 1e-9 {
		t.Fatalf("cashIn=%v want 5000", d.CashIn())
	}
	if math.Abs(d.NetCreditChange()+1500) > 1e-9 {
		t.Fatalf("netCreditChange=%v want -1500", d.NetCreditChange())
	}
}

func TestSummarize(t *testing.T) {
	in1, _ := json.Marshal(DayClose{CashSales: 3000, DuesCollected: 2000, NewCredit: 500})
	s := Summarize([]Record{{Input: in1}})
	if math.Abs(s.TotalTrueSales-3500) > 1e-9 || math.Abs(s.NetCreditChange+1500) > 1e-9 {
		t.Fatalf("summary wrong: %+v", s)
	}
}

func TestLogEndpoint(t *testing.T) {
	srv := NewServer(&memStore{})
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/log",
		strings.NewReader(`{"date":"2026-08-18","cashSales":3000,"duesCollected":2000,"newCredit":500}`)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("log %d", rec.Code)
	}
}
