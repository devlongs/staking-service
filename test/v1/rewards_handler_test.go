package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/devlongs/staking-service/internal/handler/v1"
	"github.com/go-chi/chi/v5"
)

func (d *DummyStakeService) Stake(walletAddress string, amount float64) error {
	return nil
}

func (d *DummyStakeService) GetRewards(walletAddress string) (float64, error) {
	return 5.0, nil
}

func TestRewardsHandler_Success(t *testing.T) {
	dummySvc := &DummyStakeService{}
	handler := v1.RewardsHandler(dummySvc)

	r := chi.NewRouter()
	r.Get("/v1/rewards/{wallet_address}", handler)
	req, err := http.NewRequest("GET", "/v1/rewards/0x1234567890abcdef1234567890abcdef12345678", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp map[string]interface{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}
	if resp["wallet_address"] != "0x1234567890abcdef1234567890abcdef12345678" {
		t.Errorf("unexpected wallet_address: got %v", resp["wallet_address"])
	}
	if resp["rewards"] != 5.0 {
		t.Errorf("unexpected rewards: got %v, want 5.0", resp["rewards"])
	}
}
