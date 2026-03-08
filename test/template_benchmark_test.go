package contract_test

import (
	"fmt"
	"strconv"
	"testing"
)

func TestTemplateBenchmark(t *testing.T) {
	ct := SetupContractTest()

	maxGas := uint(10_000_000_000)

	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, maxGas, "")

	// --------------------------------------------------
	// Benchmark 1: MintBatch 100 unique NFTs WITH per-token properties (baseline)
	// --------------------------------------------------
	ct1 := SetupContractTest()
	CallContract(t, ct1, "init", DefaultInitPayload, nil, ownerAddress, true, maxGas, "")

	var totalWithProps int64
	for batch := 0; batch < 2; batch++ {
		ids := make([]string, 50)
		amounts := make([]uint64, 50)
		maxSupplies := make([]uint64, 50)
		properties := make([]map[string]any, 50)

		for i := 0; i < 50; i++ {
			idx := batch*50 + i
			ids[i] = "nft-" + strconv.Itoa(idx)
			amounts[i] = 1
			maxSupplies[i] = 1
			properties[i] = map[string]any{
				"name":   "NFT #" + strconv.Itoa(idx),
				"rarity": "common",
				"power":  idx % 100,
			}
		}

		payload := ToJSONRaw(map[string]any{
			"to":          "hive:tibfox",
			"ids":         ids,
			"amounts":     amounts,
			"maxSupplies": maxSupplies,
			"properties":  properties,
			"data":        "",
		})
		res := CallContract(t, ct1, "mintBatch", payload, nil, ownerAddress, true, maxGas, "")
		totalWithProps += res.RcUsed
	}

	// --------------------------------------------------
	// Benchmark 2: MintBatch 100 unique NFTs WITH template (only 1 props write)
	// --------------------------------------------------
	ct2 := SetupContractTest()
	CallContract(t, ct2, "init", DefaultInitPayload, nil, ownerAddress, true, maxGas, "")

	var totalWithTemplate int64
	for batch := 0; batch < 2; batch++ {
		ids := make([]string, 50)
		amounts := make([]uint64, 50)
		maxSupplies := make([]uint64, 50)

		for i := 0; i < 50; i++ {
			idx := batch*50 + i
			ids[i] = "tmpl-" + strconv.Itoa(idx)
			amounts[i] = 1
			maxSupplies[i] = 1
		}

		m := map[string]any{
			"to":                 "hive:tibfox",
			"ids":                ids,
			"amounts":            amounts,
			"maxSupplies":        maxSupplies,
			"propertiesTemplate": "tmpl-0",
			"data":               "",
		}

		// First batch: first token gets explicit properties
		if batch == 0 {
			props := make([]map[string]any, 1)
			props[0] = map[string]any{
				"name":   "Template Card",
				"rarity": "common",
				"power":  42,
			}
			m["properties"] = props
		}

		payload := ToJSONRaw(m)
		res := CallContract(t, ct2, "mintBatch", payload, nil, ownerAddress, true, maxGas, "")
		totalWithTemplate += res.RcUsed
	}

	// --------------------------------------------------
	// Benchmark 3: MintBatch 100 unique NFTs WITHOUT any properties
	// --------------------------------------------------
	ct3 := SetupContractTest()
	CallContract(t, ct3, "init", DefaultInitPayload, nil, ownerAddress, true, maxGas, "")

	var totalNoProps int64
	for batch := 0; batch < 2; batch++ {
		ids := make([]string, 50)
		amounts := make([]uint64, 50)
		maxSupplies := make([]uint64, 50)

		for i := 0; i < 50; i++ {
			idx := batch*50 + i
			ids[i] = "bare-" + strconv.Itoa(idx)
			amounts[i] = 1
			maxSupplies[i] = 1
		}

		payload := ToJSONRaw(map[string]any{
			"to":          "hive:tibfox",
			"ids":         ids,
			"amounts":     amounts,
			"maxSupplies": maxSupplies,
			"data":        "",
		})
		res := CallContract(t, ct3, "mintBatch", payload, nil, ownerAddress, true, maxGas, "")
		totalNoProps += res.RcUsed
	}

	// --------------------------------------------------
	// Print comparison
	// --------------------------------------------------
	fmt.Println("\n========================================")
	fmt.Println("TEMPLATE PROPERTIES BENCHMARK")
	fmt.Println("========================================")
	fmt.Printf("%-50s %d RC\n", "100 unique NFTs — per-token properties", totalWithProps)
	fmt.Printf("%-50s %d RC\n", "100 unique NFTs — template properties", totalWithTemplate)
	fmt.Printf("%-50s %d RC\n", "100 unique NFTs — no properties", totalNoProps)
	fmt.Printf("%-50s %d RC\n", "Savings (per-token vs template)", totalWithProps-totalWithTemplate)
	if totalWithProps > 0 {
		fmt.Printf("%-50s %.1f%%\n", "Reduction", float64(totalWithProps-totalWithTemplate)/float64(totalWithProps)*100)
	}
	fmt.Println("========================================")
}
