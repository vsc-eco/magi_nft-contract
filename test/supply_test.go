package contract_test

import (
	"testing"
)

// ===================================
// Supply Query Tests
// ===================================

func TestTotalSupplyQuery(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":50,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	supplyPayload := []byte(`{"id":"1"}`)
	result := CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"totalSupply":50}` {
		t.Errorf("Expected totalSupply 50, got %s", result.Ret)
	}
}

func TestMaxSupplyQuery(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	mintPayload := []byte(`{"to":"hive:tibfox","id":"1","amount":50,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	supplyPayload := []byte(`{"id":"1"}`)
	result := CallContract(t, ct, "maxSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"maxSupply":100}` {
		t.Errorf("Expected maxSupply 100, got %s", result.Ret)
	}
}

func TestTotalSupplyReturnsZeroForNonExistentToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	supplyPayload := []byte(`{"id":"nonexistent"}`)
	result := CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"totalSupply":0}` {
		t.Errorf("Expected totalSupply 0, got %s", result.Ret)
	}
}

func TestMaxSupplyReturnsZeroForNonExistentToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	supplyPayload := []byte(`{"id":"nonexistent"}`)
	result := CallContract(t, ct, "maxSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"maxSupply":0}` {
		t.Errorf("Expected maxSupply 0, got %s", result.Ret)
	}
}

func TestSupplyAfterMultipleMints(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// First mint
	mint1 := []byte(`{"to":"hive:tibfox","id":"multi","amount":30,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Second mint
	mint2 := []byte(`{"to":"hive:other","id":"multi","amount":25,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Third mint
	mint3 := []byte(`{"to":"hive:third","id":"multi","amount":15,"data":""}`)
	CallContract(t, ct, "mint", mint3, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify total supply is 70
	supplyPayload := []byte(`{"id":"multi"}`)
	result := CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"totalSupply":70}` {
		t.Errorf("Expected totalSupply 70, got %s", result.Ret)
	}

	// Max supply should still be 100
	result2 := CallContract(t, ct, "maxSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"maxSupply":100}` {
		t.Errorf("Expected maxSupply 100, got %s", result2.Ret)
	}
}

func TestSupplyAfterBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	mintPayload := []byte(`{"to":"hive:tibfox","id":"burntest","amount":100,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	burnPayload := []byte(`{"from":"hive:tibfox","id":"burntest","amount":40}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Total supply should be 60
	supplyPayload := []byte(`{"id":"burntest"}`)
	result := CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"totalSupply":60}` {
		t.Errorf("Expected totalSupply 60, got %s", result.Ret)
	}

	// Max supply should still be 100
	result2 := CallContract(t, ct, "maxSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"maxSupply":100}` {
		t.Errorf("Expected maxSupply 100, got %s", result2.Ret)
	}
}

// ===================================
// Exists Query Tests
// ===================================

func TestExistsReturnsTrueForMintedToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"exists-test","amount":10,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check exists
	existsPayload := []byte(`{"id":"exists-test"}`)
	result := CallContract(t, ct, "exists", existsPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"exists":true}` {
		t.Errorf("Expected exists true, got %s", result.Ret)
	}
}

func TestExistsReturnsFalseForNonExistentToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Check exists for token that was never minted
	existsPayload := []byte(`{"id":"never-minted"}`)
	result := CallContract(t, ct, "exists", existsPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"exists":false}` {
		t.Errorf("Expected exists false, got %s", result.Ret)
	}
}

func TestExistsReturnsTrueAfterFullBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"burn-exists","amount":10,"maxSupply":10,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Burn all tokens
	burnPayload := []byte(`{"from":"hive:tibfox","id":"burn-exists","amount":10}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Token should still exist (maxSupply was set)
	existsPayload := []byte(`{"id":"burn-exists"}`)
	result := CallContract(t, ct, "exists", existsPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"exists":true}` {
		t.Errorf("Expected exists true after burn (token type was created), got %s", result.Ret)
	}

	// But totalSupply should be 0
	supplyPayload := []byte(`{"id":"burn-exists"}`)
	result2 := CallContract(t, ct, "totalSupply", supplyPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"totalSupply":0}` {
		t.Errorf("Expected totalSupply 0, got %s", result2.Ret)
	}
}
