package contract_test

import (
	"strconv"
	"testing"

	"vsc-node/lib/test_utils"
)

// ===================================
// Lifecycle Tests
// ===================================

func TestEditionedNFTLifecycle(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	const tokenId = "edition-lifecycle"
	const alice = "hive:alice"
	const bob = "hive:bob"
	const otherContract = "hive:othercontract"

	// 1. First mint: 20 of 30 max supply to owner
	mint1 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":20,"maxSupply":30,"data":""}`)
	CallContract(t, ct, "mint", mint1, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=20, maxSupply=30, owner balance=20
	verifySupply(t, ct, tokenId, 20, 30)
	verifyBalance(t, ct, ownerAddress, tokenId, 20)

	// 2. Send 10 to alice
	transfer1 := []byte(`{"from":"` + ownerAddress + `","to":"` + alice + `","id":"` + tokenId + `","amount":10,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transfer1, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: owner=10, alice=10
	verifyBalance(t, ct, ownerAddress, tokenId, 10)
	verifyBalance(t, ct, alice, tokenId, 10)

	// 3. Mint more (5) - no maxSupply needed for existing token
	mint2 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":5,"data":""}`)
	CallContract(t, ct, "mint", mint2, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=25, owner=15
	verifySupply(t, ct, tokenId, 25, 30)
	verifyBalance(t, ct, ownerAddress, tokenId, 15)

	// 4. Send 10 to bob
	transfer2 := []byte(`{"from":"` + ownerAddress + `","to":"` + bob + `","id":"` + tokenId + `","amount":10,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transfer2, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: owner=5, bob=10
	verifyBalance(t, ct, ownerAddress, tokenId, 5)
	verifyBalance(t, ct, bob, tokenId, 10)

	// 5. Mint more (5) to fill up to 30
	mint3 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":5,"data":""}`)
	CallContract(t, ct, "mint", mint3, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify: totalSupply=30 (maxed out), owner=10
	verifySupply(t, ct, tokenId, 30, 30)
	verifyBalance(t, ct, ownerAddress, tokenId, 10)

	// 6. Bob burns 5
	burn1 := []byte(`{"from":"` + bob + `","id":"` + tokenId + `","amount":5}`)
	CallContract(t, ct, "burn", burn1, nil, bob, true, uint(150_000_000), "")

	// Verify: totalSupply=25 (burn decreases), bob=5
	verifySupply(t, ct, tokenId, 25, 30)
	verifyBalance(t, ct, bob, tokenId, 5)

	// 7. Alice approves otherContract to transfer on her behalf
	approve := []byte(`{"operator":"` + otherContract + `","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approve, nil, alice, true, uint(150_000_000), "")

	// Verify approval
	approvalCheck := []byte(`{"account":"` + alice + `","operator":"` + otherContract + `"}`)
	result := CallContract(t, ct, "isApprovedForAll", approvalCheck, nil, alice, true, uint(150_000_000), "")
	if result.Ret != `{"approved":true}` {
		t.Errorf("Expected approved true, got %s", result.Ret)
	}

	// 8. otherContract transfers 5 from alice to bob
	transfer3 := []byte(`{"from":"` + alice + `","to":"` + bob + `","id":"` + tokenId + `","amount":5,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transfer3, nil, otherContract, true, uint(150_000_000), "")

	// Verify: alice=5, bob=10
	verifyBalance(t, ct, alice, tokenId, 5)
	verifyBalance(t, ct, bob, tokenId, 10)

	// 9. Alice burns her remaining 5
	burn2 := []byte(`{"from":"` + alice + `","id":"` + tokenId + `","amount":5}`)
	CallContract(t, ct, "burn", burn2, nil, alice, true, uint(150_000_000), "")

	// Verify: totalSupply=20, alice=0
	verifySupply(t, ct, tokenId, 20, 30)
	verifyBalance(t, ct, alice, tokenId, 0)

	// 10. otherContract tries to transfer 5 more from alice to bob - should FAIL (no balance)
	transfer4 := []byte(`{"from":"` + alice + `","to":"` + bob + `","id":"` + tokenId + `","amount":5,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transfer4, nil, otherContract, false, uint(150_000_000), "")

	// 11. Mint 10 more - should succeed (totalSupply=20, max=30, so 10 more is fine)
	mint4 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":10,"data":""}`)
	CallContract(t, ct, "mint", mint4, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify final state: totalSupply=30, owner=20, bob=10, alice=0
	verifySupply(t, ct, tokenId, 30, 30)
	verifyBalance(t, ct, ownerAddress, tokenId, 20)
	verifyBalance(t, ct, bob, tokenId, 10)
	verifyBalance(t, ct, alice, tokenId, 0)

	// Bonus: Try to mint 1 more - should FAIL (would exceed max supply)
	mint5 := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenId + `","amount":1,"data":""}`)
	CallContract(t, ct, "mint", mint5, nil, ownerAddress, false, uint(150_000_000), "")
}

func TestUniqueNFTLifecycle(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	const tokenA = "unique-A"
	const tokenB = "unique-B"
	const tokenC = "unique-C"
	const alice = "hive:alice"
	const bob = "hive:bob"
	const contractX = "hive:contractx"
	const contractY = "hive:contracty"

	// 1. Mint A (unique NFT, maxSupply=1)
	mintA := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenA + `","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintA, nil, ownerAddress, true, uint(150_000_000), "")
	verifyBalance(t, ct, ownerAddress, tokenA, 1)
	verifySupply(t, ct, tokenA, 1, 1)

	// 2. Mint B (unique NFT)
	mintB := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenB + `","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintB, nil, ownerAddress, true, uint(150_000_000), "")
	verifyBalance(t, ct, ownerAddress, tokenB, 1)

	// 3. Mint C (unique NFT)
	mintC := []byte(`{"to":"` + ownerAddress + `","id":"` + tokenC + `","amount":1,"maxSupply":1,"data":""}`)
	CallContract(t, ct, "mint", mintC, nil, ownerAddress, true, uint(150_000_000), "")
	verifyBalance(t, ct, ownerAddress, tokenC, 1)

	// 4. Transfer A to bob
	transferA := []byte(`{"from":"` + ownerAddress + `","to":"` + bob + `","id":"` + tokenA + `","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferA, nil, ownerAddress, true, uint(150_000_000), "")
	verifyBalance(t, ct, ownerAddress, tokenA, 0)
	verifyBalance(t, ct, bob, tokenA, 1)

	// 5. Transfer B to alice
	transferB := []byte(`{"from":"` + ownerAddress + `","to":"` + alice + `","id":"` + tokenB + `","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferB, nil, ownerAddress, true, uint(150_000_000), "")
	verifyBalance(t, ct, ownerAddress, tokenB, 0)
	verifyBalance(t, ct, alice, tokenB, 1)

	// 6. Bob approves contractX to transfer on his behalf
	approveX := []byte(`{"operator":"` + contractX + `","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approveX, nil, bob, true, uint(150_000_000), "")

	// 7. contractX transfers A from bob to alice
	transferAtoAlice := []byte(`{"from":"` + bob + `","to":"` + alice + `","id":"` + tokenA + `","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferAtoAlice, nil, contractX, true, uint(150_000_000), "")
	verifyBalance(t, ct, bob, tokenA, 0)
	verifyBalance(t, ct, alice, tokenA, 1)

	// 8. Alice burns A
	burnA := []byte(`{"from":"` + alice + `","id":"` + tokenA + `","amount":1}`)
	CallContract(t, ct, "burn", burnA, nil, alice, true, uint(150_000_000), "")
	verifyBalance(t, ct, alice, tokenA, 0)
	verifySupply(t, ct, tokenA, 0, 1) // totalSupply=0, maxSupply still=1

	// Note: After burning, you CAN re-mint up to maxSupply again
	// For unique NFTs, burning opens up the slot for re-minting
	// If you want truly "burned forever", don't burn - just send to a dead address

	// 9. Alice transfers B to bob
	transferBtoBob := []byte(`{"from":"` + alice + `","to":"` + bob + `","id":"` + tokenB + `","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferBtoBob, nil, alice, true, uint(150_000_000), "")
	verifyBalance(t, ct, alice, tokenB, 0)
	verifyBalance(t, ct, bob, tokenB, 1)

	// 10. Bob approves contractX for all (already done in step 6, but let's be explicit)
	// Already approved from step 6

	// 11. Bob approves contractY for all
	approveY := []byte(`{"operator":"` + contractY + `","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approveY, nil, bob, true, uint(150_000_000), "")

	// Verify both approvals
	checkApprovalX := []byte(`{"account":"` + bob + `","operator":"` + contractX + `"}`)
	resultX := CallContract(t, ct, "isApprovedForAll", checkApprovalX, nil, bob, true, uint(150_000_000), "")
	if resultX.Ret != `{"approved":true}` {
		t.Errorf("Expected contractX approved, got %s", resultX.Ret)
	}

	checkApprovalY := []byte(`{"account":"` + bob + `","operator":"` + contractY + `"}`)
	resultY := CallContract(t, ct, "isApprovedForAll", checkApprovalY, nil, bob, true, uint(150_000_000), "")
	if resultY.Ret != `{"approved":true}` {
		t.Errorf("Expected contractY approved, got %s", resultY.Ret)
	}

	// 12. contractY transfers B from bob to alice
	transferBtoAliceByY := []byte(`{"from":"` + bob + `","to":"` + alice + `","id":"` + tokenB + `","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferBtoAliceByY, nil, contractY, true, uint(150_000_000), "")
	verifyBalance(t, ct, bob, tokenB, 0)
	verifyBalance(t, ct, alice, tokenB, 1)

	// 13. contractX tries to transfer B from bob to alice - should FAIL (bob no longer has B)
	transferBtoAliceByX := []byte(`{"from":"` + bob + `","to":"` + alice + `","id":"` + tokenB + `","amount":1,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transferBtoAliceByX, nil, contractX, false, uint(150_000_000), "")

	// Final state verification
	// A: burned (totalSupply=0), can never be re-minted
	// B: owned by alice
	// C: still owned by owner
	verifySupply(t, ct, tokenA, 0, 1)
	verifyBalance(t, ct, alice, tokenB, 1)
	verifyBalance(t, ct, ownerAddress, tokenC, 1)
}

func TestMultiTokenLifecycle(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	const alice = "hive:alice"
	const bob = "hive:bob"

	// Mint batch of tokens
	mintBatch := []byte(`{"to":"` + ownerAddress + `","ids":["gold","silver","bronze"],"amounts":[100,500,1000],"maxSupplies":[100,500,1000],"data":""}`)
	CallContract(t, ct, "mintBatch", mintBatch, nil, ownerAddress, true, uint(150_000_000), "")

	// Transfer batch to alice
	transferBatch := []byte(`{"from":"` + ownerAddress + `","to":"` + alice + `","ids":["gold","silver","bronze"],"amounts":[25,100,200],"data":""}`)
	CallContract(t, ct, "safeBatchTransferFrom", transferBatch, nil, ownerAddress, true, uint(150_000_000), "")

	// Verify balances
	verifyBalance(t, ct, ownerAddress, "gold", 75)
	verifyBalance(t, ct, ownerAddress, "silver", 400)
	verifyBalance(t, ct, ownerAddress, "bronze", 800)
	verifyBalance(t, ct, alice, "gold", 25)
	verifyBalance(t, ct, alice, "silver", 100)
	verifyBalance(t, ct, alice, "bronze", 200)

	// Alice approves bob
	approve := []byte(`{"operator":"` + bob + `","approved":true}`)
	CallContract(t, ct, "setApprovalForAll", approve, nil, alice, true, uint(150_000_000), "")

	// Bob transfers from alice
	transfer := []byte(`{"from":"` + alice + `","to":"` + bob + `","id":"gold","amount":10,"data":""}`)
	CallContract(t, ct, "safeTransferFrom", transfer, nil, bob, true, uint(150_000_000), "")

	// Verify
	verifyBalance(t, ct, alice, "gold", 15)
	verifyBalance(t, ct, bob, "gold", 10)

	// Burn batch from alice
	burnBatch := []byte(`{"from":"` + alice + `","ids":["silver","bronze"],"amounts":[50,100]}`)
	CallContract(t, ct, "burnBatch", burnBatch, nil, alice, true, uint(150_000_000), "")

	// Verify supplies decreased
	verifySupply(t, ct, "silver", 450, 500)
	verifySupply(t, ct, "bronze", 900, 1000)
}

// ===================================
// Helpers for lifecycle tests
// ===================================

// Helper to verify totalSupply and maxSupply
func verifySupply(t *testing.T, ct *test_utils.ContractTest, tokenId string, expectedTotal, expectedMax uint64) {
	t.Helper()

	totalPayload := []byte(`{"id":"` + tokenId + `"}`)
	result := CallContract(t, ct, "totalSupply", totalPayload, nil, ownerAddress, true, uint(150_000_000), "")
	expectedTotalStr := `{"totalSupply":` + formatUint(expectedTotal) + `}`
	if result.Ret != expectedTotalStr {
		t.Errorf("Expected totalSupply %d, got %s", expectedTotal, result.Ret)
	}

	maxPayload := []byte(`{"id":"` + tokenId + `"}`)
	result2 := CallContract(t, ct, "maxSupply", maxPayload, nil, ownerAddress, true, uint(150_000_000), "")
	expectedMaxStr := `{"maxSupply":` + formatUint(expectedMax) + `}`
	if result2.Ret != expectedMaxStr {
		t.Errorf("Expected maxSupply %d, got %s", expectedMax, result2.Ret)
	}
}

// Helper to verify balance
func verifyBalance(t *testing.T, ct *test_utils.ContractTest, account, tokenId string, expected uint64) {
	t.Helper()

	balPayload := []byte(`{"account":"` + account + `","id":"` + tokenId + `"}`)
	result := CallContract(t, ct, "balanceOf", balPayload, nil, ownerAddress, true, uint(150_000_000), "")
	expectedStr := `{"balance":` + formatUint(expected) + `}`
	if result.Ret != expectedStr {
		t.Errorf("Expected balance %d for %s, got %s", expected, account, result.Ret)
	}
}

func formatUint(n uint64) string {
	return strconv.FormatUint(n, 10)
}
