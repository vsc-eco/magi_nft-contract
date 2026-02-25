package contract_test

import (
	"fmt"
	"strconv"
	"testing"
)

func InactiveTestUniqueTemplateBenchmark(t *testing.T) {
	maxGas := uint(100_000_000_000)

	const (
		groups    = 20
		perGroup  = 1_000
		batchSize = 50
	)

	batches := perGroup / batchSize

	// --------------------------------------------------
	// Benchmark 1: Per-token properties (no template)
	// --------------------------------------------------
	var totalPerToken int64
	for g := 0; g < groups; g++ {
		ct := SetupContractTest()
		CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, maxGas, "")

		var groupRC int64
		for b := 0; b < batches; b++ {
			ids := make([]string, batchSize)
			amounts := make([]uint64, batchSize)
			maxSupplies := make([]uint64, batchSize)
			properties := make([]map[string]any, batchSize)

			for i := 0; i < batchSize; i++ {
				idx := b*batchSize + i
				ids[i] = "p" + strconv.Itoa(g) + "-" + strconv.Itoa(idx)
				amounts[i] = 1
				maxSupplies[i] = 1
				properties[i] = map[string]any{
					"name":       "Group #" + strconv.Itoa(g),
					"rarity":     "common",
					"power":      42,
					"collection": "Per-Token Collection",
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
			res, _, _ := CallContract(t, ct, "mintBatch", payload, nil, ownerAddress, true, maxGas, "")
			groupRC += res.RcUsed
		}
		totalPerToken += groupRC
		fmt.Printf("Per-token Group %2d: %d RC\n", g, groupRC)
	}

	// --------------------------------------------------
	// Benchmark 2: Template properties
	// --------------------------------------------------
	var totalTemplate int64
	for g := 0; g < groups; g++ {
		ct := SetupContractTest()
		CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, maxGas, "")

		templateId := "g" + strconv.Itoa(g) + "-0"
		var groupRC int64

		for b := 0; b < batches; b++ {
			ids := make([]string, batchSize)
			amounts := make([]uint64, batchSize)
			maxSupplies := make([]uint64, batchSize)

			for i := 0; i < batchSize; i++ {
				idx := b*batchSize + i
				ids[i] = "g" + strconv.Itoa(g) + "-" + strconv.Itoa(idx)
				amounts[i] = 1
				maxSupplies[i] = 1
			}

			m := map[string]any{
				"to":                 "hive:tibfox",
				"ids":                ids,
				"amounts":            amounts,
				"maxSupplies":        maxSupplies,
				"propertiesTemplate": templateId,
				"data":               "",
			}

			// First batch: first token gets explicit properties
			if b == 0 {
				props := make([]map[string]any, 1)
				props[0] = map[string]any{
					"name":       "Group #" + strconv.Itoa(g),
					"rarity":     "common",
					"power":      42,
					"collection": "Template Collection",
				}
				m["properties"] = props
			}

			payload := ToJSONRaw(m)
			res, _, _ := CallContract(t, ct, "mintBatch", payload, nil, ownerAddress, true, maxGas, "")
			groupRC += res.RcUsed
		}

		totalTemplate += groupRC
		fmt.Printf("Template Group %2d: %d RC\n", g, groupRC)
	}

	// --------------------------------------------------
	// Print comparison
	// --------------------------------------------------
	totalTokens := groups * perGroup

	fmt.Println("\n========================================")
	fmt.Println("UNIQUE NFT BENCHMARK (20 × 1,000)")
	fmt.Println("========================================")
	fmt.Printf("%-45s %d RC\n", "Per-token properties — total", totalPerToken)
	fmt.Printf("%-45s %.2f RC\n", "Per-token properties — per token", float64(totalPerToken)/float64(totalTokens))
	fmt.Printf("%-45s %d RC\n", "Template properties — total", totalTemplate)
	fmt.Printf("%-45s %.2f RC\n", "Template properties — per token", float64(totalTemplate)/float64(totalTokens))
	fmt.Printf("%-45s %d RC\n", "Savings", totalPerToken-totalTemplate)
	if totalPerToken > 0 {
		fmt.Printf("%-45s %.1f%%\n", "Reduction", float64(totalPerToken-totalTemplate)/float64(totalPerToken)*100)
	}
	fmt.Println("========================================")
}
