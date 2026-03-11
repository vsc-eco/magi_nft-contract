package contract_test

import (
	"testing"
)

// ===================================
// Collection Metadata Tests
// ===================================

// --- Positive Tests ---

func TestInitWithMetadata(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","metadata":{"description":"A cool collection","image":"https://example.com/logo.png"}}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(300_000_000), "")

	// Verify metadata was stored
	result := CallContract(t, ct, "getCollectionMetadata", nil, nil, ownerAddress, true, uint(300_000_000), "")
	expected := `{"metadata":{"description":"A cool collection","image":"https://example.com/logo.png"}}`
	if result.Ret != expected {
		t.Errorf("Expected %s, got %s", expected, result.Ret)
	}
}

func TestInitWithoutMetadata(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(300_000_000), "")

	// Metadata should be null when not provided
	result := CallContract(t, ct, "getCollectionMetadata", nil, nil, ownerAddress, true, uint(300_000_000), "")
	expected := `{"metadata":null}`
	if result.Ret != expected {
		t.Errorf("Expected %s, got %s", expected, result.Ret)
	}
}

func TestSetCollectionMetadata(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(300_000_000), "")

	// Set metadata after init
	payload := []byte(`{"metadata":{"description":"Updated collection","image":"https://example.com/new-logo.png","banner":"https://example.com/banner.png"}}`)
	CallContract(t, ct, "setCollectionMetadata", payload, nil, ownerAddress, true, uint(300_000_000), "")

	// Verify
	result := CallContract(t, ct, "getCollectionMetadata", nil, nil, ownerAddress, true, uint(300_000_000), "")
	expected := `{"metadata":{"description":"Updated collection","image":"https://example.com/new-logo.png","banner":"https://example.com/banner.png"}}`
	if result.Ret != expected {
		t.Errorf("Expected %s, got %s", expected, result.Ret)
	}
}

func TestSetCollectionMetadataOverridesInit(t *testing.T) {
	ct := SetupContractTest()
	// Init with metadata
	initPayload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"","metadata":{"description":"Original"}}`)
	CallContract(t, ct, "init", initPayload, nil, ownerAddress, true, uint(300_000_000), "")

	// Override with setCollectionMetadata
	setPayload := []byte(`{"metadata":{"description":"Overridden","image":"https://example.com/img.png"}}`)
	CallContract(t, ct, "setCollectionMetadata", setPayload, nil, ownerAddress, true, uint(300_000_000), "")

	// Should have the overridden value
	result := CallContract(t, ct, "getCollectionMetadata", nil, nil, ownerAddress, true, uint(300_000_000), "")
	expected := `{"metadata":{"description":"Overridden","image":"https://example.com/img.png"}}`
	if result.Ret != expected {
		t.Errorf("Expected %s, got %s", expected, result.Ret)
	}
}

func TestSetCollectionMetadataWithNestedJSON(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(300_000_000), "")

	// Complex nested JSON
	payload := []byte(`{"metadata":{"description":"Complex","social":{"twitter":"@magi","discord":"magi.gg"},"tags":["nft","art","gaming"]}}`)
	CallContract(t, ct, "setCollectionMetadata", payload, nil, ownerAddress, true, uint(300_000_000), "")

	result := CallContract(t, ct, "getCollectionMetadata", nil, nil, ownerAddress, true, uint(300_000_000), "")
	expected := `{"metadata":{"description":"Complex","social":{"twitter":"@magi","discord":"magi.gg"},"tags":["nft","art","gaming"]}}`
	if result.Ret != expected {
		t.Errorf("Expected %s, got %s", expected, result.Ret)
	}
}

// --- Negative Tests ---

func TestSetCollectionMetadataFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(300_000_000), "")

	payload := []byte(`{"metadata":{"description":"Hacked"}}`)
	CallContract(t, ct, "setCollectionMetadata", payload, nil, "hive:attacker", false, uint(300_000_000), "")
}

func TestSetCollectionMetadataFailsIfNotInit(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"metadata":{"description":"Too early"}}`)
	CallContract(t, ct, "setCollectionMetadata", payload, nil, ownerAddress, false, uint(300_000_000), "")
}

func TestSetCollectionMetadataFailsWithEmptyMetadata(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(300_000_000), "")

	// Empty payload
	CallContract(t, ct, "setCollectionMetadata", nil, nil, ownerAddress, false, uint(300_000_000), "")
}

func TestSetCollectionMetadataFailsWithEmptyMetadataField(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(300_000_000), "")

	// metadata field present but empty object — this passes as valid JSON raw so it should store it
	// But a completely missing metadata key should fail
	payload := []byte(`{}`)
	CallContract(t, ct, "setCollectionMetadata", payload, nil, ownerAddress, false, uint(300_000_000), "")
}

func TestGetCollectionMetadataFailsIfNotInit(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "getCollectionMetadata", nil, nil, ownerAddress, false, uint(300_000_000), "")
}
