package contract_test

import (
	"testing"
)

// ===================================
// Properties Tests - Mint with Properties
// ===================================

func TestMintWithProperties(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint with properties
	mintPayload := []byte(`{"to":"hive:tibfox","id":"prop-1","amount":1,"maxSupply":1,"properties":{"color":"red","rarity":"legendary"},"data":""}`)
	result, _, _ := CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected success, got %s", result.Ret)
	}

	// Verify properties can be read back
	getPayload := []byte(`{"id":"prop-1"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"color":"red","rarity":"legendary"}}` {
		t.Errorf("Expected properties, got %s", getResult.Ret)
	}
}

func TestMintWithoutProperties(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint without properties
	mintPayload := []byte(`{"to":"hive:tibfox","id":"no-props","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Properties should be null
	getPayload := []byte(`{"id":"no-props"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":null}` {
		t.Errorf("Expected null properties, got %s", getResult.Ret)
	}
}

func TestMintWithComplexProperties(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint with nested/complex properties
	mintPayload := []byte(`{"to":"hive:tibfox","id":"complex-props","amount":1,"maxSupply":1,"properties":{"name":"Sword of Fire","stats":{"attack":50,"defense":10},"tags":["fire","rare"],"level":5},"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(250_000_000), "")

	getPayload := []byte(`{"id":"complex-props"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"name":"Sword of Fire","stats":{"attack":50,"defense":10},"tags":["fire","rare"],"level":5}}` {
		t.Errorf("Expected complex properties, got %s", getResult.Ret)
	}
}

func TestPropertiesOnlySetOnFirstMint(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// First mint sets properties
	mint1 := []byte(`{"to":"hive:tibfox","id":"first-props","amount":5,"maxSupply":100,"properties":{"color":"blue"},"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Second mint with different properties - should be ignored (not first mint)
	mint2 := []byte(`{"to":"hive:tibfox","id":"first-props","amount":5,"properties":{"color":"red"},"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Properties should still be from first mint
	getPayload := []byte(`{"id":"first-props"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"color":"blue"}}` {
		t.Errorf("Expected first mint properties, got %s", getResult.Ret)
	}
}

// ===================================
// Properties Tests - MintBatch with Properties
// ===================================

func TestMintBatchWithProperties(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint batch with per-token properties
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["batch-p1","batch-p2"],"amounts":[1,1],"maxSupplies":[1,1],"properties":[{"color":"red"},{"color":"blue","size":42}],"data":""}`)
	result, _, _ := CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(200_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected success, got %s", result.Ret)
	}

	// Verify each token's properties
	get1 := []byte(`{"id":"batch-p1"}`)
	get1Result, _, _ := CallContract(t, ct, "getProperties", get1, nil, ownerAddress, true, uint(150_000_000), "")
	if get1Result.Ret != `{"properties":{"color":"red"}}` {
		t.Errorf("Expected batch-p1 properties, got %s", get1Result.Ret)
	}

	get2 := []byte(`{"id":"batch-p2"}`)
	get2Result, _, _ := CallContract(t, ct, "getProperties", get2, nil, ownerAddress, true, uint(150_000_000), "")
	if get2Result.Ret != `{"properties":{"color":"blue","size":42}}` {
		t.Errorf("Expected batch-p2 properties, got %s", get2Result.Ret)
	}
}

func TestMintBatchWithoutProperties(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint batch without properties array
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["batch-np1","batch-np2"],"amounts":[1,1],"maxSupplies":[1,1],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Both should have null properties
	get1 := []byte(`{"id":"batch-np1"}`)
	get1Result, _, _ := CallContract(t, ct, "getProperties", get1, nil, ownerAddress, true, uint(150_000_000), "")
	if get1Result.Ret != `{"properties":null}` {
		t.Errorf("Expected null properties, got %s", get1Result.Ret)
	}

	get2 := []byte(`{"id":"batch-np2"}`)
	get2Result, _, _ := CallContract(t, ct, "getProperties", get2, nil, ownerAddress, true, uint(150_000_000), "")
	if get2Result.Ret != `{"properties":null}` {
		t.Errorf("Expected null properties, got %s", get2Result.Ret)
	}
}

// ===================================
// Properties Tests - SetProperties
// ===================================

func TestSetPropertiesSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint without properties
	mintPayload := []byte(`{"to":"hive:tibfox","id":"set-props","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Set properties after mint
	setPayload := []byte(`{"id":"set-props","properties":{"color":"green","level":3}}`)
	result, _, _ := CallContract(t, ct, "setProperties", setPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected success, got %s", result.Ret)
	}

	// Verify properties
	getPayload := []byte(`{"id":"set-props"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"color":"green","level":3}}` {
		t.Errorf("Expected set properties, got %s", getResult.Ret)
	}
}

func TestSetPropertiesOverwrite(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint with initial properties
	mintPayload := []byte(`{"to":"hive:tibfox","id":"overwrite-props","amount":1,"maxSupply":1,"properties":{"version":1},"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Overwrite properties
	setPayload := []byte(`{"id":"overwrite-props","properties":{"version":2,"updated":true}}`)
	CallContract(t, ct, "setProperties", setPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify overwritten properties
	getPayload := []byte(`{"id":"overwrite-props"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"version":2,"updated":true}}` {
		t.Errorf("Expected overwritten properties, got %s", getResult.Ret)
	}
}

func TestSetPropertiesFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"owner-only-props","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Non-owner tries to set properties - should fail
	setPayload := []byte(`{"id":"owner-only-props","properties":{"hacked":true}}`)
	CallContract(t, ct, "setProperties", setPayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestSetPropertiesFailsIfNotInit(t *testing.T) {
	ct := SetupContractTest()

	// Try to set properties without init - should fail
	setPayload := []byte(`{"id":"1","properties":{"test":true}}`)
	CallContract(t, ct, "setProperties", setPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetPropertiesFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	setPayload := []byte(`{"id":"","properties":{"test":true}}`)
	CallContract(t, ct, "setProperties", setPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetPropertiesFailsWithEmptyProperties(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"empty-props","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to set empty properties - should fail
	setPayload := []byte(`{"id":"empty-props"}`)
	CallContract(t, ct, "setProperties", setPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestSetPropertiesFailsWithEmptyPayload(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	CallContract(t, ct, "setProperties", nil, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Properties Tests - GetProperties
// ===================================

func TestGetPropertiesNonExistentToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Query properties for a token that doesn't exist
	getPayload := []byte(`{"id":"nonexistent"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":null}` {
		t.Errorf("Expected null properties for non-existent token, got %s", getResult.Ret)
	}
}

func TestGetPropertiesFailsIfNotInit(t *testing.T) {
	ct := SetupContractTest()

	getPayload := []byte(`{"id":"1"}`)
	CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestGetPropertiesFailsWithEmptyId(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	getPayload := []byte(`{"id":""}`)
	CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestGetPropertiesFailsWithEmptyPayload(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	CallContract(t, ct, "getProperties", nil, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Properties Tests - Persistence
// ===================================

func TestPropertiesPersistAfterBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint with properties
	mintPayload := []byte(`{"to":"hive:tibfox","id":"burn-props","amount":5,"maxSupply":10,"properties":{"rarity":"epic"},"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Burn all tokens
	burnPayload := []byte(`{"from":"hive:tibfox","id":"burn-props","amount":5}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Properties should still be readable (they belong to the token type, not individual tokens)
	getPayload := []byte(`{"id":"burn-props"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"rarity":"epic"}}` {
		t.Errorf("Expected properties to persist after burn, got %s", getResult.Ret)
	}
}

func TestPropertiesWithSoulbound(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint soulbound token with properties
	mintPayload := []byte(`{"to":"hive:tibfox","id":"sb-props","amount":1,"maxSupply":1,"soulbound":true,"properties":{"type":"achievement","title":"First Kill"},"data":""}`)
	result, _, _ := CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected success, got %s", result.Ret)
	}

	// Verify both soulbound and properties
	sbPayload := []byte(`{"id":"sb-props"}`)
	sbResult, _, _ := CallContract(t, ct, "isSoulbound", sbPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if sbResult.Ret != `{"soulbound":true}` {
		t.Errorf("Expected soulbound true, got %s", sbResult.Ret)
	}

	getPayload := []byte(`{"id":"sb-props"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"type":"achievement","title":"First Kill"}}` {
		t.Errorf("Expected properties, got %s", getResult.Ret)
	}
}

func TestSetPropertiesOnNonExistentToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Owner can set properties on any token ID (even unminted ones - flexible design)
	setPayload := []byte(`{"id":"future-token","properties":{"preloaded":true}}`)
	result, _, _ := CallContract(t, ct, "setProperties", setPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected success, got %s", result.Ret)
	}

	// Properties should be readable
	getPayload := []byte(`{"id":"future-token"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"preloaded":true}}` {
		t.Errorf("Expected preloaded properties, got %s", getResult.Ret)
	}
}

func TestGetPropertiesAnyoneCanRead(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint with properties
	mintPayload := []byte(`{"to":"hive:tibfox","id":"public-props","amount":1,"maxSupply":1,"properties":{"visible":true},"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Non-owner can read properties
	getPayload := []byte(`{"id":"public-props"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, "hive:other", true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"visible":true}}` {
		t.Errorf("Expected anyone to read properties, got %s", getResult.Ret)
	}
}

func TestMintBatchPropertiesOnlySetOnFirstMint(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// First batch mint with properties
	mint1 := []byte(`{"to":"hive:tibfox","ids":["bp-first"],"amounts":[5],"maxSupplies":[100],"properties":[{"original":true}],"data":""}`)
	CallContract(t, ct, "mintBatch", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Second batch mint with different properties for same token - should be ignored
	mint2 := []byte(`{"to":"hive:tibfox","ids":["bp-first"],"amounts":[5],"properties":[{"original":false}],"data":""}`)
	CallContract(t, ct, "mintBatch", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Properties should still be from first mint
	getPayload := []byte(`{"id":"bp-first"}`)
	getResult, _, _ := CallContract(t, ct, "getProperties", getPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if getResult.Ret != `{"properties":{"original":true}}` {
		t.Errorf("Expected first mint properties to persist, got %s", getResult.Ret)
	}
}
