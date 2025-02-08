package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "github.com/devlongs/staking-service/internal/handler/v1"
)

type DummyStakeService struct{}

func TestStakeHandler_Success(t *testing.T) {
	dummySvc := &DummyStakeService{}
	handler := v1.StakeHandler(dummySvc)

	payload := map[string]interface{}{
		"wallet_address": "0x1234567890abcdef1234567890abcdef12345678",
		"amount":         100,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/v1/stake", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var resp map[string]string
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}
	if msg, ok := resp["message"]; !ok || msg != "Staking successful" {
		t.Errorf("unexpected response message: %v", resp)
	}
}

func TestStakeHandler_InvalidWallet(t *testing.T) {
	dummySvc := &DummyStakeService{}
	handler := v1.StakeHandler(dummySvc)

	payload := map[string]interface{}{
		"wallet_address": "invalid",
		"amount":         100,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/v1/stake", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code for invalid wallet: got %v want %v", status, http.StatusBadRequest)
	}
}
