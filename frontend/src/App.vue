<script lang="ts" setup>
import {reactive, onMounted, watch} from 'vue'
import {FetchPairs, DownloadOHLC} from '../wailsjs/go/main/App'
import {EventsOn} from '../wailsjs/runtime'

const data = reactive({
  marketType: 'spot',
  pair: '',
  pairs: [] as string[],
  timeframe: '1h',
  timeframes: ['1m', '3m', '5m', '15m', '30m', '1h', '2h', '4h', '6h', '8h', '12h', '1d', '3d', '1w', '1M'],
  status: 'Select options and click Download.',
  startTime: '2021-01-01',
  endTime: new Date().toISOString().split('T')[0],
  downloading: false,
  progress: 0,
  outputFormat: 'json',
  outputFormats: ['json', 'csv'],
})

function fetchPairs() {
  data.status = `Fetching pairs for ${data.marketType}...`
  FetchPairs(data.marketType).then((result: string[]) => {
    data.pairs = result
    if (result.length > 0) {
      data.pair = result[0]
      data.status = 'Ready.'
    } else {
      data.status = `No pairs found for ${data.marketType}.`
    }
  }).catch((err: any) => {
    data.status = `Error fetching pairs: ${err}`
  })
}

function download() {
  if (!data.pair) {
    data.status = 'Please select a trading pair.'
    return
  }
  data.downloading = true
  data.progress = 0
  data.status = `Downloading ${data.pair} ${data.timeframe} data...`
  DownloadOHLC(data.pair, data.timeframe, data.marketType, data.startTime, data.endTime, data.outputFormat).then((result: string) => {
    data.status = result
    data.downloading = false
  }).catch((err: any) => {
    data.status = `Error: ${err}`
    data.downloading = false
  })
}

watch(() => data.marketType, () => {
  fetchPairs()
})

onMounted(() => {
  fetchPairs()
  EventsOn('download_progress', (progress: number) => {
    data.progress = progress
  })
})

</script>

<template>
  <main>
    <div class="container">
      <h1>Binance OHLC Downloader</h1>

      <div class="control-group">
        <label for="marketType">Market Type:</label>
        <select id="marketType" v-model="data.marketType">
          <option value="spot">Spot</option>
          <option value="futures">Futures</option>
        </select>
      </div>

      <div class="control-group">
        <label for="pair">Trading Pair:</label>
        <select id="pair" v-model="data.pair">
          <option v-for="p in data.pairs" :key="p" :value="p">{{ p }}</option>
        </select>
      </div>

      <div class="control-group">
        <label for="timeframe">Timeframe:</label>
        <select id="timeframe" v-model="data.timeframe">
          <option v-for="tf in data.timeframes" :key="tf" :value="tf">{{ tf }}</option>
        </select>
      </div>

      <div class="control-group">
        <label for="startTime">Start Date:</label>
        <input type="date" id="startTime" v-model="data.startTime" class="input" />
      </div>

      <div class="control-group">
        <label for="endTime">End Date:</label>
        <input type="date" id="endTime" v-model="data.endTime" class="input" />
      </div>

      <div class="control-group">
        <label for="outputFormat">Output Format:</label>
        <select id="outputFormat" v-model="data.outputFormat" class="input">
          <option v-for="format in data.outputFormats" :key="format" :value="format">{{ format.toUpperCase() }}</option>
        </select>
      </div>

      <div class="control-group">
        <button class="btn" @click="download" :disabled="data.downloading">Download</button>
      </div>

      <div class="progress-container" v-if="data.downloading">
        <progress :value="data.progress" max="100"></progress>
        <span>{{ data.progress }}%</span>
      </div>

      <div id="status" class="status">{{ data.status }}</div>
    </div>
  </main>
</template>

<style scoped>
.container {
  padding: 2rem;
  text-align: center;
}

h1 {
  margin-bottom: 2rem;
}

.control-group {
  margin-bottom: 1.5rem;
  display: flex;
  justify-content: center;
  align-items: center;
}

label {
  margin-right: 1rem;
  min-width: 120px;
  text-align: right;
}

select,
.btn,
.input {
  min-width: 200px;
  height: 30px;
  line-height: 30px;
  border-radius: 3px;
  border: 1px solid #ccc;
  padding: 0 8px;
}

.btn {
  cursor: pointer;
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.progress-container {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 1rem;
  justify-content: center;
  margin-top: 1rem;
}

progress {
  width: 70%;
}

.status {
  margin-top: 2rem;
  font-style: italic;
  color: #555;
}
</style>
