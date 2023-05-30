package main

import (
    /*
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
    */
    log "github.com/sirupsen/logrus"
    "poop.fi/poop-server/internal/mlog"
    "poop.fi/poop-server/internal/config"
    "poop.fi/poop-server/internal/database"
    "poop.fi/poop-server/internal/events"
)

func main() {
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("config error: %s", err)
    }

    err = mlog.InitLog(cfg)
    if err != nil {
        log.Fatalf("log error: %s", err)
    }

    err = database.InitDatabase(cfg)
    if err != nil {
        log.Fatalf("database error: %s", err)
    }

    err = events.Run(cfg)
    if err != nil {
        log.Fatalf("events error: %s", err)
    }

    /*
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
                transferEvent, err := abi.Unpack("Transfer", vLog.Data)
                if err != nil {
                    log.Fatal(err)
                }
                fmt.Printf("Transfer Event: %v\n", transferEvent);
        }
    }
    */

}
