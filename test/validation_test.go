package contract_test

import (
	"fmt"
	"strings"
	"testing"
)

// ===================================
// Input Validation Tests (Negative)
// ===================================

// helpers

func longString(n int) string { return strings.Repeat("a", n) }

func pipeAddress() string { return "hive:user|evil" }

// ===================================
// Address Validation
// ===================================

func TestMintFailsAddressTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"to":%q,"id":"1","amount":1,"maxSupply":10,"data":""}`, addr))
	CallContract(t, ct, "mint", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestMintFailsAddressWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(fmt.Sprintf(`{"to":%q,"id":"1","amount":1,"maxSupply":10,"data":""}`, pipeAddress()))
	CallContract(t, ct, "mint", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTransferFromFailsFromAddressTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"1","amount":5,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"from":%q,"to":"hive:other","ids":["1"],"amounts":[1],"data":""}`, addr))
	CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTransferFromFailsToAddressTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"1","amount":5,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"from":"hive:tibfox","to":%q,"id":"1","amount":1,"data":""}`, addr))
	CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTransferFromFailsFromAddressWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"1","amount":5,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(fmt.Sprintf(`{"from":%q,"to":"hive:other","id":"1","amount":1,"data":""}`, pipeAddress()))
	CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTransferFromFailsToAddressWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"1","amount":5,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(fmt.Sprintf(`{"from":"hive:tibfox","to":%q,"id":"1","amount":1,"data":""}`, pipeAddress()))
	CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnFailsFromAddressTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"from":%q,"id":"1","amount":1}`, addr))
	CallContract(t, ct, "burn", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnFailsFromAddressWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(fmt.Sprintf(`{"from":%q,"id":"1","amount":1}`, pipeAddress()))
	CallContract(t, ct, "burn", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetApprovalFailsOperatorTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"operator":%q,"approved":true}`, addr))
	CallContract(t, ct, "setApprovalForAll", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetApprovalFailsOperatorWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(fmt.Sprintf(`{"operator":%q,"approved":true}`, pipeAddress()))
	CallContract(t, ct, "setApprovalForAll", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestChangeOwnerFailsNewOwnerTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	addr := longString(257)
	payload := []byte(fmt.Sprintf(`{"newOwner":%q}`, addr))
	CallContract(t, ct, "changeOwner", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestChangeOwnerFailsNewOwnerWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(fmt.Sprintf(`{"newOwner":%q}`, pipeAddress()))
	CallContract(t, ct, "changeOwner", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Token ID Validation
// ===================================

func TestMintFailsTokenIdTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	id := longString(257)
	payload := []byte(fmt.Sprintf(`{"to":"hive:tibfox","id":%q,"amount":1,"maxSupply":10,"data":""}`, id))
	CallContract(t, ct, "mint", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestMintFailsTokenIdWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(`{"to":"hive:tibfox","id":"token|evil","amount":1,"maxSupply":10,"data":""}`)
	CallContract(t, ct, "mint", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTransferFromFailsTokenIdTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	id := longString(257)
	payload := []byte(fmt.Sprintf(`{"from":"hive:tibfox","to":"hive:other","id":%q,"amount":1,"data":""}`, id))
	CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTransferFromFailsTokenIdWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(`{"from":"hive:tibfox","to":"hive:other","id":"tok|evil","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnFailsTokenIdTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	id := longString(257)
	payload := []byte(fmt.Sprintf(`{"from":"hive:tibfox","id":%q,"amount":1}`, id))
	CallContract(t, ct, "burn", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnFailsTokenIdWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(`{"from":"hive:tibfox","id":"tok|evil","amount":1}`)
	CallContract(t, ct, "burn", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetURIFailsTokenIdTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	id := longString(257)
	payload := []byte(fmt.Sprintf(`{"id":%q,"uri":"https://example.com/"}`, id))
	CallContract(t, ct, "setURI", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetURIFailsTokenIdWithPipe(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(`{"id":"tok|evil","uri":"https://example.com/"}`)
	CallContract(t, ct, "setURI", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// URI Validation
// ===================================

func TestSetURIFailsURITooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	uri := "https://example.com/" + longString(1010)
	payload := []byte(fmt.Sprintf(`{"id":"1","uri":%q}`, uri))
	CallContract(t, ct, "setURI", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetBaseURIFailsNoTrailingSlash(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	payload := []byte(`{"baseUri":"https://example.com/nft"}`)
	CallContract(t, ct, "setBaseURI", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetBaseURIFailsTooLong(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	uri := "https://example.com/" + longString(1010) + "/"
	payload := []byte(fmt.Sprintf(`{"baseUri":%q}`, uri))
	CallContract(t, ct, "setBaseURI", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Init: Name / Symbol Length
// ===================================

func TestInitFailsNameTooLong(t *testing.T) {
	ct := SetupContractTest()
	name := longString(65)
	payload := []byte(fmt.Sprintf(`{"name":%q,"symbol":"MNFT","baseUri":"https://api.magi.network/metadata/"}`, name))
	CallContract(t, ct, "init", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestInitFailsSymbolTooLong(t *testing.T) {
	ct := SetupContractTest()
	symbol := longString(17)
	payload := []byte(fmt.Sprintf(`{"name":"Magi NFT","symbol":%q,"baseUri":"https://api.magi.network/metadata/"}`, symbol))
	CallContract(t, ct, "init", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestInitFailsBaseURINoTrailingSlash(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, false, uint(150_000_000), "")
}
