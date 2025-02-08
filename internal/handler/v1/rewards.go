package v1

import (
	"encoding/json"
	"net/http"

	"github.com/devlongs/internal/service"
	"github.com/go-chi/chi/v5"
)

// RewardsHandler returns a handler that calculates and returns rewards
func RewardsHandler(stakeSvc service.StakeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		walletAddress := chi.URLParam(r, "wallet_address")

		if walletAddress == "" {
			http.Error(w, "Wallet address is required", http.StatusBadRequest)
			return
		}
		rewards, err := stakeSvc.GetRewards(walletAddress)
		if err != nil {
			http.Error(w, "Failed to retrieve rewards", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"wallet_address": walletAddress,
			"rewards":        rewards,
		})
	}
}
