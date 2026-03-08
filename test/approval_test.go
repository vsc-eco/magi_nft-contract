package contract_test

import (
	"testing"
)

// ===================================
// Approval Tests
// ===================================

func TestSetApprovalForAll(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	approvalPayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvalPayload, nil, ownerAddress, true, uint(150_000_000), "")

	checkPayload := []byte(`{"account":"hive:tibfox","operator":"hive:operator"}`)
	result := CallContract(t, ct, "isApprovedForAll", checkPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"approved":true}` {
		t.Errorf("Expected approved true, got %s", result.Ret)
	}
}

func TestRevokeApprovalForAll(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Approve
	approvePayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvePayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify approved
	checkPayload := []byte(`{"account":"hive:tibfox","operator":"hive:operator"}`)
	result := CallContract(t, ct, "isApprovedForAll", checkPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"approved":true}` {
		t.Errorf("Expected approved true, got %s", result.Ret)
	}

	// Revoke
	revokePayload := []byte(`{"operator":"hive:operator","approved":false}`)
	CallContract(t, ct, "setApprovalForAll", revokePayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify revoked
	result2 := CallContract(t, ct, "isApprovedForAll", checkPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"approved":false}` {
		t.Errorf("Expected approved false, got %s", result2.Ret)
	}
}

func TestApprovedOperatorCanTransfer(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Approve operator
	approvalPayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvalPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Operator transfers tokens
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"1","amount":50,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, "hive:operator", true, uint(150_000_000), "")
}

func TestRevokedOperatorCannotTransfer(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Approve operator
	approvalPayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvalPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Revoke operator
	revokePayload := []byte(`{"operator":"hive:operator","approved":false}`)
	CallContract(t, ct, "setApprovalForAll", revokePayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Operator tries to transfer - should fail
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"1","amount":50,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, "hive:operator", false, uint(150_000_000), "")
}

func TestApprovedOperatorCanBatchTransfer(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["1","2"],"amounts":[100,200],"maxSupplies":[100,200],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Approve operator
	approvalPayload := []byte(`{"operator":"hive:operator","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvalPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Operator batch transfers tokens
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","ids":["1","2"],"amounts":[50,100],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, "hive:operator", true, uint(150_000_000), "")
}

func TestIsApprovedForAllReturnsFalseForNonApproved(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check approval for operator that was never approved
	checkPayload := []byte(`{"account":"hive:tibfox","operator":"hive:unknown"}`)
	result := CallContract(t, ct, "isApprovedForAll", checkPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"approved":false}` {
		t.Errorf("Expected approved false, got %s", result.Ret)
	}
}

func TestSetApprovalForAllFailsWithEmptyOperator(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	approvalPayload := []byte(`{"operator":"","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvalPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestMultipleOperatorsCanBeApproved(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Approve operator1
	approve1 := []byte(`{"operator":"hive:operator1","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approve1, nil, ownerAddress, true, uint(150_000_000), "")

	// Approve operator2
	approve2 := []byte(`{"operator":"hive:operator2","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approve2, nil, ownerAddress, true, uint(150_000_000), "")

	// Both should be approved
	check1 := []byte(`{"account":"hive:tibfox","operator":"hive:operator1"}`)
	result1 := CallContract(t, ct, "isApprovedForAll", check1, nil, ownerAddress, true, uint(150_000_000), "")
	if result1.Ret != `{"approved":true}` {
		t.Errorf("Expected operator1 approved true, got %s", result1.Ret)
	}

	check2 := []byte(`{"account":"hive:tibfox","operator":"hive:operator2"}`)
	result2 := CallContract(t, ct, "isApprovedForAll", check2, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"approved":true}` {
		t.Errorf("Expected operator2 approved true, got %s", result2.Ret)
	}
}

func TestSetApprovalForAllFailsIfApprovingSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to approve self as operator
	approvalPayload := []byte(`{"operator":"hive:tibfox","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvalPayload, nil, ownerAddress, false, uint(150_000_000), "")
}
