package contract_test

import (
	"testing"
)

// ===================================
// ERC-165 supportsInterface Tests
// ===================================

func TestSupportsInterfaceERC1155(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check ERC-1155 interface (0xd9b67a26)
	payload := []byte(`{"interfaceId":"0xd9b67a26"}`)
	result, _, _ := CallContract(t, ct, "supportsInterface", payload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"supported":true}` {
		t.Errorf("Expected ERC-1155 interface supported, got %s", result.Ret)
	}
}

func TestSupportsInterfaceERC165(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check ERC-165 interface (0x01ffc9a7)
	payload := []byte(`{"interfaceId":"0x01ffc9a7"}`)
	result, _, _ := CallContract(t, ct, "supportsInterface", payload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"supported":true}` {
		t.Errorf("Expected ERC-165 interface supported, got %s", result.Ret)
	}
}

func TestSupportsInterfaceUnknown(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check unknown interface
	payload := []byte(`{"interfaceId":"0xdeadbeef"}`)
	result, _, _ := CallContract(t, ct, "supportsInterface", payload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"supported":false}` {
		t.Errorf("Expected unknown interface not supported, got %s", result.Ret)
	}
}

func TestSupportsInterfaceFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Empty interface ID should fail
	payload := []byte(`{"interfaceId":""}`)
	CallContract(t, ct, "supportsInterface", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSupportsInterfaceFailsWithoutInit(t *testing.T) {
	ct := SetupContractTest()

	// Should fail if contract not initialized
	payload := []byte(`{"interfaceId":"0xd9b67a26"}`)
	CallContract(t, ct, "supportsInterface", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSupportsInterfaceERC20(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check ERC-20 interface (0x36372b07) - should NOT be supported
	payload := []byte(`{"interfaceId":"0x36372b07"}`)
	result, _, _ := CallContract(t, ct, "supportsInterface", payload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"supported":false}` {
		t.Errorf("Expected ERC-20 interface not supported, got %s", result.Ret)
	}
}

func TestSupportsInterfaceERC721(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check ERC-721 interface (0x80ac58cd) - should NOT be supported
	payload := []byte(`{"interfaceId":"0x80ac58cd"}`)
	result, _, _ := CallContract(t, ct, "supportsInterface", payload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"supported":false}` {
		t.Errorf("Expected ERC-721 interface not supported, got %s", result.Ret)
	}
}
