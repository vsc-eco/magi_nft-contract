package contract_test

import (
	"testing"
)

// ===================================
// Init Tests
// ===================================

func TestInitSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
}

func TestInitWithTrackMinted(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":true}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(150_000_000), "")

	result := CallContract(t, ct, "getInfo", nil, nil, ownerAddress, true, uint(150_000_000), "")
	expected := `{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":true}`
	if result.Ret != expected {
		t.Errorf("Expected %s, got %s", expected, result.Ret)
	}
}

func TestInitFailsIfAlreadyInit(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestInitFailsIfNotOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, "hive:other", false, uint(150_000_000), "")
}

func TestInitFailsWithEmptyName(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestInitFailsWithEmptySymbol(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Magi NFT","symbol":"","baseUri":"https://api.magi.network/metadata/"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestGetInfo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")
	result := CallContract(t, ct, "getInfo", nil, nil, ownerAddress, true, uint(150_000_000), "")
	expected := `{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":false}`
	if result.Ret != expected {
		t.Errorf("Expected %s, got %s", expected, result.Ret)
	}
}

func TestGetInfoFailsIfNotInit(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "getInfo", nil, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestGetOwnerFailsIfNotInit(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "getOwner", nil, nil, ownerAddress, false, uint(150_000_000), "")
}
