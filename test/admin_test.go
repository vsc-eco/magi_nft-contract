package contract_test

import (
	"testing"
)

// ===================================
// Pause Tests
// ===================================

func TestPauseAndUnpause(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Pause
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	// Check paused
	result := CallContract(t, ct, "isPaused", nil, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"paused":true}` {
		t.Errorf("Expected paused true, got %s", result.Ret)
	}

	// Transfers should fail when paused
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, false, uint(150_000_000), "")

	// Unpause
	CallContract(t, ct, "unpause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	// Now mint should succeed
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
}

func TestPauseFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	CallContract(t, ct, "pause", nil, nil, "hive:other", false, uint(150_000_000), "")
}

func TestUnpauseFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	CallContract(t, ct, "unpause", nil, nil, "hive:other", false, uint(150_000_000), "")
}

func TestPauseFailsIfAlreadyPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "pause", nil, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to pause again
	CallContract(t, ct, "pause", nil, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestUnpauseFailsIfNotPaused(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to unpause when not paused
	CallContract(t, ct, "unpause", nil, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestIsPausedReturnsFalseByDefault(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	result := CallContract(t, ct, "isPaused", nil, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"paused":false}` {
		t.Errorf("Expected paused false, got %s", result.Ret)
	}
}

// ===================================
// Owner Tests
// ===================================

func TestChangeOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Change owner
	changePayload := []byte(`{"newOwner":"hive:newowner"}`)
	CallContract(t, ct, "changeOwner", changePayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check new owner
	result := CallContract(t, ct, "getOwner", nil, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"owner":"hive:newowner"}` {
		t.Errorf("Expected new owner, got %s", result.Ret)
	}

	// Old owner can't mint anymore
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, false, uint(150_000_000), "")

	// New owner can mint
	CallContract(t, ct, "mint", mintPayload, nil, "hive:newowner", true, uint(150_000_000), "")
}

func TestChangeOwnerFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	changePayload := []byte(`{"newOwner":"hive:newowner"}`)
	CallContract(t, ct, "changeOwner", changePayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestChangeOwnerFailsWithEmptyAddress(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	changePayload := []byte(`{"newOwner":""}`)
	CallContract(t, ct, "changeOwner", changePayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestChangeOwnerToSameOwnerFails(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Change owner to same address should fail
	changePayload := []byte(`{"newOwner":"hive:tibfox"}`)
	CallContract(t, ct, "changeOwner", changePayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestNewOwnerCanPause(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Change owner
	changePayload := []byte(`{"newOwner":"hive:newowner"}`)
	CallContract(t, ct, "changeOwner", changePayload, nil, ownerAddress, true, uint(150_000_000), "")

	// New owner can pause
	CallContract(t, ct, "pause", nil, nil, "hive:newowner", true, uint(150_000_000), "")

	// Old owner cannot unpause
	CallContract(t, ct, "unpause", nil, nil, ownerAddress, false, uint(150_000_000), "")

	// New owner can unpause
	CallContract(t, ct, "unpause", nil, nil, "hive:newowner", true, uint(150_000_000), "")
}

func TestNewOwnerCanChangeOwnerAgain(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// First ownership change
	change1 := []byte(`{"newOwner":"hive:owner2"}`)
	CallContract(t, ct, "changeOwner", change1, nil, ownerAddress, true, uint(150_000_000), "")

	// Second ownership change
	change2 := []byte(`{"newOwner":"hive:owner3"}`)
	CallContract(t, ct, "changeOwner", change2, nil, "hive:owner2", true, uint(150_000_000), "")

	// Verify final owner
	result := CallContract(t, ct, "getOwner", nil, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"owner":"hive:owner3"}` {
		t.Errorf("Expected owner3, got %s", result.Ret)
	}
}
