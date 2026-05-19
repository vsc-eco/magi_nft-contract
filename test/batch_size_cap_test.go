package contract_test

import (
	"fmt"
	"strings"
	"testing"
)

// ===================================
// NFT-10: Batch Size Cap (maxBatchSize = 256)
// ===================================
//
// A/B test for the unbounded-batch finding.
//
// Pre-fix (RED): an oversized mintBatch (maxBatchSize+1 = 257 ids) was
//   accepted, so the call below would succeed and the
//   ExpectSuccess=false assertion fails the test.
// Post-fix (GREEN): the same oversized batch aborts with
//   "Batch size exceeds maximum", so the call fails as expected and
//   the test passes. A normal small batch still succeeds (no regression).

const testMaxBatchSize = 256

// buildBatchPayload builds a mintBatch payload with n distinct token ids,
// each amount=1 and maxSupply=1.
func buildBatchPayload(n int) []byte {
	ids := make([]string, n)
	amounts := make([]string, n)
	maxSupplies := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = fmt.Sprintf("%q", fmt.Sprintf("tok-%d", i))
		amounts[i] = "1"
		maxSupplies[i] = "1"
	}
	return []byte(fmt.Sprintf(
		`{"to":"hive:tibfox","ids":[%s],"amounts":[%s],"maxSupplies":[%s],"data":""}`,
		strings.Join(ids, ","),
		strings.Join(amounts, ","),
		strings.Join(maxSupplies, ","),
	))
}

// TestMintBatchRejectsOversizedBatch is the RED/GREEN assertion: an
// oversized batch (maxBatchSize+1 entries) must be rejected.
func TestMintBatchRejectsOversizedBatch(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	oversized := buildBatchPayload(testMaxBatchSize + 1) // 257 ids
	// Pre-fix: the batch is processed and the call succeeds, so the
	//   ExpectSuccess=false assertion fails the test (RED).
	// Post-fix: the call aborts with "Batch size exceeds maximum", so
	//   the call fails as expected and the test passes (GREEN).
	// High gas ceiling so the meaningful RED failure is the semantic
	// "did not fail" assertion, not an incidental gas-limit breach.
	CallContract(t, ct, "mintBatch", oversized, nil, ownerAddress, false, uint(8_000_000_000), "")
}

// TestMintBatchAcceptsNormalBatch guards against regression: a small,
// in-bounds batch must still succeed after the cap is added.
func TestMintBatchAcceptsNormalBatch(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	small := buildBatchPayload(3)
	CallContract(t, ct, "mintBatch", small, nil, ownerAddress, true, uint(200_000_000), "")
}

// TestMintBatchAcceptsBoundaryBatch verifies the boundary is inclusive:
// exactly maxBatchSize entries is still accepted.
func TestMintBatchAcceptsBoundaryBatch(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(150_000_000), "")

	boundary := buildBatchPayload(testMaxBatchSize) // exactly 256
	// A full 256-entry batch legitimately consumes a lot of gas (~6.9B);
	// the high ceiling here only asserts functional success, not a gas bound.
	CallContract(t, ct, "mintBatch", boundary, nil, ownerAddress, true, uint(8_000_000_000), "")
}
