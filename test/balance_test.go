package contract_test

import (
	"testing"
)

// ===================================
// Balance Tests
// ===================================

func TestBalanceOf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	balancePayload := []byte(`{"account":"hive:tibfox","id":"1"}`)
	result, _, _ := CallContract(t, ct, "balanceOf", balancePayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":100}` {
		t.Errorf("Expected balance 100, got %s", result.Ret)
	}
}

func TestBalanceOfBatch(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	batchPayload := []byte(`{"accounts":["hive:tibfox","hive:tibfox"],"ids":["1","2"]}`)
	result, _, _ := CallContract(t, ct, "balanceOfBatch", batchPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balances":[100,200]}` {
		t.Errorf("Expected balances [100,200], got %s", result.Ret)
	}
}

func TestBalanceOfReturnsZeroForNonHolder(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check balance for someone who doesn't own the token
	balancePayload := []byte(`{"account":"hive:other","id":"1"}`)
	result, _, _ := CallContract(t, ct, "balanceOf", balancePayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":0}` {
		t.Errorf("Expected balance 0, got %s", result.Ret)
	}
}

func TestBalanceOfReturnsZeroForNonExistentToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	balancePayload := []byte(`{"account":"hive:tibfox","id":"nonexistent"}`)
	result, _, _ := CallContract(t, ct, "balanceOf", balancePayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":0}` {
		t.Errorf("Expected balance 0, got %s", result.Ret)
	}
}

func TestBalanceOfBatchWithDifferentAccounts(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint to different accounts
	mint1 := []byte(`{"to":"hive:alice","id":"1","amount":50,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	mint2 := []byte(`{"to":"hive:bob","id":"2","amount":75,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Query balances for different account/token combinations
	batchPayload := []byte(`{"accounts":["hive:alice","hive:bob"],"ids":["1","2"]}`)
	result, _, _ := CallContract(t, ct, "balanceOfBatch", batchPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balances":[50,75]}` {
		t.Errorf("Expected balances [50,75], got %s", result.Ret)
	}
}

func TestBalanceOfBatchFailsWithMismatchedLengths(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// 2 accounts, 3 ids
	batchPayload := []byte(`{"accounts":["hive:alice","hive:bob"],"ids":["1","2","3"]}`)
	CallContract(t, ct, "balanceOfBatch", batchPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBalanceAfterTransfer(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"1","amount":30,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check sender balance
	senderPayload := []byte(`{"account":"hive:tibfox","id":"1"}`)
	result, _, _ := CallContract(t, ct, "balanceOf", senderPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":70}` {
		t.Errorf("Expected sender balance 70, got %s", result.Ret)
	}

	// Check recipient balance
	recipientPayload := []byte(`{"account":"hive:recipient","id":"1"}`)
	result2, _, _ := CallContract(t, ct, "balanceOf", recipientPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"balance":30}` {
		t.Errorf("Expected recipient balance 30, got %s", result2.Ret)
	}
}

func TestBalanceAfterBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	burnPayload := []byte(`{"from":"hive:tibfox","id":"1","amount":40}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	balancePayload := []byte(`{"account":"hive:tibfox","id":"1"}`)
	result, _, _ := CallContract(t, ct, "balanceOf", balancePayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":60}` {
		t.Errorf("Expected balance 60, got %s", result.Ret)
	}
}
