package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
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

func writeCSV(end int, s *[]BlockData) {
	// create CSV
	csvFile, err := os.Create(fmt.Sprint(end, ".csv"))
	if err != nil {
		log.Fatal(err)
	}
	// csvFile.Close()

	// write to CSV
	w := csv.NewWriter(csvFile)
	defer w.Flush()

	blkData := *s

	for _, record := range blkData {
		d := []string{
			fmt.Sprint(record.Height),
			record.CoinbaseData,
			record.CoinbaseAddress,
		}
		w.Write(d)
	}
}

func queryBlocks(start int, end int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	baseURL := "https://litepool.space/api/v1/blocks/"
	var s []BlockData

	for blkNum := start; blkNum <= end; blkNum++ {
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

		r := BlockData{
			blkResponse[0].Height,
			blkResponse[0].Extras.CoinbaseData,
			blkResponse[0].Extras.CoinbaseAddress,
		}

		s = append(s, r)
	}

	// dump slice to csv file
	writeCSV(end, &s)
	ch <- start
}

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	// spawn goroutines
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go queryBlocks(i*100000, i*100000+99999, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}
}
