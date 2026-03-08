package contract_test

import (
	"testing"
)

// ===================================
// URI Tests
// ===================================

func TestURI(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Default URI uses baseURI + id
	uriPayload := []byte(`{"id":"123"}`)
	result := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"uri":"https://api.magi.network/metadata/123"}` {
		t.Errorf("Expected URI with baseUri + id, got %s", result.Ret)
	}
}

func TestSetURI(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Set custom URI for token
	setPayload := []byte(`{"id":"1","uri":"https://custom.example.com/token1.json"}`)
	CallContract(t, ct, "setURI", setPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Get custom URI
	uriPayload := []byte(`{"id":"1"}`)
	result := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"uri":"https://custom.example.com/token1.json"}` {
		t.Errorf("Expected custom URI, got %s", result.Ret)
	}
}

func TestSetURIFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	setPayload := []byte(`{"id":"1","uri":"https://custom.example.com/token1.json"}`)
	CallContract(t, ct, "setURI", setPayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestSetURIOverridesBaseURI(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Before setting custom URI, it uses baseURI
	uriPayload := []byte(`{"id":"test-token"}`)
	result := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"uri":"https://api.magi.network/metadata/test-token"}` {
		t.Errorf("Expected baseUri + id, got %s", result.Ret)
	}

	// Set custom URI
	setPayload := []byte(`{"id":"test-token","uri":"https://custom.example.com/special.json"}`)
	CallContract(t, ct, "setURI", setPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Now it should return custom URI
	result2 := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"uri":"https://custom.example.com/special.json"}` {
		t.Errorf("Expected custom URI, got %s", result2.Ret)
	}
}

func TestSetURICanBeUpdated(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Set first custom URI
	setPayload1 := []byte(`{"id":"1","uri":"https://first.example.com/token.json"}`)
	CallContract(t, ct, "setURI", setPayload1, nil, ownerAddress, true, uint(150_000_000), "")

	// Update to second custom URI
	setPayload2 := []byte(`{"id":"1","uri":"https://second.example.com/updated.json"}`)
	CallContract(t, ct, "setURI", setPayload2, nil, ownerAddress, true, uint(150_000_000), "")

	// Should return updated URI
	uriPayload := []byte(`{"id":"1"}`)
	result := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"uri":"https://second.example.com/updated.json"}` {
		t.Errorf("Expected updated URI, got %s", result.Ret)
	}
}

func TestURIWithEmptyBaseUri(t *testing.T) {
	ct := SetupContractTest()
	// Init with empty baseUri
	initPayload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":""}`)
	CallContract(t, ct, "init", initPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Token without custom URI should return empty string
	uriPayload := []byte(`{"id":"123"}`)
	result := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"uri":""}` {
		t.Errorf("Expected empty URI, got %s", result.Ret)
	}
}

func TestURIDifferentTokensDifferentURIs(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Set custom URI for token 1
	setPayload1 := []byte(`{"id":"1","uri":"https://custom.example.com/token1.json"}`)
	CallContract(t, ct, "setURI", setPayload1, nil, ownerAddress, true, uint(150_000_000), "")

	// Set custom URI for token 2
	setPayload2 := []byte(`{"id":"2","uri":"https://custom.example.com/token2.json"}`)
	CallContract(t, ct, "setURI", setPayload2, nil, ownerAddress, true, uint(150_000_000), "")

	// Token 1 should have its URI
	uriPayload1 := []byte(`{"id":"1"}`)
	result1 := CallContract(t, ct, "uri", uriPayload1, nil, ownerAddress, true, uint(150_000_000), "")
	if result1.Ret != `{"uri":"https://custom.example.com/token1.json"}` {
		t.Errorf("Expected token1 URI, got %s", result1.Ret)
	}

	// Token 2 should have its URI
	uriPayload2 := []byte(`{"id":"2"}`)
	result2 := CallContract(t, ct, "uri", uriPayload2, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"uri":"https://custom.example.com/token2.json"}` {
		t.Errorf("Expected token2 URI, got %s", result2.Ret)
	}

	// Token 3 without custom URI should use baseUri
	uriPayload3 := []byte(`{"id":"3"}`)
	result3 := CallContract(t, ct, "uri", uriPayload3, nil, ownerAddress, true, uint(150_000_000), "")
	if result3.Ret != `{"uri":"https://api.magi.network/metadata/3"}` {
		t.Errorf("Expected baseUri + id for token3, got %s", result3.Ret)
	}
}

// ===================================
// SetBaseURI Tests
// ===================================

func TestSetBaseURI(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Initial URI uses original baseUri
	uriPayload := []byte(`{"id":"123"}`)
	result := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"uri":"https://api.magi.network/metadata/123"}` {
		t.Errorf("Expected original baseUri, got %s", result.Ret)
	}

	// Update baseUri
	setPayload := []byte(`{"baseUri":"https://newapi.example.com/nft/"}`)
	CallContract(t, ct, "setBaseURI", setPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// URI should now use new baseUri
	result2 := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"uri":"https://newapi.example.com/nft/123"}` {
		t.Errorf("Expected new baseUri, got %s", result2.Ret)
	}
}

func TestSetBaseURIFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	setPayload := []byte(`{"baseUri":"https://newapi.example.com/nft/"}`)
	CallContract(t, ct, "setBaseURI", setPayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestSetBaseURIToEmpty(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Set baseUri to empty
	setPayload := []byte(`{"baseUri":""}`)
	CallContract(t, ct, "setBaseURI", setPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// URI should now be empty for tokens without custom URI
	uriPayload := []byte(`{"id":"123"}`)
	result := CallContract(t, ct, "uri", uriPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"uri":""}` {
		t.Errorf("Expected empty URI, got %s", result.Ret)
	}
}

func TestSetBaseURIDoesNotAffectCustomURIs(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Set custom URI for token 1
	setCustom := []byte(`{"id":"1","uri":"https://custom.example.com/token1.json"}`)
	CallContract(t, ct, "setURI", setCustom, nil, ownerAddress, true, uint(150_000_000), "")

	// Update baseUri
	setBase := []byte(`{"baseUri":"https://newapi.example.com/nft/"}`)
	CallContract(t, ct, "setBaseURI", setBase, nil, ownerAddress, true, uint(150_000_000), "")

	// Token 1 should still have custom URI (not affected by baseUri change)
	uri1 := []byte(`{"id":"1"}`)
	result1 := CallContract(t, ct, "uri", uri1, nil, ownerAddress, true, uint(150_000_000), "")
	if result1.Ret != `{"uri":"https://custom.example.com/token1.json"}` {
		t.Errorf("Expected custom URI unchanged, got %s", result1.Ret)
	}

	// Token 2 should use new baseUri
	uri2 := []byte(`{"id":"2"}`)
	result2 := CallContract(t, ct, "uri", uri2, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"uri":"https://newapi.example.com/nft/2"}` {
		t.Errorf("Expected new baseUri, got %s", result2.Ret)
	}
}

func TestSetBaseURIReflectedInGetInfo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Update baseUri
	setPayload := []byte(`{"baseUri":"https://updated.example.com/metadata/"}`)
	CallContract(t, ct, "setBaseURI", setPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// getInfo should reflect updated baseUri
	result := CallContract(t, ct, "getInfo", nil, nil, ownerAddress, true, uint(150_000_000), "")
	expected := `{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://updated.example.com/metadata/","trackMinted":false}`
	if result.Ret != expected {
		t.Errorf("Expected updated baseUri in getInfo, got %s", result.Ret)
	}
}
