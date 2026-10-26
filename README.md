# StoreCoach 商店教練

StoreCoach 為店家提供智慧化分析，讓店家能夠更了解自己的營運狀況，制定更好的經營策略。

網址：https://storecoach.hazelnut-paradise.com/

## 套件設定位置

Go 後端在根目錄使用 `go.mod` 與 `go.sum` 管理依賴。Bun 前端在 `frontend/` 使用 `package.json` 與 `bun.lock`，前端套件指令都在該目錄執行。

更新直接與間接依賴時，Bun 的間接依賴仍須符合上游套件宣告的版本範圍。跨越該範圍前必須確認相容性，並重新執行檢查、建置與瀏覽器驗證。

`go.mod` 保留 MongoDB v1 的間接版本要求，將上游依賴圖固定在已修補的版本。應用程式的資料庫程式使用 v2 SDK。執行 `go mod tidy` 或更新依賴後，需確認依賴圖中的 v1 版本仍不低於 `v1.17.7`。

## 開發環境

使用 Bun 1.4.4 與 Go 1.27.3。前端的安裝、開發伺服器與打包都由 Bun 執行，不需要另外安裝 Node.js。

在專案根目錄複製 `.env.example` 為 `.env`，填入 `GEMINI_API_KEY`。已有 `.env` 時保留原檔，不要覆寫。`.env` 不加入版控，也不會放進 Docker 建置內容或映像。

```sh
cd frontend
bun install --frozen-lockfile
bun run dev
```

打包與檢查：

```sh
cd frontend
bun run lint
bun run build
cd ..
go build .
go vet ./...
go test ./...
```

現有評論探勘整合測試需要可連線的 Ollama 服務與 `gemma3:12b` 模型，以及根目錄 `.env` 的有效 Gemini 金鑰。Gemini 呼叫可能產生 API 費用。不呼叫外部服務的升級回歸測試可執行：

```sh
go test ./... -run '^TestCallLLM|^TestLLM(IdleCPU|CPUProbeHelper|Queue)|^TestGemini|^TestStatsDependencyUpgrade$|^Test_splitIntoChunks$|^TestReportPageFallback$'
```

## LLM 請求佇列

使用 Insyra `engine/ring` 維護 Gemini 與 Ollama 各自的先進先出佇列。所有佇列操作共用同一把鎖，空佇列透過 `sync.Cond` 等待新工作，每筆請求透過自己的 channel 接收結果與錯誤。

同一時間最多執行一筆模型請求，兩個模型都有等待工作時輪流處理。Gemini 維持原有每筆請求前後合計 2,100 毫秒的等待間隔，不因更換佇列而增加 API 呼叫頻率。

## Docker Compose

`Dockerfile` 放在專案根目錄，建置內容涵蓋前端與 Go 後端，前端產物會嵌入 Go 執行檔。Compose 的 `web` 服務從原始碼建置，不需要先準備應用程式映像。

Compose 預設讀取根目錄 `.env`，也接受 shell 環境變數。`GEMINI_API_KEY` 必須設定且不能為空。

```sh
docker compose config --quiet
docker compose up -d --build
```

服務沿用 `127.0.0.1:6868`、MongoDB 資料目錄 `/mnt/hdd0/storecoach/db`，以及提供 Ollama 的既有 `infra-net` 外部網路，Compose 內以 `ollama-net` 引用。啟動前需確認資料目錄與該網路可用。MongoDB 使用官方映像，應用程式使用根目錄 Dockerfile 建置。

既有 MongoDB 資料若仍使用 7.0 相容格式，不能直接啟動 9.0。請先停止寫入並完整備份，在隔離副本驗證 8.0 過渡升級，將相容格式設為 8.0 後，才以 9.0 啟動。正式切換前比對各集合筆數與內容。參考 [MongoDB 7.0 升級至 8.0](https://www.mongodb.com/docs/manual/release-notes/8.0-upgrade-standalone/) 與 [8.0 升級至 9.0](https://www.mongodb.com/docs/manual/release-notes/9.0-upgrade-from-8.0-standalone/)。

過渡升級使用的 MongoDB 8.0.34 若將 Ubuntu 的 `7.0.0` 核心誤判為不相容，先核對 `/proc/version_signature`。本次主機實際為已修正的 7.0.14，臨時 8.0 容器使用 `GLIBC_TUNABLES=glibc.pthread.rseq=1`，並確認 `db.serverStatus().tcmalloc.usingPerCPUCaches` 為 `false`。正式 MongoDB 9.0.4 沿用預設設定。參考 [MongoDB 記憶體配置說明](https://www.mongodb.com/docs/manual/administration/tcmalloc-performance/)。
