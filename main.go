package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type BlockResponse []struct {
	Id     string `json:"id"`
	Height int    `json:"height"`
	Extras struct {
		Pool struct {
			Name string `json:"name"`
		}
		CoinbaseAddress string `json:"CoinbaseAddress"`
		CoinbaseData    string `json:"CoinbaseSignatureAscii"`
	}
}

type BlockData struct {
	Height          int
	Id              string
	CoinbaseData    string
	CoinbaseAddress string
}

func containsElement(slice *[]BlockData, address string) (exists bool) {
	result := false
	blockData := *slice

	for _, element := range blockData {
		if element.CoinbaseAddress == address {
			result = true
			break
		}
	}

	return result
}

func checkBlock(blkResponse *BlockResponse, s *[]BlockData) (exists bool) {
	blk := *blkResponse
	if blk[0].Extras.Pool.Name != "Unknown" {
		return true
	}

	// some older blocks don't have identified coinbaseAddresses
	if blk[0].Extras.CoinbaseAddress == "null" {
		return false
	}

	// check if unidentified block is held in slice
	// check for duplication by Coinbase Address reuse
	return containsElement(s, blk[0].Extras.CoinbaseAddress)
}

func main() {
	baseURL := "https://litepool.space/api/v1/blocks/"

	// create CSV
	// csvFile, err := os.Create("100k.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// csvFile.Close()

	var s []BlockData

	// loop queries blocks
	for blkNum := 500000; blkNum <= 501000; blkNum++ {
		queryURL := fmt.Sprint(baseURL, blkNum)

		response, err := http.Get(queryURL)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var blkResponse BlockResponse
		json.Unmarshal(body, &blkResponse)

		// check block if unknown miner exists
		blockExists := checkBlock(&blkResponse, &s)
		if blockExists {
			continue
		}

		r := BlockData{blkResponse[0].Height,
			blkResponse[0].Id,
			blkResponse[0].Extras.CoinbaseData,
			blkResponse[0].Extras.CoinbaseAddress,
		}

		s = append(s, r)
	}

}
