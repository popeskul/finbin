# Finbin Project 📈

## Overview 🚀
Finbin is a robust application designed to interact with the Binance API, facilitating the monitoring and analysis of cryptocurrency market data.

## Features 🌟
- Fetch real-time price data for various cryptocurrency symbols.
- Efficiently process and store data for later analysis.
- Graceful shutdown for smooth termination of the application.

## Prerequisites 📋
Before you begin, ensure you have met the following requirements:
- You have installed the latest version of Go.
- You have a Binance account with API keys.

## Installation 💻
Clone the repository to your local machine:
```sh
git clone https://github.com/popeskul/finbin.git
cd finbin
```

## Configuration ⚙️
Create a .env file in the project directory with your Binance API credentials:
```sh
BINANCE_API_KEY=your_api_key
BINANCE_SECRET=your_api_secret
GRACEFUL_TIMEOUT=5s
LOG_LEVEL=info
```
You can adjust the GRACEFUL_TIMEOUT and LOG_LEVEL values as needed.

## Usage 🚀
```sh
make run
```
