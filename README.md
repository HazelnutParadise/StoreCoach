# StoreCoach 商店教練

StoreCoach 為店家提供智慧化分析，讓店家能夠更了解自己的營運狀況，制定更好的經營策略。

網址：https://storecoach.hazelnut-paradise.com/

## 套件設定位置

Go 後端在根目錄使用 `go.mod` 與 `go.sum` 管理依賴。Bun 前端在 `frontend/` 使用 `package.json` 與 `bun.lock`，前端套件指令都在該目錄執行。

更新直接與間接依賴時，Bun 的間接依賴仍須符合上游套件宣告的版本範圍。跨越該範圍前必須確認相容性，並重新執行檢查、建置與瀏覽器驗證。

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

現有評論探勘整合測試會讀取根目錄 `.env` 並呼叫真實 Gemini API，需要有效金鑰及可用模型，也可能產生 API 費用。不呼叫外部服務的升級回歸測試可執行：

```sh
go test ./... -run '^TestGemini|^TestStatsDependencyUpgrade$|^Test_splitIntoChunks$'
```

## Docker Compose

`Dockerfile` 放在專案根目錄，建置內容涵蓋前端與 Go 後端，前端產物會嵌入 Go 執行檔。Compose 的 `web` 服務從原始碼建置，不需要先準備應用程式映像。

Compose 預設讀取根目錄 `.env`，也接受 shell 環境變數。`GEMINI_API_KEY` 必須設定且不能為空。

```sh
docker compose config --quiet
docker compose up -d --build
```

服務沿用 `127.0.0.1:6868`、MongoDB 資料目錄 `/mnt/hdd0/storecoach/db`，以及既有的外部 `ollama-net` 網路。啟動前需確認資料目錄可用，並已有提供 Ollama 的該網路。MongoDB 使用官方映像，應用程式使用根目錄 Dockerfile 建置。
