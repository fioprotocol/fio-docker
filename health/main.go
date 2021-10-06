package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)


const (
	chainIdMainnet = `21dcae42c0182200e93f954a074011f9048a7624c6fe81d3c9541a614a88bd1c`
	chainIdTestnet = `b20901380af44ef59c5918439a1f9a41d83669020319a80574b804a5f95cbd7e`
	lagMinutes     = 2 * time.Minute
)

var (
	url string
	chainId string
	isLambda bool
)

type getInfo struct {
	ChainId string `json:"chain_id"`
	HeadBlockTime string `json:"head_block_time"`
}

func main()  {
	if e := check(); e != nil {
		log.Fatalln(e)
	}
}

func check() error {
	opts()
	resp, err := http.Get(url+"/v1/chain/get_info")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got invalid http status code %d", resp.StatusCode)
	}
	info := &getInfo{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, info)
	if err != nil {
		return err
	}
	if info.ChainId != chainId {
		return fmt.Errorf("invalid chain id when checking host, wanted %s, got %s", chainId, info.ChainId)
	}
	headTime, err := time.Parse("2006-01-02T15:04:05.999", info.HeadBlockTime)
	if err != nil {
		return err
	}
	if headTime.Before(time.Now().Add(-lagMinutes)) {
		return fmt.Errorf("headblock is more than %v minutes behind", lagMinutes)
	}
	return nil
}

func opts() () {
	flag.StringVar(&url, "u", "http://127.0.0.1:8888", "nodeos url")
	flag.StringVar(&chainId, "c", "mainnet", "chain id for node, hex string (without 0x,) or use mainnet / testnet")
	flag.Parse()
	switch chainId {
	case "mainnet", "":
		chainId = chainIdMainnet
	case "testnet":
		chainId = chainIdTestnet
	}
	return
}