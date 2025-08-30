# Binance OHLC Data Downloader

A desktop application for downloading historical OHLC (Open, High, Low, Close) data from Binance exchange. Built with Wails, Go, and Vue.js.

## Features

- Download OHLC data for any trading pair on Binance
- Support for both Spot and Futures markets
- Multiple timeframes available (1m, 5m, 15m, 1h, 4h, 1d, etc.)
- Date range filtering
- Export data in JSON or CSV format
- Real-time download progress tracking

## Installation

### Prerequisites

- Go 1.18 or higher
- Node.js 16 or higher
- npm or yarn

### Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/binance-ohlc-downloader.git
   cd binance-ohlc-downloader
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   cd frontend
   npm install
   cd ..
   ```

3. Build the application:
   ```bash
   wails build
   ```

4. The built application will be available in the `build/bin` directory.

## Usage

1. Launch the application
2. Select market type (Spot or Futures)
3. Choose a trading pair from the dropdown
4. Select your desired timeframe
5. (Optional) Set a date range
6. Choose output format (JSON or CSV)
7. Click "Download" to save the data

## Development

### Running in Development Mode

```bash
wails dev
```

This will start the development server with hot-reload for both frontend and backend.

### Project Structure

- `app.go` - Main Go application with Binance API integration
- `frontend/` - Vue.js frontend source code
  - `src/App.vue` - Main application component
  - `src/main.ts` - Application entry point

## License

MIT

## Disclaimer

This software is for educational purposes only. Use it at your own risk. The author is not responsible for any losses incurred while using this software.
