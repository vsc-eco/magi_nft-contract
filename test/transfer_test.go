package contract_test

import (
	"testing"
)

// ===================================
// Transfer Tests
// ===================================

func TestSafeTransferFromSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"1","amount":50,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check balances
	balancePayload := []byte(`{"account":"hive:tibfox","id":"1"}`)
	result := CallContract(t, ct, "balanceOf", balancePayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":50}` {
		t.Errorf("Expected sender balance 50, got %s", result.Ret)
	}

	recipientPayload := []byte(`{"account":"hive:recipient","id":"1"}`)
	result = CallContract(t, ct, "balanceOf", recipientPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":50}` {
		t.Errorf("Expected recipient balance 50, got %s", result.Ret)
	}
}

func TestSafeTransferFromFailsIfNotAuthorized(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Other user trying to transfer tibfox's tokens
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:other","id":"1","amount":50,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestSafeTransferFromFailsIfPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"1","amount":50,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeTransferFromFailsWithInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":50,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to transfer more than available
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"1","amount":100,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeTransferFromWithZeroAmountFails(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Zero amount transfer should fail
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"1","amount":0,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeTransferFromFailsWithEmptyRecipient(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	transferPayload := []byte(`{"from":"hive:tibfox","to":"","id":"1","amount":50,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeTransferFromSelfToSelfFails(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Transfer to self should fail
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:tibfox","id":"1","amount":50,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// SafeBatchTransferFrom Tests
// ===================================

func TestSafeBatchTransferFromSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","ids":["1","2"],"amounts":[30,50],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")
}

func TestSafeBatchTransferFromFailsIfNotAuthorized(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","ids":["1","2"],"amounts":[30,50],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestSafeBatchTransferFromFailsWithMismatchedLengths(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// 2 ids, 3 amounts
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","ids":["1","2"],"amounts":[30,50,70],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeBatchTransferFromFailsIfPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","ids":["1","2"],"amounts":[30,50],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSafeBatchTransferFromFailsWithInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[50,100],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to transfer more than available for token 1
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","ids":["1","2"],"amounts":[100,50],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, ownerAddress, false, uint(150_000_000), "")
}
