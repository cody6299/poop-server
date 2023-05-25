package main

import (
    "os"
    "log"
    "fmt"
    "strings"
    "math/big"
    "context"
    "github.com/ethereum/go-ethereum"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/crypto"
    //"github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/accounts/abi"
)

type LogTransfer struct {
    From common.Address
    To common.Address
    Tokens *big.Int
}

func main() {
    fmt.Println("Hello Sync")

    abiBytes, err := os.ReadFile("./config/abi/WETH.json")
    if err != nil {
        log.Fatal(err)
    }
    abiStr := string(abiBytes)
    //fmt.Println(abiStr)
    
    abi, err := abi.JSON(strings.NewReader(abiStr))
    if err != nil {
        log.Fatal(err)
    }
    


    client, err := ethclient.Dial("http://35.85.128.168")
    if err != nil {
        log.Fatal(err)
    }

    contractAddress := common.HexToAddress("0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c")
    query := ethereum.FilterQuery{
        FromBlock: big.NewInt(28499829),
        ToBlock: big.NewInt(28499830),
        Addresses: []common.Address{
            contractAddress,
        },
    }

    //logs := make(chan types.Log)
    logs, err := client.FilterLogs(context.Background(), query)
    if err != nil {
        log.Fatal("fat1: ", err)
    }

    logTransferSig := []byte("Transfer(address,address,uint256)")
    logTransferSigHash := crypto.Keccak256Hash(logTransferSig)

    for _, vLog := range logs {
        switch vLog.Topics[0].Hex() {
            case logTransferSigHash.Hex():
                /*
                transferEvent, err := abi.Unpack("Transfer", vLog.Data)
                if err != nil {
                    log.Fatal(err)
                }
                fmt.Printf("Transfer Event: %v\n", transferEvent);
                */
                var transferEvent LogTransfer
                err := abi.UnpackIntoInterface(&transferEvent, "Transfer", vLog.Data)
                if err != nil {
                    log.Fatal(err)
                }
                fmt.Println("Transfer Event from=%s to=%s amount=%s\n", transferEvent.From.Hex(), transferEvent.To.Hex(), transferEvent.Tokens.String())
        }
    }

}
