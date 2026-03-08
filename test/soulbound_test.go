package contract_test

import (
	"testing"
)

// ===================================
// Soulbound Token Tests
// ===================================

func TestMintSoulboundToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a soulbound token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"soulbound-1","amount":1,"maxSupply":1,"soulbound":true,"data":""}`)
	result := CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected success, got %s", result.Ret)
	}

	// Verify token was minted
	balPayload := []byte(`{"account":"hive:tibfox","id":"soulbound-1"}`)
	balResult := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":1}` {
		t.Errorf("Expected balance 1, got %s", balResult.Ret)
	}
}

func TestSoulboundTokenCannotBeTransferred(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a soulbound token to a recipient
	mintPayload := []byte(`{"to":"hive:recipient","id":"soulbound-2","amount":1,"maxSupply":1,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Recipient tries to transfer - should fail (soulbound)
	transferPayload := []byte(`{"from":"hive:recipient","to":"hive:other","id":"soulbound-2","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, "hive:recipient", false, uint(150_000_000), "")
}

func TestOwnerCanTransferSoulboundToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint soulbound token to owner (for distribution)
	mintPayload := []byte(`{"to":"hive:tibfox","id":"soulbound-dist","amount":5,"maxSupply":10,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Owner can transfer (distribute) soulbound tokens
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","id":"soulbound-dist","amount":2,"data":""}`)
	result := CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected owner to transfer soulbound token, got %s", result.Ret)
	}

	// Verify recipient got the tokens
	balPayload := []byte(`{"account":"hive:recipient","id":"soulbound-dist"}`)
	balResult := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":2}` {
		t.Errorf("Expected balance 2, got %s", balResult.Ret)
	}

	// Recipient cannot transfer them further
	transfer2 := []byte(`{"from":"hive:recipient","to":"hive:other","id":"soulbound-dist","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transfer2, nil, "hive:recipient", false, uint(150_000_000), "")
}

func TestSoulboundTokenCanBeBurned(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a soulbound token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"soulbound-burn","amount":1,"maxSupply":1,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Burn should succeed
	burnPayload := []byte(`{"from":"hive:tibfox","id":"soulbound-burn","amount":1}`)
	result := CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected burn success, got %s", result.Ret)
	}

	// Verify balance is 0
	balPayload := []byte(`{"account":"hive:tibfox","id":"soulbound-burn"}`)
	balResult := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":0}` {
		t.Errorf("Expected balance 0 after burn, got %s", balResult.Ret)
	}
}

func TestBatchTransferFailsIfAnySoulbound(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint regular token to recipient
	mint1 := []byte(`{"to":"hive:holder","id":"regular-1","amount":5,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint soulbound token to recipient
	mint2 := []byte(`{"to":"hive:holder","id":"soulbound-batch","amount":5,"maxSupply":100,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Recipient tries to batch transfer both - should fail because of soulbound
	transferPayload := []byte(`{"from":"hive:holder","to":"hive:other","ids":["regular-1","soulbound-batch"],"amounts":[1,1],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, "hive:holder", false, uint(150_000_000), "")
}

func TestOwnerCanBatchTransferSoulbound(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint tokens to owner for distribution
	mint1 := []byte(`{"to":"hive:tibfox","id":"sb-batch-1","amount":5,"maxSupply":100,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	mint2 := []byte(`{"to":"hive:tibfox","id":"sb-batch-2","amount":5,"maxSupply":100,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Owner can batch transfer soulbound tokens (for distribution)
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:recipient","ids":["sb-batch-1","sb-batch-2"],"amounts":[2,3],"data":""}`)
	result := CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected owner batch transfer to succeed, got %s", result.Ret)
	}

	// Verify balances
	bal1 := []byte(`{"account":"hive:recipient","id":"sb-batch-1"}`)
	bal1Result := CallContract(t, ct, "balanceOf", bal1, nil, ownerAddress, true, uint(150_000_000), "")
	if bal1Result.Ret != `{"balance":2}` {
		t.Errorf("Expected balance 2, got %s", bal1Result.Ret)
	}

	// Recipient cannot batch transfer them
	transfer2 := []byte(`{"from":"hive:recipient","to":"hive:other","ids":["sb-batch-1","sb-batch-2"],"amounts":[1,1],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transfer2, nil, "hive:recipient", false, uint(150_000_000), "")
}

func TestIsSoulboundQuery(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a soulbound token
	mintSoulbound := []byte(`{"to":"hive:tibfox","id":"sb-query","amount":1,"maxSupply":1,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mintSoulbound, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a regular token
	mintRegular := []byte(`{"to":"hive:tibfox","id":"regular-query","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintRegular, nil, ownerAddress, true, uint(150_000_000), "")

	// Query soulbound status
	sbPayload := []byte(`{"id":"sb-query"}`)
	sbResult := CallContract(t, ct, "isSoulbound", sbPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if sbResult.Ret != `{"soulbound":true}` {
		t.Errorf("Expected soulbound true, got %s", sbResult.Ret)
	}

	// Query regular token
	regPayload := []byte(`{"id":"regular-query"}`)
	regResult := CallContract(t, ct, "isSoulbound", regPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if regResult.Ret != `{"soulbound":false}` {
		t.Errorf("Expected soulbound false, got %s", regResult.Ret)
	}

	// Query non-existent token
	nonPayload := []byte(`{"id":"nonexistent"}`)
	nonResult := CallContract(t, ct, "isSoulbound", nonPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if nonResult.Ret != `{"soulbound":false}` {
		t.Errorf("Expected soulbound false for non-existent, got %s", nonResult.Ret)
	}
}

func TestSoulboundOnlySetOnFirstMint(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// First mint as regular token
	mint1 := []byte(`{"to":"hive:tibfox","id":"first-mint","amount":5,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Second mint tries to set soulbound - should be ignored
	mint2 := []byte(`{"to":"hive:tibfox","id":"first-mint","amount":5,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Token should NOT be soulbound (was set on first mint without soulbound)
	sbPayload := []byte(`{"id":"first-mint"}`)
	sbResult := CallContract(t, ct, "isSoulbound", sbPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if sbResult.Ret != `{"soulbound":false}` {
		t.Errorf("Expected soulbound false (set on first mint), got %s", sbResult.Ret)
	}

	// Transfer should work since it's not soulbound
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:other","id":"first-mint","amount":1,"data":""}`)
	result := CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected transfer success for non-soulbound token, got %s", result.Ret)
	}
}

func TestNonSoulboundTokenCanBeTransferred(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint a regular (non-soulbound) token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"transferable","amount":10,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Transfer should succeed
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:other","id":"transferable","amount":5,"data":""}`)
	result := CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected transfer success, got %s", result.Ret)
	}

	// Verify balances
	bal1 := []byte(`{"account":"hive:tibfox","id":"transferable"}`)
	bal1Result := CallContract(t, ct, "balanceOf", bal1, nil, ownerAddress, true, uint(150_000_000), "")
	if bal1Result.Ret != `{"balance":5}` {
		t.Errorf("Expected balance 5, got %s", bal1Result.Ret)
	}

	bal2 := []byte(`{"account":"hive:other","id":"transferable"}`)
	bal2Result := CallContract(t, ct, "balanceOf", bal2, nil, ownerAddress, true, uint(150_000_000), "")
	if bal2Result.Ret != `{"balance":5}` {
		t.Errorf("Expected balance 5, got %s", bal2Result.Ret)
	}
}

func TestMintBatchWithSoulbound(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint batch with mixed soulbound flags
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["batch-sb","batch-regular"],"amounts":[1,1],"maxSupplies":[1,100],"soulbound":[true,false],"data":""}`)
	result := CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected success, got %s", result.Ret)
	}

	// Verify soulbound status
	sb1 := []byte(`{"id":"batch-sb"}`)
	sb1Result := CallContract(t, ct, "isSoulbound", sb1, nil, ownerAddress, true, uint(150_000_000), "")
	if sb1Result.Ret != `{"soulbound":true}` {
		t.Errorf("Expected batch-sb soulbound true, got %s", sb1Result.Ret)
	}

	sb2 := []byte(`{"id":"batch-regular"}`)
	sb2Result := CallContract(t, ct, "isSoulbound", sb2, nil, ownerAddress, true, uint(150_000_000), "")
	if sb2Result.Ret != `{"soulbound":false}` {
		t.Errorf("Expected batch-regular soulbound false, got %s", sb2Result.Ret)
	}

	// Owner can transfer soulbound (for distribution)
	transferSb := []byte(`{"from":"hive:tibfox","to":"hive:other","id":"batch-sb","amount":1,"data":""}`)
	transferSbResult := CallContract(t, ct, "safeTransferFrom", transferSb, nil, ownerAddress, true, uint(150_000_000), "")
	if transferSbResult.Ret != `{"success":true}` {
		t.Errorf("Expected owner to transfer soulbound, got %s", transferSbResult.Ret)
	}

	// Transfer regular - should succeed
	transferReg := []byte(`{"from":"hive:tibfox","to":"hive:other","id":"batch-regular","amount":1,"data":""}`)
	transferRegResult := CallContract(t, ct, "safeTransferFrom", transferReg, nil, ownerAddress, true, uint(150_000_000), "")
	if transferRegResult.Ret != `{"success":true}` {
		t.Errorf("Expected success for regular token transfer, got %s", transferRegResult.Ret)
	}

	// Recipient cannot transfer soulbound further
	transferSb2 := []byte(`{"from":"hive:other","to":"hive:buyer","id":"batch-sb","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferSb2, nil, "hive:other", false, uint(150_000_000), "")
}

func TestSoulboundTransferFailsViaOperator(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint soulbound token to user
	mintPayload := []byte(`{"to":"hive:user1","id":"sb-operator","amount":1,"maxSupply":1,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// User approves operator
	approvePayload := []byte(`{"operator":"hive:marketplace","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approvePayload, nil, "hive:user1", true, uint(150_000_000), "")

	// Operator tries to transfer - should fail because soulbound
	transferPayload := []byte(`{"from":"hive:user1","to":"hive:buyer","id":"sb-operator","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, "hive:marketplace", false, uint(150_000_000), "")
}

func TestSoulboundBatchBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint batch of soulbound tokens
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["sb-burn1","sb-burn2"],"amounts":[5,10],"maxSupplies":[10,20],"soulbound":[true,true],"data":""}`)
	CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Batch burn should succeed
	burnPayload := []byte(`{"from":"hive:tibfox","ids":["sb-burn1","sb-burn2"],"amounts":[3,5]}`)
	result := CallContract(t, ct, "burnBatch", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected burn batch success, got %s", result.Ret)
	}

	// Verify balances after burn
	bal1 := []byte(`{"account":"hive:tibfox","id":"sb-burn1"}`)
	bal1Result := CallContract(t, ct, "balanceOf", bal1, nil, ownerAddress, true, uint(150_000_000), "")
	if bal1Result.Ret != `{"balance":2}` {
		t.Errorf("Expected balance 2, got %s", bal1Result.Ret)
	}

	bal2 := []byte(`{"account":"hive:tibfox","id":"sb-burn2"}`)
	bal2Result := CallContract(t, ct, "balanceOf", bal2, nil, ownerAddress, true, uint(150_000_000), "")
	if bal2Result.Ret != `{"balance":5}` {
		t.Errorf("Expected balance 5, got %s", bal2Result.Ret)
	}
}

func TestSoulboundWithExplicitFalse(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint with explicit soulbound: false
	mintPayload := []byte(`{"to":"hive:tibfox","id":"explicit-false","amount":1,"maxSupply":10,"soulbound":false,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Should NOT be soulbound
	sbPayload := []byte(`{"id":"explicit-false"}`)
	sbResult := CallContract(t, ct, "isSoulbound", sbPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if sbResult.Ret != `{"soulbound":false}` {
		t.Errorf("Expected soulbound false, got %s", sbResult.Ret)
	}

	// Transfer should work
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:other","id":"explicit-false","amount":1,"data":""}`)
	result := CallContract(t, ct, "safeTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected transfer success, got %s", result.Ret)
	}
}

func TestIsSoulboundRemainsTrueAfterBurn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint soulbound token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"sb-after-burn","amount":5,"maxSupply":10,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Burn all tokens
	burnPayload := []byte(`{"from":"hive:tibfox","id":"sb-after-burn","amount":5}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// isSoulbound should still return true (property of the token type)
	sbPayload := []byte(`{"id":"sb-after-burn"}`)
	sbResult := CallContract(t, ct, "isSoulbound", sbPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if sbResult.Ret != `{"soulbound":true}` {
		t.Errorf("Expected soulbound true after burn, got %s", sbResult.Ret)
	}

	// Verify balance is 0
	balPayload := []byte(`{"account":"hive:tibfox","id":"sb-after-burn"}`)
	balResult := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if balResult.Ret != `{"balance":0}` {
		t.Errorf("Expected balance 0, got %s", balResult.Ret)
	}
}

func TestSoulboundHolderCannotTransferOwn(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint soulbound token to user
	mintPayload := []byte(`{"to":"hive:holder","id":"sb-holder","amount":1,"maxSupply":1,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Holder tries to transfer their own token - should fail
	transferPayload := []byte(`{"from":"hive:holder","to":"hive:other","id":"sb-holder","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferPayload, nil, "hive:holder", false, uint(150_000_000), "")
}

func TestSoulboundMultipleMintsSameToken(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// First mint as soulbound
	mint1 := []byte(`{"to":"hive:tibfox","id":"sb-multi","amount":5,"maxSupply":100,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Second mint to different address
	mint2 := []byte(`{"to":"hive:other","id":"sb-multi","amount":10,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Both balances should be correct
	bal1 := []byte(`{"account":"hive:tibfox","id":"sb-multi"}`)
	bal1Result := CallContract(t, ct, "balanceOf", bal1, nil, ownerAddress, true, uint(150_000_000), "")
	if bal1Result.Ret != `{"balance":5}` {
		t.Errorf("Expected balance 5, got %s", bal1Result.Ret)
	}

	bal2 := []byte(`{"account":"hive:other","id":"sb-multi"}`)
	bal2Result := CallContract(t, ct, "balanceOf", bal2, nil, ownerAddress, true, uint(150_000_000), "")
	if bal2Result.Ret != `{"balance":10}` {
		t.Errorf("Expected balance 10, got %s", bal2Result.Ret)
	}

	// Token should still be soulbound
	sbPayload := []byte(`{"id":"sb-multi"}`)
	sbResult := CallContract(t, ct, "isSoulbound", sbPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if sbResult.Ret != `{"soulbound":true}` {
		t.Errorf("Expected soulbound true, got %s", sbResult.Ret)
	}

	// Owner CAN transfer soulbound tokens (for distribution)
	transfer1 := []byte(`{"from":"hive:tibfox","to":"hive:buyer","id":"sb-multi","amount":1,"data":""}`)
	transferResult := CallContract(t, ct, "safeTransferFrom", transfer1, nil, ownerAddress, true, uint(150_000_000), "")
	if transferResult.Ret != `{"success":true}` {
		t.Errorf("Expected owner transfer to succeed, got %s", transferResult.Ret)
	}

	// Recipient cannot transfer their soulbound tokens
	transfer2 := []byte(`{"from":"hive:other","to":"hive:buyer","id":"sb-multi","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transfer2, nil, "hive:other", false, uint(150_000_000), "")
}

func TestSoulboundWithTrackMinted(t *testing.T) {
	ct := SetupContractTest()

	// Init with trackMinted enabled
	initPayload := []byte(`{"name":"Soulbound NFT","symbol":"SBNFT","baseUri":"https://api.test.com/","trackMinted":true}`)
	CallContract(t, ct, "init", initPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint soulbound token
	mintPayload := []byte(`{"to":"hive:tibfox","id":"sb-tracked","amount":5,"maxSupply":10,"soulbound":true,"data":""}`)
	CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Burn some
	burnPayload := []byte(`{"from":"hive:tibfox","id":"sb-tracked","amount":3}`)
	CallContract(t, ct, "burn", burnPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Try to mint more (trackMinted means burned tokens count against supply)
	mint2 := []byte(`{"to":"hive:other","id":"sb-tracked","amount":6,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, false, uint(150_000_000), "") // Should fail - 5 minted + 6 = 11 > 10

	// Can still mint up to remaining
	mint3 := []byte(`{"to":"hive:other","id":"sb-tracked","amount":5,"data":""}`)
	CallContract(t, ct, "mint", mint3, nil, ownerAddress, true, uint(150_000_000), "") // Should succeed - 5 + 5 = 10

	// Token should still be soulbound
	sbPayload := []byte(`{"id":"sb-tracked"}`)
	sbResult := CallContract(t, ct, "isSoulbound", sbPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if sbResult.Ret != `{"soulbound":true}` {
		t.Errorf("Expected soulbound true with trackMinted, got %s", sbResult.Ret)
	}
}

func TestSoulboundMintBatchWithoutSoulboundArray(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint batch without soulbound array (should default to false)
	mintPayload := []byte(`{"to":"hive:tibfox","ids":["no-sb-1","no-sb-2"],"amounts":[1,1],"maxSupplies":[1,1],"data":""}`)
	result := CallContract(t, ct, "mintBatch", mintPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected success, got %s", result.Ret)
	}

	// Both should NOT be soulbound
	sb1 := []byte(`{"id":"no-sb-1"}`)
	sb1Result := CallContract(t, ct, "isSoulbound", sb1, nil, ownerAddress, true, uint(150_000_000), "")
	if sb1Result.Ret != `{"soulbound":false}` {
		t.Errorf("Expected soulbound false, got %s", sb1Result.Ret)
	}

	sb2 := []byte(`{"id":"no-sb-2"}`)
	sb2Result := CallContract(t, ct, "isSoulbound", sb2, nil, ownerAddress, true, uint(150_000_000), "")
	if sb2Result.Ret != `{"soulbound":false}` {
		t.Errorf("Expected soulbound false, got %s", sb2Result.Ret)
	}
}

func TestSoulboundBatchTransferOnlyRegularTokens(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	// Mint two regular tokens
	mint1 := []byte(`{"to":"hive:tibfox","id":"reg-batch-1","amount":5,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	mint2 := []byte(`{"to":"hive:tibfox","id":"reg-batch-2","amount":5,"maxSupply":100,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Batch transfer should succeed (no soulbound tokens)
	transferPayload := []byte(`{"from":"hive:tibfox","to":"hive:other","ids":["reg-batch-1","reg-batch-2"],"amounts":[2,3],"data":""}`)
	result := CallContract(t, ct, "safeBatchTransferFrom", transferPayload, nil, ownerAddress, true, uint(150_000_000), "")
	if result.Ret != `{"success":true}` {
		t.Errorf("Expected batch transfer success, got %s", result.Ret)
	}

	// Verify balances
	bal1 := []byte(`{"account":"hive:other","id":"reg-batch-1"}`)
	bal1Result := CallContract(t, ct, "balanceOf", bal1, nil, ownerAddress, true, uint(150_000_000), "")
	if bal1Result.Ret != `{"balance":2}` {
		t.Errorf("Expected balance 2, got %s", bal1Result.Ret)
	}

	bal2 := []byte(`{"account":"hive:other","id":"reg-batch-2"}`)
	bal2Result := CallContract(t, ct, "balanceOf", bal2, nil, ownerAddress, true, uint(150_000_000), "")
	if bal2Result.Ret != `{"balance":3}` {
		t.Errorf("Expected balance 3, got %s", bal2Result.Ret)
	}
}
