package contract_test

import (
	"fmt"
	"strconv"
	"testing"
)

func TestEditionsBenchmark(t *testing.T) {
	ct := SetupContractTest()

	maxGas := uint(10_000_000_000)

	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, maxGas, "")

	var totalRC int64
	for i := 0; i < 20; i++ {
		payload := []byte(`{"to":"hive:tibfox","id":"edition-` + strconv.Itoa(i) + `","amount":1000,"maxSupply":1000,"properties":{"collection":"Edition #` + strconv.Itoa(i) + `","type":"collectible"},"data":""}`)
		res, _, _ := CallContract(t, ct, "mint", payload, nil, ownerAddress, true, maxGas, "")
		totalRC += res.RcUsed
	}

	fmt.Println("\n========================================")
	fmt.Printf("Mint 20 edition NFTs (1,000 each) — total RC: %d\n", totalRC)
	fmt.Printf("Average per mint: %d RC\n", totalRC/20)
	fmt.Println("========================================")
}
