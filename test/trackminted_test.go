package contract_test

import (
	"strconv"
	"testing"

	"vsc-node/lib/test_utils"
)

// ===================================
// TrackMinted Tests
// ===================================

func TestTrackMintedEnabled(t *testing.T) {
	ct := SetupContractTest()
	// Init with trackMinted enabled
	initPayload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":true}`)
	CallContract(t, ct, "init", initPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify trackMinted is enabled in getInfo
	result, _, _ := CallContract(t, ct, "getInfo", nil, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":true}` {
		t.Errorf("Expected trackMinted true in getInfo, got %s", result.Ret)
	}

	const tokenId = "track-test"

	// 1. Mint 5 of 10 max supply
	mint1 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":5,"maxSupply":10,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=5, totalMinted=5
	verifySupplyWithMinted(t, ct, tokenId, 5, 10, 5)

	// 2. Burn 3
	burn := []byte(`{"from":"` + ownerAddress + `","id":"` + tokenId + `","amount":3}`)
	CallContract(t, ct, "burn", burn, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=2, totalMinted=5 (unchanged)
	verifySupplyWithMinted(t, ct, tokenId, 2, 10, 5)

	// 3. Mint 5 more - should succeed (totalMinted would be 10)
	mint2 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":5,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=7, totalMinted=10
	verifySupplyWithMinted(t, ct, tokenId, 7, 10, 10)

	// 4. Try to mint 1 more - should FAIL (totalMinted already at max)
	mint3 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":1,"data":""}`)
	CallContract(t, ct, "mint", mint3, nil, ownerAddress, false, uint(150_000_000), "")

	// Even though totalSupply is only 7 (3 were burned), we can't mint more
	// because totalMinted is at max
}

func TestTrackMintedDisabled(t *testing.T) {
	ct := SetupContractTest()
	// Init without trackMinted (default behavior)
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify trackMinted is disabled in getInfo
	result, _, _ := CallContract(t, ct, "getInfo", nil, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":false}` {
		t.Errorf("Expected trackMinted false in getInfo, got %s", result.Ret)
	}

	const tokenId = "no-track-test"

	// 1. Mint 10 of 10 max supply (fill up)
	mint1 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":10,"maxSupply":10,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=10, totalMinted=0 (not tracked)
	verifySupplyWithMinted(t, ct, tokenId, 10, 10, 0)

	// 2. Burn 5
	burn := []byte(`{"from":"` + ownerAddress + `","id":"` + tokenId + `","amount":5}`)
	CallContract(t, ct, "burn", burn, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=5
	verifySupply(t, ct, tokenId, 5, 10)

	// 3. Mint 5 more - should succeed (re-minting burned slots)
	mint2 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":5,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=10 again
	verifySupply(t, ct, tokenId, 10, 10)
}

func TestTotalMintedQuery(t *testing.T) {
	ct := SetupContractTest()
	initPayload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":true}`)
	CallContract(t, ct, "init", initPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Query totalMinted for non-existent token
	payload := []byte(`{"id":"nonexistent"}`)
	result, _, _ := CallContract(t, ct, "totalMinted", payload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"totalMinted":0}` {
		t.Errorf("Expected totalMinted 0 for non-existent token, got %s", result.Ret)
	}

	// Mint some tokens
	mintPayload := []byte(`{"to":"` + ownerAddress + `","id":"test-token","amount":25,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Query totalMinted
	payload2 := []byte(`{"id":"test-token"}`)
	result2, _, _ := CallContract(t, ct, "totalMinted", payload2, nil, ownerAddress, true, uint(150_000_000), "")
	if result2.Ret != `{"totalMinted":25}` {
		t.Errorf("Expected totalMinted 25, got %s", result2.Ret)
	}
}

func TestTrackMintedWithBatchMint(t *testing.T) {
	ct := SetupContractTest()
	initPayload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":true}`)
	CallContract(t, ct, "init", initPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Batch mint
	mintBatch := []byte(`{"to":"` + ownerAddress + `","ids":["a","b","c"],"amounts":[10,20,30],"maxSupplies":[50,50,50],"data":""}`)
	CallContract(t, ct, "mintBatch", mintBatch, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify totalMinted for each
	verifyTotalMinted(t, ct, "a", 10)
	verifyTotalMinted(t, ct, "b", 20)
	verifyTotalMinted(t, ct, "c", 30)

	// Burn some
	burnBatch := []byte(`{"from":"` + ownerAddress + `","ids":["a","b"],"amounts":[5,10]}`)
	CallContract(t, ct, "burnBatch", burnBatch, nil, ownerAddress, true, uint(150_000_000), "")

	// totalMinted should remain the same (burning doesn't decrease it)
	verifyTotalMinted(t, ct, "a", 10)
	verifyTotalMinted(t, ct, "b", 20)

	// totalSupply should decrease
	verifySupply(t, ct, "a", 5, 50)
	verifySupply(t, ct, "b", 10, 50)
}

func TestTrackMintedUniqueNFTCannotRemint(t *testing.T) {
	ct := SetupContractTest()
	initPayload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":true}`)
	CallContract(t, ct, "init", initPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint unique NFT
	mint := []byte(`{"to":"` + ownerAddress + `","id":"unique","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mint, nil, ownerAddress, true, uint(150_000_000), "")

	// Burn it
	burn := []byte(`{"from":"` + ownerAddress + `","id":"unique","amount":1}`)
	CallContract(t, ct, "burn", burn, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify state
	verifySupply(t, ct, "unique", 0, 1)
	verifyTotalMinted(t, ct, "unique", 1)

	// Try to re-mint - should FAIL (totalMinted already at max)
	remint := []byte(`{"to":"` + ownerAddress + `","id":"unique","amount":1,"data":""}`)
	CallContract(t, ct, "mint", remint, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestTrackMintedLifecycle(t *testing.T) {
	ct := SetupContractTest()
	initPayload := []byte(`{"name":"Magi NFT","symbol":"MNFT","baseUri":"https://api.magi.network/metadata/","trackMinted":true}`)
	CallContract(t, ct, "init", initPayload, nil, ownerAddress, true, uint(150_000_000), "")

	const alice = "hive:alice"
	const tokenId = "permanent"

	// 1. Mint 50 of 100
	mint1 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":50,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")
	verifySupplyWithMinted(t, ct, tokenId, 50, 100, 50)

	// 2. Transfer 20 to alice
	transfer := []byte(`{"from":"` + ownerAddress + `","to":"` + alice + `","id":"` + tokenId + `","amount":20,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transfer, nil, ownerAddress, true, uint(150_000_000), "")

	// totalMinted shouldn't change on transfer
	verifyTotalMinted(t, ct, tokenId, 50)

	// 3. Alice burns 10
	burn := []byte(`{"from":"` + alice + `","id":"` + tokenId + `","amount":10}`)
	CallContract(t, ct, "burn", burn, nil, alice, true, uint(150_000_000), "")

	// totalMinted still 50, totalSupply now 40
	verifySupplyWithMinted(t, ct, tokenId, 40, 100, 50)

	// 4. Mint 50 more (totalMinted will be 100)
	mint2 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":50,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")
	verifySupplyWithMinted(t, ct, tokenId, 90, 100, 100)

	// 5. Try to mint 1 more - should FAIL
	mint3 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":1,"data":""}`)
	CallContract(t, ct, "mint", mint3, nil, ownerAddress, false, uint(150_000_000), "")

	// 6. Even if we burn everything, we still can't mint more
	burnAll := []byte(`{"from":"` + ownerAddress + `","id":"` + tokenId + `","amount":80}`)
	CallContract(t, ct, "burn", burnAll, nil, ownerAddress, true, uint(150_000_000), "")

	burnAlice := []byte(`{"from":"` + alice + `","id":"` + tokenId + `","amount":10}`)
	CallContract(t, ct, "burn", burnAlice, nil, alice, true, uint(150_000_000), "")

	// totalSupply is now 0, but totalMinted is still 100
	verifySupplyWithMinted(t, ct, tokenId, 0, 100, 100)

	// Still can't mint
	mint4 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":1,"data":""}`)
	CallContract(t, ct, "mint", mint4, nil, ownerAddress, false, uint(150_000_000), "")
}

// ===================================
// Helper for trackMinted tests
// ===================================

// Helper to verify totalMinted
func verifyTotalMinted(t *testing.T, ct *test_utils.ContractTest, tokenId string, expected uint64) {
	t.Helper()

	payload := []byte(`{"id":"` + tokenId + `"}`)
	result, _, _ := CallContract(t, ct, "totalMinted", payload, nil, ownerAddress, true, uint(150_000_000), "")
	expectedStr := `{"totalMinted":` + strconv.FormatUint(expected, 10) + `}`
	if result.Ret != expectedStr {
		t.Errorf("Expected totalMinted %d, got %s", expected, result.Ret)
	}
}

// Helper to verify totalSupply, maxSupply, and totalMinted together
func verifySupplyWithMinted(t *testing.T, ct *test_utils.ContractTest, tokenId string, expectedTotal, expectedMax, expectedMinted uint64) {
	t.Helper()
	verifySupply(t, ct, tokenId, expectedTotal, expectedMax)
	verifyTotalMinted(t, ct, tokenId, expectedMinted)
}
