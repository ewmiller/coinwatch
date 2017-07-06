package main

import (
	"fmt" 
	"net/http" 
	"io/ioutil"
	"encoding/json"
	"os"
	"time"
	"strconv"
)


// exchange must be 'gemini', 'cex', or 'coinbase' 
// TODO: add more exchanges
func getExchangeData(exchange string) (map[string]interface{}, map[string]interface{}) {
	btcUrl := ""
	ethUrl := ""

	switch exchange {
	case "gemini":
		btcUrl = "https://api.gemini.com/v1/pubticker/btcusd"
		ethUrl = "https://api.gemini.com/v1/pubticker/ethusd"
	case "coinbase":
		btcUrl = "https://api.coinbase.com/v2/prices/BTC-USD/spot"
		ethUrl = "https://api.coinbase.com/v2/prices/ETH-USD/spot"
	case "cex":
		btcUrl = "https://cex.io/api/ticker/BTC/USD"
		ethUrl = "https://cex.io/api/ticker/ETH/USD"
	default:
		fmt.Println("No exchange provided")
		btcUrl = "https://api.gemini.com/v1/pubticker/btcusd"
		ethUrl = "https://api.gemini.com/v1/pubticker/ethusd"
	}

	btc, err := http.Get(btcUrl)
	eth, err := http.Get(ethUrl)
	if(err != nil){
		fmt.Println("Error requesting coin prices")
		os.Exit(1)
	}

	btcReader, err := ioutil.ReadAll(btc.Body)
	ethReader, err := ioutil.ReadAll(eth.Body)
	btc.Body.Close()
	eth.Body.Close()

	if(exchange != "coinbase"){
		var btcMap map[string]interface{}
		var ethMap map[string]interface{}

		if err := json.Unmarshal(btcReader, &btcMap); err != nil {
			panic(err)
		}

		if err := json.Unmarshal(ethReader, &ethMap); err != nil {
			panic(err)
		}

		return btcMap, ethMap

	// coinbase api is extra nested
	} else {
		
		var btcMap map[string]map[string]interface{}
		var ethMap map[string]map[string]interface{}

		if err := json.Unmarshal(btcReader, &btcMap); err != nil {
			// btcMap = make(map[string]map[string]interface{})
			// btcMap["data"] = make(map[string]interface{})
			// btcMap["data"]["amount"] = "(error getting price)"
			panic(err)
		}

		if err := json.Unmarshal(ethReader, &ethMap); err != nil {
			panic(err)
		}

		return btcMap["data"], ethMap["data"]
	}

}

func printPrices() {
	
	btcMapGemini, ethMapGemini := getExchangeData("gemini")
	btcMapCex, ethMapCex := getExchangeData("cex")
	
	fmt.Printf("BTC:\n     Gemini   $%s\n     CEX      $%s\n", btcMapGemini["last"], btcMapCex["last"])
	fmt.Printf("ETH:\n     Gemini   $%s\n     CEX      $%s\n", ethMapGemini["last"], ethMapCex["last"])
	fmt.Println("--------------------")
}

func main() {

	seconds := 60


	if(len(os.Args[:1]) > 0) {
		i := 1
		for( i <= len(os.Args[1:])) {
			switch os.Args[i] { 
			case "--help":
				fmt.Println("Usage: coinwatch <options> \n where <options> are one or more of the following: " + 
						"\n--help : print this menu" + 
						"\n--interval <seconds> : set the time interval to check on (default is 60 seconds)")
				os.Exit(0)
			case "--interval":
				if(len(os.Args[i:]) == 1){
					fmt.Println("No time interval provided. Use --help for more info.")
					os.Exit(2)
				} else {
					sec, err := strconv.Atoi(os.Args[i+1])
					if(err != nil){
						fmt.Println("Invalid time argument provided.")
						os.Exit(2)
					} else {
						seconds = sec
					}
					i++
				}
			default:
				i++
			}
		}
	}
	fmt.Println("\n ----- Welcome to CoinWatch! ----- \n Getting prices...")
	fmt.Println("--------------------")

	for {
		printPrices()
		time.Sleep(time.Duration(seconds) * time.Second)
	}

}