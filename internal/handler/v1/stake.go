package v1

import (
	"encoding/json"
	"net/http"

	"github.com/devlongs/staking-service/internal/service"
	"github.com/devlongs/staking-service/utils"
)

type stakeRequest struct {
	WalletAddress string  `json:"wallet_address"`
	Amount        float64 `json:"amount"`
}

// StakeHandler returns a handler for processing stake requests
func StakeHandler(stakeSvc service.StakeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req stakeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		// Validate Ethereum address
		if !utils.IsValidEthereumAddress(req.WalletAddress) {
			http.Error(w, "Invalid Ethereum address", http.StatusBadRequest)
			return
		}
		// Validate amount
		if req.Amount < 0 {
			http.Error(w, "Amount must be non-negative", http.StatusBadRequest)
			return
		}
		if err := stakeSvc.Stake(req.WalletAddress, req.Amount); err != nil {
			http.Error(w, "Failed to process staking", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Staking successful"})
	}
}
