package contract_test

import (
	"testing"
)

// ===================================
// Per-Token Approval Tests (ERC-6909)
// ===================================

func TestApproveAndTransfer(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"card-1","amount":10,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Approve spender for 5 of card-1
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"card-1","amount":5}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Check allowance
	result := CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:marketplace","id":"card-1"}`), nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"amount":5}` {
		t.Errorf("Expected allowance 5, got %s", result.Ret)
	}

	// Spender transfers 3 (within allowance)
	CallContract(t, ct, "safeTransferFrom", []byte(`{"from":"hive:tibfox","to":"hive:buyer","id":"card-1","amount":3,"data":""}`), nil, "hive:marketplace", true, uint(150_000_000), "")

	// Allowance should be decremented to 2
	result = CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:marketplace","id":"card-1"}`), nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"amount":2}` {
		t.Errorf("Expected allowance 2 after transfer, got %s", result.Ret)
	}

	// Verify balance
	result = CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:buyer","id":"card-1"}`), nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"balance":3}` {
		t.Errorf("Expected balance 3, got %s", result.Ret)
	}
}

func TestApproveExceedsAllowanceFails(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"card-1","amount":10,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Approve for 2
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"card-1","amount":2}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Try to transfer 5 — exceeds allowance
	CallContract(t, ct, "safeTransferFrom", []byte(`{"from":"hive:tibfox","to":"hive:buyer","id":"card-1","amount":5,"data":""}`), nil, "hive:marketplace", false, uint(150_000_000), "")
}

func TestApproveZeroRevokes(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"card-1","amount":10,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Approve then revoke
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"card-1","amount":5}`), nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"card-1","amount":0}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Transfer should fail
	CallContract(t, ct, "safeTransferFrom", []byte(`{"from":"hive:tibfox","to":"hive:buyer","id":"card-1","amount":1,"data":""}`), nil, "hive:marketplace", false, uint(150_000_000), "")
}

func TestApproveDoesNotAffectOtherTokens(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"card-1","amount":10,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"card-2","amount":10,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Approve for card-1 only
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"card-1","amount":5}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Transfer card-2 should fail (no allowance)
	CallContract(t, ct, "safeTransferFrom", []byte(`{"from":"hive:tibfox","to":"hive:buyer","id":"card-2","amount":1,"data":""}`), nil, "hive:marketplace", false, uint(150_000_000), "")
}

func TestApproveBatchTransfer(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mintBatch", []byte(`{"to":"hive:tibfox","ids":["a","b"],"amounts":[10,10],"maxSupplies":[10,10],"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Approve both tokens
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"a","amount":5}`), nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"b","amount":3}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Batch transfer within allowance
	CallContract(t, ct, "safeBatchTransferFrom", []byte(`{"from":"hive:tibfox","to":"hive:buyer","ids":["a","b"],"amounts":[2,3],"data":""}`), nil, "hive:marketplace", true, uint(150_000_000), "")

	// Check remaining allowances
	result := CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:marketplace","id":"a"}`), nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"amount":3}` {
		t.Errorf("Expected allowance 3 for a, got %s", result.Ret)
	}
	result = CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:marketplace","id":"b"}`), nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"amount":0}` {
		t.Errorf("Expected allowance 0 for b, got %s", result.Ret)
	}
}

func TestApproveBatchTransferFailsIfOneExceeds(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mintBatch", []byte(`{"to":"hive:tibfox","ids":["a","b"],"amounts":[10,10],"maxSupplies":[10,10],"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"a","amount":5}`), nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"b","amount":1}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Batch transfer — b exceeds allowance
	CallContract(t, ct, "safeBatchTransferFrom", []byte(`{"from":"hive:tibfox","to":"hive:buyer","ids":["a","b"],"amounts":[2,3],"data":""}`), nil, "hive:marketplace", false, uint(150_000_000), "")
}

func TestApproveFailsSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	CallContract(t, ct, "approve", []byte(`{"spender":"hive:tibfox","id":"card-1","amount":5}`), nil, ownerAddress, false, uint(150_000_000), "")
}

func TestApproveFailsEmptySpender(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	CallContract(t, ct, "approve", []byte(`{"spender":"","id":"card-1","amount":5}`), nil, ownerAddress, false, uint(150_000_000), "")
}

func TestApproveFailsEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"","amount":5}`), nil, ownerAddress, false, uint(150_000_000), "")
}

func TestOperatorBypassesAllowance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"card-1","amount":10,"maxSupply":10,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Set blanket operator approval (no per-token allowance needed)
	CallContract(t, ct, "setApprovalForAll", []byte(`{"operator":"hive:operator","approved":true}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Operator can transfer without per-token allowance
	CallContract(t, ct, "safeTransferFrom", []byte(`{"from":"hive:tibfox","to":"hive:buyer","id":"card-1","amount":5,"data":""}`), nil, "hive:operator", true, uint(150_000_000), "")
}

func TestAllowanceDefaultsToZero(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	result := CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:marketplace","id":"card-1"}`), nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"amount":0}` {
		t.Errorf("Expected allowance 0, got %s", result.Ret)
	}
}

func TestAllowanceOnNonExistentToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Query allowance for a token that was never minted — should return 0
	result := CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:marketplace","id":"does-not-exist"}`), nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"amount":0}` {
		t.Errorf("Expected allowance 0 for non-existent token, got %s", result.Ret)
	}
}

func TestApproveSoulboundTokenSucceeds(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"to":"hive:tibfox","id":"soul-1","amount":1,"maxSupply":1,"soulbound":true,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Approve should succeed (soulbound blocks transfer, not approval)
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"soul-1","amount":1}`), nil, ownerAddress, true, uint(150_000_000), "")

	result := CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:marketplace","id":"soul-1"}`), nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"amount":1}` {
		t.Errorf("Expected allowance 1, got %s", result.Ret)
	}
}

func TestTransferSoulboundViaAllowanceFails(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint soulbound to a non-owner account, then transfer to them via owner
	CallContract(t, ct, "mint", []byte(`{"to":"hive:holder","id":"soul-1","amount":1,"maxSupply":1,"soulbound":true,"data":""}`), nil, ownerAddress, true, uint(150_000_000), "")

	// Holder approves marketplace
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:marketplace","id":"soul-1","amount":1}`), nil, "hive:holder", true, uint(150_000_000), "")

	// Marketplace tries to transfer holder's soulbound token — should fail (holder is not contract owner)
	CallContract(t, ct, "safeTransferFrom", []byte(`{"from":"hive:holder","to":"hive:buyer","id":"soul-1","amount":1,"data":""}`), nil, "hive:marketplace", false, uint(150_000_000), "")
}
