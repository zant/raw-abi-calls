package main

import (
  "bytes"
  "encoding/hex"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "time"

  "github.com/ethereum/go-ethereum/common/hexutil"
  "github.com/ethereum/go-ethereum/crypto"
)

const (
  j           = "application/json"
  flashbotURL = "https://mainnet.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161"
  method      = "eth_call"
  p           = "POST"
)

func pad32(data []byte) []byte {
  if len(data) > 32 {
    log.Fatal("Nope")
  }
  var arr [32]uint

  padLength := len(arr) - len(data)
  pad := []byte{}
  for i := 0; i < padLength; i++ {
    pad = append(pad, 0)
  }
  padded := append(pad, data...)
  return padded
}

func main() {
  mevHTTPClient := &http.Client{
    Timeout: time.Second * 3,
  }

  function := hexutil.Encode(crypto.Keccak256([]byte("balanceOf(address)"))[:4])
  account, _ := hex.DecodeString("0x0b6a7872221067fe3047de725015bd0826314074")
  padded := pad32(account)
  to := "0x6b175474e89094c44da98b954eedeac495271d0f"

  data := function + hex.EncodeToString(padded)
  one := map[string]interface{}{
    "data": data,
    "to":   to,
  }

  params := map[string]interface{}{
    "jsonrpc": "2.0",
    "id":      1,
    "method":  method,
    "params": []interface{}{
      one,
      "latest",
    },
  }

  payload, _ := json.Marshal(params)
  fmt.Println(string(payload))

  req, _ := http.NewRequest(p, flashbotURL, bytes.NewBuffer(payload))
  req.Header.Add("content-type", j)
  req.Header.Add("Accept", j)

  resp, _ := mevHTTPClient.Do(req)
  res, _ := ioutil.ReadAll(resp.Body)

  fmt.Println(string(res))
}
