package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	var byteArray []byte
	var magicNumber = make([]byte, 4)
	magicNumber[0] = 0xfb
	magicNumber[1] = 0xc0
	magicNumber[2] = 0xb6
	magicNumber[3] = 0xdb

	for blk := 1; blk <= 256; blk++ {
		blkHashUrl := fmt.Sprint("https://litepool.space/api/block-height/", blk)
		blkHashRes, err := http.Get(blkHashUrl)
		if err != nil {
			log.Fatal(err)
		}
		defer blkHashRes.Body.Close()

		blockHash, err := io.ReadAll(blkHashRes.Body)
		if err != nil {
			log.Fatal(err)
		}

		rawBlockURL := fmt.Sprint("https://litepool.space/api/block/", string(blockHash), "/raw")
		rawBlockRes, err := http.Get(rawBlockURL)
		if err != nil {
			log.Fatal(err)
		}
		defer rawBlockRes.Body.Close()

		rawBlock, err := io.ReadAll(rawBlockRes.Body)
		if err != nil {
			log.Fatal(err)
		}
		rawBlockLength := make([]byte, 4)

		binary.LittleEndian.PutUint32(rawBlockLength, uint32(len(rawBlock)))

		byteArray = append(byteArray, magicNumber...)
		byteArray = append(byteArray, rawBlockLength...)
		byteArray = append(byteArray, rawBlock...)

	}

	fmt.Print(fmt.Sprintf("%x", byteArray))
}
