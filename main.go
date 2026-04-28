package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

const (
	// Public Mainnet WebSocket RPC
	rpcURL    = "wss://solana-rpc.publicnode.com"
	// The official Jupiter v6 Program ID
	jupV6Prog = "JUP6LkbZbjS1jKKwapdHNy74zcZ3tLUZoi5QNyVTaV4"
)

// Structs to parse the nested JSON-RPC WebSocket responses
type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type RPCResponse struct {
	Method string `json:"method"`
	Params struct {
		Result struct {
			Value struct {
				Signature string      `json:"signature"`
				Err       interface{} `json:"err"`
				Logs      []string    `json:"logs"`
			} `json:"value"`
		} `json:"result"`
	} `json:"params"`
}

func main() {
	// 1. Setup UI Colors
	title := color.New(color.FgHiCyan, color.Bold).SprintFunc()
	success := color.New(color.FgHiGreen, color.Bold).SprintFunc()
	fail := color.New(color.FgHiRed, color.Bold).SprintFunc()
	dim := color.New(color.FgHiBlack).SprintFunc()
	yellow := color.New(color.FgHiYellow).SprintFunc()
	magenta := color.New(color.FgHiMagenta).SprintFunc()

	fmt.Printf("\n%s\n", title("⚡ IONIC GOD-MODE INDEXER (Golang)"))
	fmt.Printf("%s %s\n", dim("RPC Node:"), rpcURL)
	fmt.Printf("%s %s %s\n", dim("Target:  "), yellow(jupV6Prog), dim("(Jupiter v6)"))
	fmt.Println(dim("----------------------------------------------------------------------"))

	// 2. Open WebSocket Connection
	c, _, err := websocket.DefaultDialer.Dial(rpcURL, nil)
	if err != nil {
		log.Fatal("🚨 Failed to connect to Solana RPC:", err)
	}
	defer c.Close()

	// 3. Subscribe strictly to Jupiter v6 transactions
	req := RPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "logsSubscribe",
		Params: []interface{}{
			map[string]interface{}{"mentions": []string{jupV6Prog}},
			map[string]interface{}{"commitment": "processed"},
		},
	}

	if err := c.WriteJSON(req); err != nil {
		log.Fatal("🚨 Subscription failed:", err)
	}

	// Regex to extract compute units from the raw log strings
	re := regexp.MustCompile(`consumed (\d+) of`)
	
	// Variables for our Moving Average math
	var recentCUs []int
	maxHistory := 20

	// 4. The Infinite Stream Loop
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Connection dropped:", err)
			return
		}

		var res RPCResponse
		// If the JSON doesn't match our struct (like the initial success message), ignore it
		if err := json.Unmarshal(message, &res); err != nil || res.Method != "logsNotification" {
			continue
		}

		val := res.Params.Result.Value
		
		// Format the signature to look clean (first 6...last 6)
		sig := val.Signature
		if len(sig) > 16 {
			sig = sig[:6] + "..." + sig[len(sig)-6:]
		}

		// Parse the logs to find the CU consumption
		var cuConsumed int
		for _, logLine := range val.Logs {
			matches := re.FindStringSubmatch(logLine)
			if len(matches) > 1 {
				cu, _ := strconv.Atoi(matches[1])
				cuConsumed = cu
				break
			}
		}

		// Skip if it was a weird transaction that didn't output CUs
		if cuConsumed == 0 {
			continue 
		}

		// 5. Calculate Live Moving Average
		recentCUs = append(recentCUs, cuConsumed)
		if len(recentCUs) > maxHistory {
			recentCUs = recentCUs[1:] // Pop the oldest item
		}

		sum := 0
		for _, c := range recentCUs {
			sum += c
		}
		avg := sum / len(recentCUs)

		// 6. Print the Dashboard UI
		status := success("🟢 SWAP EXEC")
		if val.Err != nil {
			status = fail("🔴 SWAP FAIL")
		}

		fmt.Printf("[Tx: %s] %s | CUs: %-6s | 20-Tx Avg: %-6s\n", 
			yellow(sig), 
			status, 
			magenta(strconv.Itoa(cuConsumed)), 
			title(strconv.Itoa(avg)))
	}
}