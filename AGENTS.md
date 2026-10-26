# Follow-ups

- [P1] `app/llm.go` 的 `CallLLM` 沒有將傳入的 `model` 寫入 `llmReq`，指定本機模型的請求仍會選到預設 Gemini。調整前需確認原本的模型分流需求，並補上 Gemini／Ollama 分流測試。

- [P2] `frontend/src/pages/ReviewMiningResult.jsx` 在報告缺少 `attributes` 時，顯示提示並呼叫 `navigate` 後仍執行 `rmAttributes.forEach`，會造成執行錯誤。另案補上缺漏資料的畫面測試並修正提前返回。
