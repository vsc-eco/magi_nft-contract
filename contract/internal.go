package main

import (
	"encoding/binary"
	"magi_nft/sdk"

	"github.com/CosmWasm/tinyjson/jwriter"
)

// ===================================
// Internal Helper Functions (ERC-1155)
// ===================================

// ===================================
// Input Validation
// ===================================

// Maximum allowed lengths for user-controlled input fields.
const (
	maxAddressLen = 256
	maxTokenIdLen = 256
	maxURILen     = 1024
	maxNameLen    = 64
	maxSymbolLen  = 16
)

// validateAddress checks that an address is within length bounds and contains
// no pipe characters, which are used as state key delimiters.
func validateAddress(account string) {
	if len(account) > maxAddressLen {
		sdk.Abort("Address exceeds maximum length")
	}
	for i := 0; i < len(account); i++ {
		if account[i] == '|' {
			sdk.Abort("Invalid character in address")
		}
	}
}

// validateTokenId checks that a token ID is within length bounds and contains
// no pipe characters, which are used as state key delimiters.
func validateTokenId(id string) {
	if len(id) > maxTokenIdLen {
		sdk.Abort("Token ID exceeds maximum length")
	}
	for i := 0; i < len(id); i++ {
		if id[i] == '|' {
			sdk.Abort("Invalid character in token ID")
		}
	}
}

// validateURI checks that a URI is within length bounds.
func validateURI(uri string) {
	if len(uri) > maxURILen {
		sdk.Abort("URI exceeds maximum length")
	}
}

// validateBaseURI checks that a base URI is within length bounds and ends with
// a trailing slash when non-empty, ensuring safe concatenation with token IDs.
func validateBaseURI(uri string) {
	if uri == "" {
		return
	}
	validateURI(uri)
	if uri[len(uri)-1] != '/' {
		sdk.Abort("Base URI must end with a trailing slash")
	}
}

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
	if bal == nil || *bal == "" {
		return 0
	}
	amt := bytesToU64([]byte(*bal))
	return amt
}

// setBalance sets the balance of an address for a specific token ID.
// Deletes the key when balance reaches 0 to avoid state bloat.
func setBalance(account, tokenId string, amount uint64) {
	key := balanceKey(account, tokenId)
	if amount == 0 {
		sdk.StateDeleteObject(key)
		return
	}
	sdk.StateSetObject(key, string(u64ToBytes(amount)))
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
// Deletes the key when revoking to avoid state bloat.
func setApprovalForAllInternal(owner, operator string, approved bool) {
	key := operatorKey(owner, operator)
	if approved {
		sdk.StateSetObject(key, "1")
	} else {
		sdk.StateDeleteObject(key)
	}
}

// ===================================
// Token Allowance Management (ERC-6909)
// ===================================

// allowanceKey returns the state key for per-token allowance.
func allowanceKey(owner, spender, tokenId string) string {
	return "allow|" + owner + "|" + spender + "|" + tokenId
}

// getAllowance returns the allowance for a spender on a specific token.
func getAllowance(owner, spender, tokenId string) uint64 {
	val := sdk.StateGetObject(allowanceKey(owner, spender, tokenId))
	if val == nil || *val == "" {
		return 0
	}
	return bytesToU64([]byte(*val))
}

// setAllowance sets the allowance for a spender on a specific token.
// Deletes the key when allowance reaches 0 to avoid state bloat.
func setAllowance(owner, spender, tokenId string, amount uint64) {
	key := allowanceKey(owner, spender, tokenId)
	if amount == 0 {
		sdk.StateDeleteObject(key)
		return
	}
	sdk.StateSetObject(key, string(u64ToBytes(amount)))
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
	if val == nil || *val == "" {
		return 0
	}
	amt := bytesToU64([]byte(*val))
	return amt
}

// setMaxSupply sets the max supply for a token ID (can only be set once).
func setMaxSupply(tokenId string, maxSupply uint64) {
	sdk.StateSetObject(maxSupplyKey(tokenId), string(u64ToBytes(maxSupply)))
}

// getTotalSupply retrieves the total minted supply for a token ID.
func getTotalSupply(tokenId string) uint64 {
	val := sdk.StateGetObject(totalSupplyKey(tokenId))
	if val == nil || *val == "" {
		return 0
	}
	amt := bytesToU64([]byte(*val))
	return amt
}

// incTotalSupply increments the total supply for a token ID.
func incTotalSupply(tokenId string, amount uint64) {
	oldSupply := getTotalSupply(tokenId)
	newSupply := safeAdd(oldSupply, amount)
	sdk.StateSetObject(totalSupplyKey(tokenId), string(u64ToBytes(newSupply)))
}

// decTotalSupply decrements the total supply for a token ID (on burn).
func decTotalSupply(tokenId string, amount uint64) {
	oldSupply := getTotalSupply(tokenId)
	newSupply := safeSub(oldSupply, amount)
	sdk.StateSetObject(totalSupplyKey(tokenId), string(u64ToBytes(newSupply)))
}

// totalMintedKey returns the state key for a token's total minted count (never decreases).
func totalMintedKey(tokenId string) string {
	return "minted|" + tokenId
}

// getTotalMinted retrieves the total ever minted for a token ID.
func getTotalMinted(tokenId string) uint64 {
	val := sdk.StateGetObject(totalMintedKey(tokenId))
	if val == nil || *val == "" {
		return 0
	}
	amt := bytesToU64([]byte(*val))
	return amt
}

// incTotalMinted increments the total minted count for a token ID (only increases).
func incTotalMinted(tokenId string, amount uint64) {
	oldMinted := getTotalMinted(tokenId)
	newMinted := safeAdd(oldMinted, amount)
	sdk.StateSetObject(totalMintedKey(tokenId), string(u64ToBytes(newMinted)))
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
// Collection Metadata Management
// ===================================

// getCollectionMetadata retrieves the raw JSON collection metadata from state.
func getCollectionMetadata() string {
	val := sdk.StateGetObject("collection_metadata")
	if val == nil {
		return ""
	}
	return *val
}

// setCollectionMetadata stores raw JSON collection metadata in state.
func setCollectionMetadata(metadata string) {
	sdk.StateSetObject("collection_metadata", metadata)
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
// uint64 <-> []byte Helper
// ===================================

func u64ToBytes(val uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)

	// In Little Endian, leading zeros (most significant bytes) are at the end of the slice.
	// We iterate backwards to find the last non-zero byte.
	lastNonZeroIndex := len(b) - 1
	for lastNonZeroIndex >= 0 {
		if b[lastNonZeroIndex] != 0 {
			break
		}
		lastNonZeroIndex--
	}

	// If the value was 0, ensure we return at least one byte (0x00) instead of an empty slice.
	if lastNonZeroIndex < 0 {
		return []byte{0}
	}

	return b[:lastNonZeroIndex+1]
}

func bytesToU64(b []byte) uint64 {
	if len(b) > 8 {
		sdk.Abort("byte length less than or equal to 8")
	}

	// Create an 8-byte buffer initialized to zeros.
	buf := make([]byte, 8)

	// In Little Endian, the existing bytes are the least significant and go at the start.
	// Copy the input slice into the beginning of the buffer.
	copy(buf, b)

	return binary.LittleEndian.Uint64(buf)
}

// ===================================
// uint64 -> string Helper
// ===================================

// uint64ToStr converts a uint64 to its decimal string representation.
func uint64ToStr(n uint64) string {
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 20) // max 20 decimal digits for uint64
	pos := 20
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

// ===================================
// JSON Response Helper
// ===================================

func jsonResponse(marshaler interface{ MarshalTinyJSON(*jwriter.Writer) }) *string {
	w := jwriter.Writer{}
	marshaler.MarshalTinyJSON(&w)
	if w.Error != nil {
		sdk.Abort("JSON marshal error")
	}
	result := string(w.Buffer.BuildBytes())
	return &result
}
