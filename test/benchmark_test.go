package contract_test

import (
	"fmt"
	"strconv"
	"testing"
)

// ===================================
// Benchmark: Real-World Scenario
// ===================================
//
// This test simulates a realistic NFT workflow and reports RC consumption:
// 1. Init the contract
// 2. Mint 100 unique NFTs (1/1) with template properties (in batches of 50)
// 3. Mint 1 editioned NFT with 10,000 editions and properties
// 4. Transfer some unique NFTs
// 5. Transfer some editions
// 6. Burn some unique NFTs
// 7. Burn some editions

func TestBenchmarkScenario(t *testing.T) {
	ct := SetupContractTest()

	type rcEntry struct {
		Step string
		RC   int64
	}
	var rcLog []rcEntry

	maxGas := uint(10_000_000_000) // very high limit so nothing is constrained

	// --------------------------------------------------
	// Step 1: Init
	// --------------------------------------------------
	result := CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"init", result.RcUsed})

	// --------------------------------------------------
	// Step 1a: mint 10 unique NFTs via regular mint (one at a time)
	// --------------------------------------------------
	var totalMintSingleRC int64
	for i := 0; i < 10; i++ {
		id := "single-" + strconv.Itoa(i)
		mintPayload := []byte(`{"to":"hive:tibfox","id":"` + id + `","amount":1,"maxSupply":1,"data":""}`)
		mRes := CallContract(t, ct, "mint", mintPayload, nil, ownerAddress, true, maxGas, "")
		totalMintSingleRC += mRes.RcUsed
	}
	rcLog = append(rcLog, rcEntry{"mint 10 unique NFTs (single) — total", totalMintSingleRC})
	rcLog = append(rcLog, rcEntry{"mint 10 unique NFTs (single) — avg per mint", totalMintSingleRC / 10})

	// --------------------------------------------------
	// Step 1b: mintSeries — 10 unique NFTs
	// --------------------------------------------------
	seriesPayload := []byte(`{"to":"hive:tibfox","idPrefix":"series-","startNumber":1,"count":10,"amount":1,"maxSupply":1}`)
	seriesRes := CallContract(t, ct, "mintSeries", seriesPayload, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"mintSeries 10 unique NFTs", seriesRes.RcUsed})

	seriesPayload50 := []byte(`{"to":"hive:tibfox","idPrefix":"series50-","startNumber":1,"count":50,"amount":1,"maxSupply":1}`)
	seriesRes50 := CallContract(t, ct, "mintSeries", seriesPayload50, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"mintSeries 50 unique NFTs", seriesRes50.RcUsed})

	// --------------------------------------------------
	// Step 1c: mintSeries with template properties — 50 unique NFTs
	// First ID (tpl-1) gets properties, rest inherit via template event
	// --------------------------------------------------
	seriesTemplatePayload := []byte(`{"to":"hive:tibfox","idPrefix":"tpl-","startNumber":1,"count":50,"amount":1,"maxSupply":1,"properties":{"name":"Game Card","rarity":"common","power":42},"propertiesTemplate":"tpl-1"}`)
	seriesTemplateRes := CallContract(t, ct, "mintSeries", seriesTemplatePayload, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"mintSeries 50 unique NFTs with template properties", seriesTemplateRes.RcUsed})

	// --------------------------------------------------
	// Step 2: Mint 100 unique NFTs with template properties (batches of 50)
	// nft-0 gets explicit properties and serves as the template for all others
	// --------------------------------------------------
	batchSize := 50
	totalUnique := 100
	var totalMintUniqueRC int64

	for batch := 0; batch < totalUnique/batchSize; batch++ {
		ids := make([]string, batchSize)
		amounts := make([]uint64, batchSize)
		maxSupplies := make([]uint64, batchSize)

		for i := 0; i < batchSize; i++ {
			idx := batch*batchSize + i
			ids[i] = "nft-" + strconv.Itoa(idx)
			amounts[i] = 1
			maxSupplies[i] = 1
		}

		m := map[string]any{
			"to":                 "hive:tibfox",
			"ids":                ids,
			"amounts":            amounts,
			"maxSupplies":        maxSupplies,
			"propertiesTemplate": "nft-0",
			"data":               "",
		}

		// First batch: first token gets explicit properties
		if batch == 0 {
			props := make([]map[string]any, 1)
			props[0] = map[string]any{
				"name":   "Game Card",
				"rarity": "common",
				"power":  42,
			}
			m["properties"] = props
		}

		payload := ToJSONRaw(m)
		res := CallContract(t, ct, "mintBatch", payload, nil, ownerAddress, true, maxGas, "")
		rcLog = append(rcLog, rcEntry{fmt.Sprintf("mintBatch 50 unique NFTs with template — batch %d", batch+1), res.RcUsed})
		totalMintUniqueRC += res.RcUsed
	}
	rcLog = append(rcLog, rcEntry{"mintBatch 100 unique NFTs (2x50) with template — total", totalMintUniqueRC})
	rcLog = append(rcLog, rcEntry{"mintBatch 100 unique NFTs (2x50) with template — avg per batch", totalMintUniqueRC / int64(totalUnique/batchSize)})

	// --------------------------------------------------
	// Step 3: Mint 1 editioned NFT with 10,000 editions and properties
	// --------------------------------------------------
	editionPayload := []byte(`{"to":"hive:tibfox","id":"edition-mega","amount":10000,"maxSupply":10000,"properties":{"collection":"Mega Edition","type":"collectible","series":1},"data":""}`)
	res := CallContract(t, ct, "mint", editionPayload, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"mint 10,000 editions with properties", res.RcUsed})

	// --------------------------------------------------
	// Step 4: Transfer some unique NFTs (10 single transfers + 1 batch of 50)
	// Note: nft-0 is a template and cannot be transferred, so start from nft-1
	// --------------------------------------------------
	var totalTransferSingleRC int64
	for i := 1; i <= 10; i++ {
		payload := ToJSONRaw(map[string]any{
			"from":   "hive:tibfox",
			"to":     "hive:collector",
			"id":     "nft-" + strconv.Itoa(i),
			"amount": 1,
			"data":   "",
		})
		res := CallContract(t, ct, "safeTransferFrom", payload, nil, ownerAddress, true, maxGas, "")
		totalTransferSingleRC += res.RcUsed
	}
	rcLog = append(rcLog, rcEntry{"safeTransferFrom 10 unique NFTs (owner) — total", totalTransferSingleRC})
	rcLog = append(rcLog, rcEntry{"safeTransferFrom 10 unique NFTs (owner) — avg per transfer", totalTransferSingleRC / 10})

	// --------------------------------------------------
	// Step 4b: Approve + transfer via per-token allowance (ERC-6909)
	// Transfer nft-61 to nft-70 via approve (they are still owned by tibfox)
	// --------------------------------------------------
	var totalApproveRC int64
	for i := 61; i <= 70; i++ {
		id := "nft-" + strconv.Itoa(i)
		approvePayload := ToJSONRaw(map[string]any{
			"spender": "hive:marketplace",
			"id":      id,
			"amount":  1,
		})
		aRes := CallContract(t, ct, "approve", approvePayload, nil, ownerAddress, true, maxGas, "")
		totalApproveRC += aRes.RcUsed
	}
	rcLog = append(rcLog, rcEntry{"approve 10 unique NFTs — total", totalApproveRC})
	rcLog = append(rcLog, rcEntry{"approve 10 unique NFTs — avg per call", totalApproveRC / 10})

	var totalTransferAllowanceRC int64
	for i := 61; i <= 70; i++ {
		id := "nft-" + strconv.Itoa(i)
		payload := ToJSONRaw(map[string]any{
			"from":   "hive:tibfox",
			"to":     "hive:collector",
			"id":     id,
			"amount": 1,
			"data":   "",
		})
		tRes := CallContract(t, ct, "safeTransferFrom", payload, nil, "hive:marketplace", true, maxGas, "")
		totalTransferAllowanceRC += tRes.RcUsed
	}
	rcLog = append(rcLog, rcEntry{"safeTransferFrom 10 unique NFTs (allowance) — total", totalTransferAllowanceRC})
	rcLog = append(rcLog, rcEntry{"safeTransferFrom 10 unique NFTs (allowance) — avg per transfer", totalTransferAllowanceRC / 10})

	// Batch transfer 50 unique NFTs (nft-11 to nft-60, since 1-10 already transferred)
	batchTransferIds := make([]string, 50)
	batchTransferAmounts := make([]uint64, 50)
	for i := 0; i < 50; i++ {
		batchTransferIds[i] = "nft-" + strconv.Itoa(11+i)
		batchTransferAmounts[i] = 1
	}
	batchTransferPayload := ToJSONRaw(map[string]any{
		"from":    "hive:tibfox",
		"to":      "hive:collector",
		"ids":     batchTransferIds,
		"amounts": batchTransferAmounts,
		"data":    "",
	})
	res = CallContract(t, ct, "safeBatchTransferFrom", batchTransferPayload, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"safeBatchTransferFrom 50 unique NFTs", res.RcUsed})

	// --------------------------------------------------
	// Step 5: Transfer some editions
	// --------------------------------------------------
	editionTransferPayload := ToJSONRaw(map[string]any{
		"from":   "hive:tibfox",
		"to":     "hive:collector",
		"id":     "edition-mega",
		"amount": 1000,
		"data":   "",
	})
	res = CallContract(t, ct, "safeTransferFrom", editionTransferPayload, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"safeTransferFrom 1,000 editions", res.RcUsed})

	editionTransferPayload2 := ToJSONRaw(map[string]any{
		"from":   "hive:tibfox",
		"to":     "hive:buyer",
		"id":     "edition-mega",
		"amount": 500,
		"data":   "",
	})
	res = CallContract(t, ct, "safeTransferFrom", editionTransferPayload2, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"safeTransferFrom 500 editions", res.RcUsed})

	// --------------------------------------------------
	// Step 6: Burn some unique NFTs (5 single burns + 1 batch of 20)
	// --------------------------------------------------
	// nft-71 to nft-75 (still owned by tibfox)
	var totalBurnSingleRC int64
	for i := 0; i < 5; i++ {
		payload := ToJSONRaw(map[string]any{
			"from":   "hive:tibfox",
			"id":     "nft-" + strconv.Itoa(71+i),
			"amount": 1,
		})
		res := CallContract(t, ct, "burn", payload, nil, ownerAddress, true, maxGas, "")
		totalBurnSingleRC += res.RcUsed
	}
	rcLog = append(rcLog, rcEntry{"burn 5 unique NFTs — total", totalBurnSingleRC})
	rcLog = append(rcLog, rcEntry{"burn 5 unique NFTs — avg per burn", totalBurnSingleRC / 5})

	// Batch burn 20 unique NFTs (nft-76 to nft-95)
	batchBurnIds := make([]string, 20)
	batchBurnAmounts := make([]uint64, 20)
	for i := 0; i < 20; i++ {
		batchBurnIds[i] = "nft-" + strconv.Itoa(76+i)
		batchBurnAmounts[i] = 1
	}
	batchBurnPayload := ToJSONRaw(map[string]any{
		"from":    "hive:tibfox",
		"ids":     batchBurnIds,
		"amounts": batchBurnAmounts,
	})
	res = CallContract(t, ct, "burnBatch", batchBurnPayload, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"burnBatch 20 unique NFTs — total", res.RcUsed})
	rcLog = append(rcLog, rcEntry{"burnBatch 20 unique NFTs — avg per burn", res.RcUsed / 20})

	// --------------------------------------------------
	// Step 7: Burn some editions
	// --------------------------------------------------
	burnEditionPayload := ToJSONRaw(map[string]any{
		"from":   "hive:tibfox",
		"id":     "edition-mega",
		"amount": 100,
	})
	res = CallContract(t, ct, "burn", burnEditionPayload, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"burn 100 editions", res.RcUsed})

	burnEditionPayload2 := ToJSONRaw(map[string]any{
		"from":   "hive:tibfox",
		"id":     "edition-mega",
		"amount": 1000,
	})
	res = CallContract(t, ct, "burn", burnEditionPayload2, nil, ownerAddress, true, maxGas, "")
	rcLog = append(rcLog, rcEntry{"burn 1000 editions", res.RcUsed})

	// --------------------------------------------------
	// Print RC Summary
	// --------------------------------------------------
	fmt.Println("\n========================================")
	fmt.Println("RC CONSUMPTION SUMMARY")
	fmt.Println("========================================")
	for _, entry := range rcLog {
		fmt.Printf("%-65s %d RC\n", entry.Step, entry.RC)
	}
	fmt.Println("========================================")
}
