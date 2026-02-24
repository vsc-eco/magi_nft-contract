package main

import (
	"magi_nft/sdk"
	"strconv"

	"github.com/CosmWasm/tinyjson/jwriter"
)

// ===================================
// Internal Helper Functions (ERC-1155)
// ===================================

// ===================================
// Safe Math Utilities
// ===================================

// safeAdd performs a + b addition. Aborts execution if an overflow is detected.
func safeAdd(a, b uint64) uint64 {
	sum := a + b
	if sum < a {
		sdk.Abort("safeAdd overflow")
	}
	return sum
}

// safeSub performs a - b subtraction. Aborts execution if an underflow is detected.
func safeSub(a, b uint64) uint64 {
	if b > a {
		sdk.Abort("safeSub underflow")
	}
	return a - b
}

// ===================================
// Balance Management (ERC-1155)
// ===================================

// balanceKey returns the state key for an account's balance of a specific token ID.
func balanceKey(account, tokenId string) string {
	return "bal|" + account + "|" + tokenId
}

// getBalance retrieves token balance of an address for a specific token ID.
func getBalance(account, tokenId string) uint64 {
	bal := sdk.StateGetObject(balanceKey(account, tokenId))
	if bal == nil {
		return 0
	}
	amt, _ := strconv.ParseUint(*bal, 10, 64)
	return amt
}

// setBalance sets the balance of an address for a specific token ID.
func setBalance(account, tokenId string, amount uint64) {
	sdk.StateSetObject(balanceKey(account, tokenId), strconv.FormatUint(amount, 10))
}

// incBalance increments token balance of an address for a specific token ID.
func incBalance(account, tokenId string, amount uint64) {
	oldBal := getBalance(account, tokenId)
	newBal := safeAdd(oldBal, amount)
	setBalance(account, tokenId, newBal)
}

// decBalance decrements token balance of an address for a specific token ID.
// Aborts if insufficient balance.
func decBalance(account, tokenId string, amount uint64) {
	oldBal := getBalance(account, tokenId)
	if oldBal < amount {
		sdk.Abort("Insufficient balance")
	}
	newBal := safeSub(oldBal, amount)
	setBalance(account, tokenId, newBal)
}

// ===================================
// Operator Approval Management (ERC-1155)
// ===================================

// operatorKey returns the state key for operator approval.
func operatorKey(owner, operator string) string {
	return "op|" + owner + "|" + operator
}

// isApprovedForAllInternal checks if an operator is approved for all tokens of an owner.
func isApprovedForAllInternal(owner, operator string) bool {
	approved := sdk.StateGetObject(operatorKey(owner, operator))
	if approved == nil {
		return false
	}
	return *approved == "1"
}

// setApprovalForAllInternal sets operator approval for all tokens of an owner.
func setApprovalForAllInternal(owner, operator string, approved bool) {
	if approved {
		sdk.StateSetObject(operatorKey(owner, operator), "1")
	} else {
		sdk.StateSetObject(operatorKey(owner, operator), "0")
	}
}

// ===================================
// Token URI Management
// ===================================

// uriKey returns the state key for a token's URI.
func uriKey(tokenId string) string {
	return "uri|" + tokenId
}

// getTokenURI retrieves the URI for a specific token ID.
// Falls back to baseURI + tokenId if no specific URI is set.
func getTokenURI(tokenId string) string {
	// Check for token-specific URI first
	uri := sdk.StateGetObject(uriKey(tokenId))
	if uri != nil && *uri != "" {
		return *uri
	}
	// Fall back to baseURI + tokenId
	baseURI := getBaseURI()
	if baseURI != "" {
		return baseURI + tokenId
	}
	return ""
}

// setTokenURI sets a specific URI for a token ID.
func setTokenURI(tokenId, uri string) {
	sdk.StateSetObject(uriKey(tokenId), uri)
}

// ===================================
// Supply Management
// ===================================

// maxSupplyKey returns the state key for a token's max supply.
func maxSupplyKey(tokenId string) string {
	return "max|" + tokenId
}

// totalSupplyKey returns the state key for a token's total minted supply.
func totalSupplyKey(tokenId string) string {
	return "tot|" + tokenId
}

// getMaxSupply retrieves the max supply for a token ID (0 = not set).
func getMaxSupply(tokenId string) uint64 {
	val := sdk.StateGetObject(maxSupplyKey(tokenId))
	if val == nil {
		return 0
	}
	amt, _ := strconv.ParseUint(*val, 10, 64)
	return amt
}

// setMaxSupply sets the max supply for a token ID (can only be set once).
func setMaxSupply(tokenId string, maxSupply uint64) {
	sdk.StateSetObject(maxSupplyKey(tokenId), strconv.FormatUint(maxSupply, 10))
}

// getTotalSupply retrieves the total minted supply for a token ID.
func getTotalSupply(tokenId string) uint64 {
	val := sdk.StateGetObject(totalSupplyKey(tokenId))
	if val == nil {
		return 0
	}
	amt, _ := strconv.ParseUint(*val, 10, 64)
	return amt
}

// incTotalSupply increments the total supply for a token ID.
func incTotalSupply(tokenId string, amount uint64) {
	oldSupply := getTotalSupply(tokenId)
	newSupply := safeAdd(oldSupply, amount)
	sdk.StateSetObject(totalSupplyKey(tokenId), strconv.FormatUint(newSupply, 10))
}

// decTotalSupply decrements the total supply for a token ID (on burn).
func decTotalSupply(tokenId string, amount uint64) {
	oldSupply := getTotalSupply(tokenId)
	newSupply := safeSub(oldSupply, amount)
	sdk.StateSetObject(totalSupplyKey(tokenId), strconv.FormatUint(newSupply, 10))
}

// totalMintedKey returns the state key for a token's total minted count (never decreases).
func totalMintedKey(tokenId string) string {
	return "minted|" + tokenId
}

// getTotalMinted retrieves the total ever minted for a token ID.
func getTotalMinted(tokenId string) uint64 {
	val := sdk.StateGetObject(totalMintedKey(tokenId))
	if val == nil {
		return 0
	}
	amt, _ := strconv.ParseUint(*val, 10, 64)
	return amt
}

// incTotalMinted increments the total minted count for a token ID (only increases).
func incTotalMinted(tokenId string, amount uint64) {
	oldMinted := getTotalMinted(tokenId)
	newMinted := safeAdd(oldMinted, amount)
	sdk.StateSetObject(totalMintedKey(tokenId), strconv.FormatUint(newMinted, 10))
}

// isTrackMintedEnabled checks if totalMinted tracking is enabled for this contract.
func isTrackMintedEnabled() bool {
	val := sdk.StateGetObject("track_minted")
	if val == nil {
		return false
	}
	return *val == "1"
}

// ===================================
// Soulbound Token Management
// ===================================

// soulboundKey returns the state key for a token's soulbound status.
func soulboundKey(tokenId string) string {
	return "sb|" + tokenId
}

// isSoulbound checks if a token is soulbound (non-transferable).
func isSoulbound(tokenId string) bool {
	val := sdk.StateGetObject(soulboundKey(tokenId))
	return val != nil && *val == "1"
}

// setSoulbound marks a token as soulbound (non-transferable).
func setSoulbound(tokenId string) {
	sdk.StateSetObject(soulboundKey(tokenId), "1")
}

// ===================================
// Token Properties Management
// ===================================

// propertiesKey returns the state key for a token's properties.
func propertiesKey(tokenId string) string {
	return "props|" + tokenId
}

// getTokenProperties retrieves the raw JSON properties for a token ID.
func getTokenProperties(tokenId string) string {
	val := sdk.StateGetObject(propertiesKey(tokenId))
	if val == nil {
		return ""
	}
	return *val
}

// setTokenProperties stores raw JSON properties for a token ID.
func setTokenProperties(tokenId, properties string) {
	sdk.StateSetObject(propertiesKey(tokenId), properties)
}

// ===================================
// Contract Properties (from state)
// ===================================

// getContractName retrieves the contract name from state.
func getContractName() string {
	n := sdk.StateGetObject("contract_name")
	if n == nil {
		return ""
	}
	return *n
}

// getContractSymbol retrieves the contract symbol from state.
func getContractSymbol() string {
	s := sdk.StateGetObject("contract_symbol")
	if s == nil {
		return ""
	}
	return *s
}

// getBaseURI retrieves the base URI from state.
func getBaseURI() string {
	u := sdk.StateGetObject("base_uri")
	if u == nil {
		return ""
	}
	return *u
}

// ===================================
// Transfer Authorization
// ===================================

// isApprovedOrOwner checks if the caller can transfer tokens from an account.
func isApprovedOrOwner(caller, from string) bool {
	if caller == from {
		return true
	}
	return isApprovedForAllInternal(from, caller)
}

// ===================================
// JSON Response Helper
// ===================================

func jsonResponse(marshaler interface{ MarshalTinyJSON(*jwriter.Writer) }) *string {
	w := jwriter.Writer{}
	marshaler.MarshalTinyJSON(&w)
	result := string(w.Buffer.BuildBytes())
	return &result
}
