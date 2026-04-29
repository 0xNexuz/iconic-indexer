# ⚡ Iconic God-Mode Indexer

A high-performance, real-time WebSocket indexer built in Golang. 

This tool streams live transaction data directly from the Solana Mainnet, filters the firehose for Jupiter v6 swaps, extracts the Compute Units (CUs) burned per transaction, and calculates a live moving average on the fly. Built to demonstrate high-concurrency data pipelining and terminal-based monitoring.

## ✨ Features

- **Live Mainnet Streaming:** Maintains a persistent WebSocket connection to Solana public RPCs, catching transactions the millisecond they are processed.
- **Precision Filtering:** Uses JSON-RPC subscriptions to isolate traffic specifically for the Jupiter v6 Program ID.
- **Compute Unit Extraction:** Parses nested log arrays via regex to pull exact CU consumption for every valid swap.
- **Real-Time Analytics:** Calculates and maintains a rolling 20-transaction moving average of CU costs.
- **Matrix-Style CLI:** A highly responsive, color-coded terminal dashboard built for immediate visual feedback.

## 🚀 Quick Start

### Prerequisites
You need [Go (Golang)](https://go.dev/dl/) installed on your machine.

### Installation
1. Clone the repository:
   ```bash
   git clone [https://github.com/YOUR_USERNAME/ionic-god-mode-indexer.git](https://github.com/YOUR_USERNAME/ionic-god-mode-indexer.git)
   cd ionic-god-mode-indexer

   Download the required dependencies (gorilla/websocket and fatih/color)
   go mod tidy

   🛠️ Usage
Because it's written in Go, the tool compiles and runs instantly.

Bash
go run main.go

