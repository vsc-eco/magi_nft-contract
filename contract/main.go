package main

import "magi_nft/sdk"

// ========================================
// MAGI Token Contract (ERC-20 Based)
// ========================================

func main() {
	// placeholder function
}

// =========================
// Initialization Functions
// =========================

// isInit returns true if the token has been initialized.
func isInit() bool {
	i := sdk.StateGetObject("isInit")
	if i == nil {
		return false
	}
	return *i != ""
}

// assertInit aborts execution if token has not been initialized.
func assertInit() {
	if !isInit() {
		sdk.Abort("Token not initialized")
	}
}

// getOwner returns the contract owner address and whether caller is the owner.
func getOwner() (string, bool) {
	i := sdk.StateGetObject("owner")
	if i == nil {
		return "", false
	}
	if *i == "" {
		return "", false
	}
	caller := sdk.GetEnvKey("msg.sender")
	if caller == nil {
		return *i, false
	}
	return *i, *i == *caller
}

// =========================
// Pausable Functions
// =========================

// isPaused returns true if the contract is paused.
func isPaused() bool {
	p := sdk.StateGetObject("paused")
	if p == nil {
		return false
	}
	return *p == "1"
}

// assertNotPaused aborts execution if contract is paused.
func assertNotPaused() {
	if isPaused() {
		sdk.Abort("Contract is paused")
	}
}
