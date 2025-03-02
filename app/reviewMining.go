package app

import (
	"fmt"
	"math/rand"
	"strings"
)

type SingleReviewMiningResult struct {
	ReviewContent string `json:"reviewContent"`
	MiningResult  []struct {
		Aspect    string `json:"aspect"`
		Sentiment string `json:"sentiment"`
	}
}

func ReviewMining(storeName string, reviews []string) ([]SingleReviewMiningResult, error) {
	// **將評論分塊並生成構面**
	aspects, err := generateAspectsFromReviews(storeName, reviews)
	if err != nil {
		return nil, err
	}

	// **分析評論**
	var results []SingleReviewMiningResult
	for _, review := range reviews {
		result, err := analyzeReview(storeName, review, aspects)
		if err != nil {
			return nil, err
		}
		result.ReviewContent = review
		results = append(results, *result)
	}

	return results, nil
}

func generateAspectsFromReviews(storeName string, reviews []string) ([]string, error) {
	var allAspects []string
	// **迭代 5 次以提升精準度**
	for range 5 {
		shuffledReviews := randomSort(reviews)         // **隨機排序評論**
		chunks := splitIntoChunks(shuffledReviews, 50) // **分塊處理評論**
		for _, chunk := range chunks {
			prompt := "您是全球頂尖的資料探勘專家，正在為名為「" + fmt.Sprintf("%v", storeName) + "」的店家製作精準的分析報告。讓我們一步一步思考來捕捉評論中的各個構面:\n" +
				"1. 首先，仔細閱讀下列評論，並從中歸納出能夠完整反映評論主旨的各項具體構面（注意：單一評論可能同時涉及多個構面）:\n" +
				"`" + fmt.Sprintf("%q\n", chunk) + "`\n" +
				"2. 務必使用精確、專業且具體的詞彙描述每個構面，確保捕捉評論中的所有關鍵意涵。\n" +
				"3. 僅選取評論中有明確提及的構面，切勿自行添加評論中未出現的構面，並且僅回傳構面名稱，不附加括號內說明或註解。\n" +
				"4. 參考以下已識別的構面，並根據評論內容進行調整、合併或刪除，以確保結果準確反映評論：\n" +
				"`" + fmt.Sprintf("%q\n", allAspects) + "`\n" +
				"5. 最後，以 JSON 陣列格式回傳結果，格式：`[\"構面1\", \"構面2\", \"構面3\"]`\n" +
				"6. 只回傳 JSON 陣列，不要回傳其他文字。\n" +
				"違反任何一條規則將會讓您蒙受鉅額損失，請務必仔細閱讀並遵守以上規則。"

			resp, err := CallLLM(prompt)
			if err != nil {
				return nil, err
			}

			// **去除 LLM 輸出中的 JSON 標記**
			resp = strings.Replace(resp, "```json", "", -1)
			resp = strings.Replace(resp, "```", "", -1)

			resSlice := make([]string, 0)
			// **解析 LLM 輸出**
			err = json.Unmarshal([]byte(resp), &resSlice)
			if err != nil {
				return nil, err
			}

			// **合併新構面，確保沒有重複**
			allAspects = mergeSlicesUnique(allAspects, resSlice)
		}
	}

	return allAspects, nil
}

func analyzeReview(storeName, review string, aspects []string) (*SingleReviewMiningResult, error) {
	prompt := "您是一位全球頂尖的資料探勘專家，正在分析「" + fmt.Sprintf("%v", storeName) + "」（一間商店）的評論。請根據以下評論，找出涉及的構面，並對每個構面進行情緒分類。\n" +
		"評論: `" + fmt.Sprintf("%v", review) + "`\n" +
		"- 請根據已識別的構面，找出評論涉及的類別，並對每個類別進行情緒分類。\n" +
		"- 全部的構面如下:\n" +
		"`" + fmt.Sprintf("%q\n", aspects) + "`\n" +
		"- 情緒分類包括：positive, negative, neutral。\n" +
		"- 回傳 JSON 陣列，格式：[{\"aspect\": \"構面1\", \"sentiment\": \"positive\"}, {\"aspect\": \"構面2\", \"sentiment\": \"negative\"}]\n" +
		"- 只回傳 JSON 陣列，不要回傳其他文字。\n" +
		"違反任何一條規則將會讓您蒙受鉅額損失，請務必仔細閱讀並遵守以上規則。"

	resp, err := CallLLM(prompt)
	if err != nil {
		return nil, err
	}

	// **去除 LLM 輸出中的 JSON 標記**
	resp = strings.Replace(resp, "```json", "", -1)
	resp = strings.Replace(resp, "```", "", -1)

	var analysisResp = SingleReviewMiningResult{}
	err = json.Unmarshal([]byte(resp), &analysisResp.MiningResult)
	if err != nil {
		return nil, err
	}
	return &analysisResp, nil
}

// **合併去重的構面函數**
func mergeSlicesUnique[T any](existing []T, newElements []T) []T {
	elementMap := make(map[any]bool)
	for _, el := range existing {
		elementMap[el] = true
	}
	for _, el := range newElements {
		if !elementMap[el] {
			existing = append(existing, el)
			elementMap[el] = true
		}
	}
	return existing
}

func randomSort[T any](slice []T) []T {
	var newSlice = make([]T, len(slice))
	copy(newSlice, slice)
	rand.Shuffle(len(slice), func(i, j int) {
		newSlice[i], newSlice[j] = newSlice[j], newSlice[i]
	})
	return slice
}
