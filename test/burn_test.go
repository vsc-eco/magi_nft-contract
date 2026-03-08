package contract_test

import (
	"testing"
)

// ===================================
// Burn Tests
// ===================================

func TestBurnSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	burnPayload := []byte(`{"from":"hive:tibfox","id":"1","amount":30}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	balancePayload := []byte(`{"account":"hive:tibfox","id":"1"}`)
	result := CallContract(t, ct, "balanceOf", balancePayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":70}` {
		t.Errorf("Expected balance 70, got %s", result.Ret)
	}

	// Check total supply decreased
	supplyPayload := []byte(`{"id":"1"}`)
	result = CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"totalSupply":70}` {
		t.Errorf("Expected totalSupply 70, got %s", result.Ret)
	}
}

func TestBurnFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Someone else tries to burn tibfox's tokens
	burnPayload := []byte(`{"from":"hive:tibfox","id":"1","amount":30}`)
	CallContract(t, ct, "burn", burnPayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestBurnFailsIfPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	burnPayload := []byte(`{"from":"hive:tibfox","id":"1","amount":30}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnFailsWithInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":50,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to burn more than available
	burnPayload := []byte(`{"from":"hive:tibfox","id":"1","amount":100}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnFailsWithZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	burnPayload := []byte(`{"from":"hive:tibfox","id":"1","amount":0}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnEntireBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Burn all 100
	burnPayload := []byte(`{"from":"hive:tibfox","id":"1","amount":100}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Balance should be 0
	balancePayload := []byte(`{"account":"hive:tibfox","id":"1"}`)
	result := CallContract(t, ct, "balanceOf", balancePayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":0}` {
		t.Errorf("Expected balance 0, got %s", result.Ret)
	}

	// Total supply should be 0
	supplyPayload := []byte(`{"id":"1"}`)
	result2 := CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"totalSupply":0}` {
		t.Errorf("Expected totalSupply 0, got %s", result2.Ret)
	}
}

func TestBurnAndRemint(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint 100 of 100
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Burn 50
	burnPayload := []byte(`{"from":"hive:tibfox","id":"1","amount":50}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Re-mint 50 (should succeed since totalSupply is now 50)
	remintPayload := []byte(`{"to":"hive:other","id":"1","amount":50,"data":""}`)
	CallContract(t, ct, "mint", remintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Total supply should be back to 100
	supplyPayload := []byte(`{"id":"1"}`)
	result := CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"totalSupply":100}` {
		t.Errorf("Expected totalSupply 100, got %s", result.Ret)
	}
}

// ===================================
// BurnBatch Tests
// ===================================

func TestBurnBatchSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	burnPayload := []byte(`{"from":"hive:tibfox","ids":["1","2"],"amounts":[30,50]}`)
	CallContract(t, ct, "burnBatch", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")
}

func TestBurnBatchFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	burnPayload := []byte(`{"from":"hive:tibfox","ids":["1","2"],"amounts":[30,50]}`)
	CallContract(t, ct, "burnBatch", burnPayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestBurnBatchFailsIfPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	burnPayload := []byte(`{"from":"hive:tibfox","ids":["1","2"],"amounts":[30,50]}`)
	CallContract(t, ct, "burnBatch", burnPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnBatchFailsWithMismatchedLengths(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// 2 ids, 3 amounts
	burnPayload := []byte(`{"from":"hive:tibfox","ids":["1","2"],"amounts":[30,50,70]}`)
	CallContract(t, ct, "burnBatch", burnPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestBurnBatchFailsWithInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[50,100],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to burn more than available for token 1
	burnPayload := []byte(`{"from":"hive:tibfox","ids":["1","2"],"amounts":[100,50]}`)
	CallContract(t, ct, "burnBatch", burnPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Operator Burn Tests
// ===================================

func TestOperatorCanBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to a user
	mintPayload := []byte(`{"to":"hive:user1","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// User approves operator
	approvePayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvePayload, nil, "hive:user1", true, uint(150_000_000), "")

	// Operator burns tokens on behalf of user
	burnPayload := []byte(`{"from":"hive:user1","id":"1","amount":30}`)
	result := CallContract(t, ct, "burn", burnPayload, nil, "hive:operator", true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected operator burn success, got %s", result.Ret)
	}

	// Verify balance decreased
	balPayload := []byte(`{"account":"hive:user1","id":"1"}`)
	balResult := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":70}` {
		t.Errorf("Expected balance 70, got %s", balResult.Ret)
	}
}

func TestOperatorCanBatchBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to a user
	mintPayload := []byte(`{"to":"hive:user1","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// User approves operator
	approvePayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvePayload, nil, "hive:user1", true, uint(150_000_000), "")

	// Operator batch burns tokens on behalf of user
	burnPayload := []byte(`{"from":"hive:user1","ids":["1","2"],"amounts":[30,50]}`)
	result := CallContract(t, ct, "burnBatch", burnPayload, nil, "hive:operator", true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected operator batch burn success, got %s", result.Ret)
	}

	// Verify balances decreased
	bal1 := []byte(`{"account":"hive:user1","id":"1"}`)
	bal1Result := CallContract(t, ct, "balanceOf", bal1, nil, ownerAddress, true, uint(150_000_000), "")
	if bal1Result.Ret != `{"balance":70}` {
		t.Errorf("Expected balance 70, got %s", bal1Result.Ret)
	}

	bal2 := []byte(`{"account":"hive:user1","id":"2"}`)
	bal2Result := CallContract(t, ct, "balanceOf", bal2, nil, ownerAddress, true, uint(150_000_000), "")
	if bal2Result.Ret != `{"balance":150}` {
		t.Errorf("Expected balance 150, got %s", bal2Result.Ret)
	}
}

func TestRevokedOperatorCannotBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to a user
	mintPayload := []byte(`{"to":"hive:user1","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// User approves operator
	approvePayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvePayload, nil, "hive:user1", true, uint(150_000_000), "")

	// User revokes operator
	revokePayload := []byte(`{"operator":"hive:operator","approved":false}`)
	CallContract(t, ct, "setApprovalForAll", revokePayload, nil, "hive:user1", true, uint(150_000_000), "")

	// Operator tries to burn - should fail
	burnPayload := []byte(`{"from":"hive:user1","id":"1","amount":30}`)
	CallContract(t, ct, "burn", burnPayload, nil, "hive:operator", false, uint(150_000_000), "")
}

// ===================================
// User Self-Burn Tests
// ===================================

func TestRegularUserCanBurnOwnTokens(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to a regular user (not the contract owner)
	mintPayload := []byte(`{"to":"hive:regularuser","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Regular user burns their own tokens
	burnPayload := []byte(`{"from":"hive:regularuser","id":"1","amount":40}`)
	result := CallContract(t, ct, "burn", burnPayload, nil, "hive:regularuser", true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected user to burn their own tokens, got %s", result.Ret)
	}

	// Verify balance
	balPayload := []byte(`{"account":"hive:regularuser","id":"1"}`)
	balResult := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":60}` {
		t.Errorf("Expected balance 60, got %s", balResult.Ret)
	}
}

func TestRegularUserCanBatchBurnOwnTokens(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to a regular user
	mintPayload := []byte(`{"to":"hive:regularuser","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Regular user batch burns their own tokens
	burnPayload := []byte(`{"from":"hive:regularuser","ids":["1","2"],"amounts":[25,75]}`)
	result := CallContract(t, ct, "burnBatch", burnPayload, nil, "hive:regularuser", true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected user batch burn success, got %s", result.Ret)
	}

	// Verify balances
	bal1 := []byte(`{"account":"hive:regularuser","id":"1"}`)
	bal1Result := CallContract(t, ct, "balanceOf", bal1, nil, ownerAddress, true, uint(150_000_000), "")
	if bal1Result.Ret != `{"balance":75}` {
		t.Errorf("Expected balance 75, got %s", bal1Result.Ret)
	}

	bal2 := []byte(`{"account":"hive:regularuser","id":"2"}`)
	bal2Result := CallContract(t, ct, "balanceOf", bal2, nil, ownerAddress, true, uint(150_000_000), "")
	if bal2Result.Ret != `{"balance":125}` {
		t.Errorf("Expected balance 125, got %s", bal2Result.Ret)
	}
}

// ===================================
// Contract Owner Burn Restriction Tests
// ===================================

func TestContractOwnerCannotBurnOthersTokens(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to a user
	mintPayload := []byte(`{"to":"hive:user1","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Contract owner tries to burn user's tokens without approval - should fail
	burnPayload := []byte(`{"from":"hive:user1","id":"1","amount":30}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, false, uint(150_000_000), "")

	// Verify balance unchanged
	balPayload := []byte(`{"account":"hive:user1","id":"1"}`)
	balResult := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":100}` {
		t.Errorf("Expected balance unchanged at 100, got %s", balResult.Ret)
	}
}

func TestContractOwnerCannotBatchBurnOthersTokens(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to a user
	mintPayload := []byte(`{"to":"hive:user1","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Contract owner tries to batch burn user's tokens without approval - should fail
	burnPayload := []byte(`{"from":"hive:user1","ids":["1","2"],"amounts":[30,50]}`)
	CallContract(t, ct, "burnBatch", burnPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestContractOwnerCanBurnIfApproved(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to a user
	mintPayload := []byte(`{"to":"hive:user1","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// User approves contract owner as operator
	approvePayload := []byte(`{"operator":"hive:tibfox","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvePayload, nil, "hive:user1", true, uint(150_000_000), "")

	// Now contract owner can burn user's tokens
	burnPayload := []byte(`{"from":"hive:user1","id":"1","amount":30}`)
	result := CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected approved contract owner to burn, got %s", result.Ret)
	}

	// Verify balance
	balPayload := []byte(`{"account":"hive:user1","id":"1"}`)
	balResult := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":70}` {
		t.Errorf("Expected balance 70, got %s", balResult.Ret)
	}
}
