package service

import (
	"github.com/devlongs/staking-service/internal/store"
)

// StakeService defines the methods for staking and rewards
type StakeService interface {
	Stake(walletAddress string, amount float64) error
	GetRewards(walletAddress string) (float64, error)
}

// stakeService implements the StakeService interface
type stakeService struct {
	storage store.Storage
}

// NewStakeService creates a new StakeService
func NewStakeService(storage store.Storage) StakeService {
	return &stakeService{storage: storage}
}

// Stake records a new staking transaction
func (s *stakeService) Stake(walletAddress string, amount float64) error {
	return s.storage.CreateStake(walletAddress, amount)
}

// GetRewards calculates rewards as 5% of the total staked amount.
func (s *stakeService) GetRewards(walletAddress string) (float64, error) {
	total, err := s.storage.GetTotalStaked(walletAddress)
	if err != nil {
		return 0, err
	}
	return total * 0.05, nil // 5% rewards
}
