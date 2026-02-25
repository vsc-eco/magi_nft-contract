package contract_test

import (
	"testing"
)

// ===================================
// SafeAdd Overflow Tests
// ===================================

func TestMintFailsOnSafeAddOverflow(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a token with a very large supply
	payload := []byte(`{"to":"hive:tibfox","id":"1","amount":18446744073709551615,"maxSupply":18446744073709551615,"data":""}`)
	CallContract(t, ct, "mint", payload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to mint 1 more — safeAdd(max_uint64, 1) should overflow
	payload2 := []byte(`{"to":"hive:tibfox","id":"1","amount":1,"data":""}`)
	CallContract(t, ct, "mint", payload2, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestMintBatchFailsOnSafeAddOverflow(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a token with a very large supply
	payload := []byte(`{"to":"hive:tibfox","ids":["1"],"amounts":[18446744073709551615],"maxSupplies":[18446744073709551615],"data":""}`)
	CallContract(t, ct, "mintBatch", payload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to mint 1 more — safeAdd overflow
	payload2 := []byte(`{"to":"hive:tibfox","ids":["1"],"amounts":[1],"maxSupplies":[18446744073709551615],"data":""}`)
	CallContract(t, ct, "mintBatch", payload2, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Transfer — Missing Payload Validation
// ===================================

func TestSafeTransferFromFailsWithEmptyFrom(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"","to":"hive:recipient","id":"1","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeTransferFromFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeBatchTransferFromFailsWithEmptyFrom(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"","to":"hive:recipient","ids":["1"],"amounts":[1],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeBatchTransferFromFailsWithEmptyTo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"hive:tibfox","to":"","ids":["1"],"amounts":[1],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeBatchTransferFromFailsWithEmptyIds(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","ids":[],"amounts":[],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Burn — Missing Payload Validation
// ===================================

func TestBurnFailsWithEmptyFrom(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"","id":"1","amount":1}`)
	CallContract(t, ct, "burn", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"hive:tibfox","id":"","amount":1}`)
	CallContract(t, ct, "burn", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnBatchFailsWithEmptyFrom(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"","ids":["1"],"amounts":[1]}`)
	CallContract(t, ct, "burnBatch", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnBatchFailsWithEmptyIds(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"hive:tibfox","ids":[],"amounts":[]}`)
	CallContract(t, ct, "burnBatch", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnBatchFailsWithZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":10,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"from":"hive:tibfox","ids":["1"],"amounts":[0]}`)
	CallContract(t, ct, "burnBatch", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// MintBatch — Missing Payload Validation
// ===================================

func TestMintBatchFailsWithEmptyTo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"to":"","ids":["1"],"amounts":[1],"maxSupplies":[1],"data":""}`)
	CallContract(t, ct, "mintBatch", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestMintBatchFailsWithMaxSuppliesLengthMismatch(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// 3 ids but 2 maxSupplies
	payload := []byte(`{"to":"hive:tibfox","ids":["1","2","3"],"amounts":[1,1,1],"maxSupplies":[1,1],"data":""}`)
	CallContract(t, ct, "mintBatch", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestMintBatchFailsWithZeroAmountInMiddle(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"to":"hive:tibfox","ids":["1","2","3"],"amounts":[1,0,1],"maxSupplies":[1,1,1],"data":""}`)
	CallContract(t, ct, "mintBatch", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestMintBatchFailsWithNoMaxSupplyForNewToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// maxSupplies omitted entirely — all tokens are new, should fail
	payload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[1,1],"data":""}`)
	CallContract(t, ct, "mintBatch", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// SetURI — Missing Payload Validation
// ===================================

func TestSetURIFailsWithEmptyTokenId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"id":"","uri":"https://example.com"}`)
	CallContract(t, ct, "setURI", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Query — Missing Payload Validation
// ===================================

func TestBalanceOfFailsWithEmptyAccount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"account":"","id":"1"}`)
	CallContract(t, ct, "balanceOf", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBalanceOfFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"account":"hive:tibfox","id":""}`)
	CallContract(t, ct, "balanceOf", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTotalSupplyFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"id":""}`)
	CallContract(t, ct, "totalSupply", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestMaxSupplyFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"id":""}`)
	CallContract(t, ct, "maxSupply", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestExistsFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"id":""}`)
	CallContract(t, ct, "exists", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTotalMintedFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"id":""}`)
	CallContract(t, ct, "totalMinted", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestIsSoulboundFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"id":""}`)
	CallContract(t, ct, "isSoulbound", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestURIFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"id":""}`)
	CallContract(t, ct, "uri", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestIsApprovedForAllFailsWithEmptyAccount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"account":"","operator":"hive:operator"}`)
	CallContract(t, ct, "isApprovedForAll", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestIsApprovedForAllFailsWithEmptyOperator(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"account":"hive:tibfox","operator":""}`)
	CallContract(t, ct, "isApprovedForAll", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// BalanceOfBatch — Missing Validation
// ===================================

func TestBalanceOfBatchFailsWithEmptyAccounts(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	payload := []byte(`{"accounts":[],"ids":[]}`)
	CallContract(t, ct, "balanceOfBatch", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// MintBatch — MaxSupply Mismatch on Existing Token
// ===================================

func TestMintBatchFailsWithMaxSupplyMismatchForExistingToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// First mint creates token "1" with maxSupply=100
	mint1 := []byte(`{"to":"hive:tibfox","id":"1","amount":10,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Second mintBatch with different maxSupply for existing token — should fail
	mint2 := []byte(`{"to":"hive:tibfox","ids":["1"],"amounts":[10],"maxSupplies":[200],"data":""}`)
	CallContract(t, ct, "mintBatch", mint2, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Pause — Read-Only Queries Still Work
// ===================================

func TestQueriesWorkWhilePaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a token before pausing
	mintPayload := []byte(`{"to":"hive:tibfox","id":"paused-query","amount":10,"maxSupply":100,"properties":{"color":"red"},"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Pause the contract
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	// All read-only queries should still work while paused
	balPayload := []byte(`{"account":"hive:tibfox","id":"paused-query"}`)
	balResult, _, _ := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":10}` {
		t.Errorf("Expected balance 10 while paused, got %s", balResult.Ret)
	}

	supplyPayload := []byte(`{"id":"paused-query"}`)
	supplyResult, _, _ := CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if supplyResult.Ret != `{"totalSupply":10}` {
		t.Errorf("Expected totalSupply 10 while paused, got %s", supplyResult.Ret)
	}

	maxPayload := []byte(`{"id":"paused-query"}`)
	maxResult, _, _ := CallContract(t, ct, "maxSupply", maxPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if maxResult.Ret != `{"maxSupply":100}` {
		t.Errorf("Expected maxSupply 100 while paused, got %s", maxResult.Ret)
	}

	existsPayload := []byte(`{"id":"paused-query"}`)
	existsResult, _, _ := CallContract(t, ct, "exists", existsPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if existsResult.Ret != `{"exists":true}` {
		t.Errorf("Expected exists true while paused, got %s", existsResult.Ret)
	}

	uriPayload := []byte(`{"id":"paused-query"}`)
	CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")

	propsPayload := []byte(`{"id":"paused-query"}`)
	propsResult, _, _ := CallContract(t, ct, "getProperties", propsPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if propsResult.Ret != `{"properties":{"color":"red"}}` {
		t.Errorf("Expected properties while paused, got %s", propsResult.Ret)
	}

	CallContract(t, ct, "getOwner", nil, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "getInfo", nil, nil, ownerAddress, true, uint(150_000_000), "")

	pausedResult, _, _ := CallContract(t, ct, "isPaused", nil, nil, ownerAddress, true, uint(150_000_000), "")
	if pausedResult.Ret != `{"paused":true}` {
		t.Errorf("Expected paused true, got %s", pausedResult.Ret)
	}
}

// ===================================
// Pause — Admin Actions Still Work
// ===================================

func TestAdminActionsWorkWhilePaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a token before pausing
	mintPayload := []byte(`{"to":"hive:tibfox","id":"paused-admin","amount":10,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Pause the contract
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	// setURI should work while paused (not a transfer operation)
	setURIPayload := []byte(`{"id":"paused-admin","uri":"https://custom.example.com/token.json"}`)
	CallContract(t, ct, "setURI", setURIPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// setBaseURI should work while paused
	setBaseURIPayload := []byte(`{"baseUri":"https://newbase.example.com/"}`)
	CallContract(t, ct, "setBaseURI", setBaseURIPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// setProperties should work while paused
	setPropsPayload := []byte(`{"id":"paused-admin","properties":{"updated":"while-paused"}}`)
	CallContract(t, ct, "setProperties", setPropsPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// setApprovalForAll should work while paused
	approvalPayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvalPayload, nil, ownerAddress, true, uint(150_000_000), "")
}

// ===================================
// Pause — ChangeOwner During Pause
// ===================================

func TestChangeOwnerWhilePausedNewOwnerUnpauses(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Pause the contract
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	// changeOwner should work while paused
	changePayload := []byte(`{"newOwner":"hive:newowner"}`)
	CallContract(t, ct, "changeOwner", changePayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Old owner cannot unpause anymore
	CallContract(t, ct, "unpause", nil, nil, ownerAddress, false, uint(150_000_000), "")

	// New owner can unpause
	CallContract(t, ct, "unpause", nil, nil, "hive:newowner", true, uint(150_000_000), "")

	// Verify operations work again after new owner unpause
	mintPayload := []byte(`{"to":"hive:tibfox","id":"after-unpause","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, "hive:newowner", true, uint(150_000_000), "")
}
