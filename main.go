package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type BlockRepsonse struct {
}

func main() {
	baseURL := "https://litepool.space/api/v1/blocks/"
	for blkNum := 0; blkNum < 1; blkNum++ {
		queryURL := fmt.Sprint(baseURL, blkNum)

		response, err := http.Get(queryURL)
		if err != nil {
			log.Fatal(err)
		}

		blkData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(blkData))
	}
}
