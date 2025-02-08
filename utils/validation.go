package utils

import (
	"regexp"
	"strings"
)

// IsValidEthereumAddress validates that an Ethereum address starts with "0x"
// and is followed by exactly 40 hexadecimal characters
func IsValidEthereumAddress(address string) bool {
	if len(address) != 42 {
		return false
	}
	if !strings.HasPrefix(address, "0x") {
		return false
	}
	matched, err := regexp.MatchString("^0x[0-9a-fA-F]{40}$", address)
	if err != nil {
		return false
	}
	return matched
}
