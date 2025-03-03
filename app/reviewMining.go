package app

import (
	"math/rand"
	"time"

	"github.com/HazelnutParadise/Go-Utils/conv"
)

type ReviewMiningStruct struct {
	StoreName  string                     `json:"storeName"`
	Attributes []string                   `json:"attributes"`
	Results    []singleReviewMiningResult `json:"results"`
	Summary    string                     `json:"summary"`
	Timestamp  int64                      `json:"timestamp"`
}

type singleReviewMiningResult struct {
	ReviewContent string `json:"reviewContent"`
	MiningResult  []struct {
		Attribute string `json:"attribute"`
		Sentiment string `json:"sentiment"`
	}
}

func ReviewMining(storeName string, reviews []string) (reviewMiningResultInfo *ReviewMiningStruct, err error) {
	// **將評論分塊並生成屬性**
	attributes, err := generateAttributesFromReviews(storeName, reviews)
	if err != nil {
		return nil, err
	}

	// **分析評論**
	var results []singleReviewMiningResult
	for _, review := range reviews {
		result, err := analyzeReview(storeName, review, attributes)
		if err != nil {
			return nil, err
		}
		result.ReviewContent = review
		results = append(results, *result)
	}

	reviewMiningStruct := ReviewMiningStruct{
		StoreName:  storeName,
		Attributes: attributes,
		Results:    results,
	}

	// **為評論分析結果添加摘要**
	err = addSummaryForReviewMining(&reviewMiningStruct)
	if err != nil {
		return nil, err
	}

	// **添加時間戳**
	reviewMiningStruct.Timestamp = time.Now().Unix()

	reviewMiningResultInfo = &reviewMiningStruct
	return reviewMiningResultInfo, nil
}

func generateAttributesFromReviews(storeName string, reviews []string) ([]string, error) {
	var allAttributes []string
	// **迭代 5 次以提升精準度**
	for range 5 {
		shuffledReviews := randomSort(reviews)         // **隨機排序評論**
		chunks := splitIntoChunks(shuffledReviews, 50) // **分塊處理評論**
		for _, chunk := range chunks {
			var chunkedReviewsStr string
			for _, review := range chunk {
				chunkedReviewsStr += "  - 「" + review + "」\n"
			}

			prompt := "您是全球頂尖的資料探勘專家，正在幫名為「" + storeName + "」的店家製作精準的分析報告。讓我們一步一步思考來捕捉評論中的各個屬性:\n" +
				"1. 首先，仔細閱讀下列評論，並從中歸納出能夠完整反映評論主旨的各項具體屬性（注意：單一評論可能同時涉及多個屬性）:\n" +
				chunkedReviewsStr +
				"2. 屬性的種類包括但不限於：\n" +
				"  - 屬性功能(例如：材質、顏色、價格、美感、保證、速度等)\n" +
				"  - 利益與用途(例如：安全、便利、易用性、銳利、耐用、效率、效能、性價比等)\n" +
				"  - 品牌個性(例如：浪漫、品味、歡樂、典雅、創新等)\n" +
				"  - 形象(與特定品牌相關的形象或印象)\n" +
				"3. 務必使用精確、專業且具體的詞彙描述每個屬性，確保捕捉評論中的所有關鍵意涵。\n" +
				"4. 僅選取評論中有明確提及的屬性，切勿自行添加評論中未出現的屬性，並且僅回傳屬性名稱，不附加括號內說明或註解。\n" +
				"5. 參考以下已識別的屬性，並根據評論內容進行調整、合併或刪除，以確保結果準確反映評論：\n" +
				"  `" + conv.ToString(allAttributes) + "`\n" +
				"6. 最後，以 JSON 陣列格式回傳結果，格式：`[\"屬性1\", \"屬性2\", \"屬性3\"]`\n" +
				"7. 只回傳 JSON 陣列，不要回傳其他文字。\n" +
				"違反任何一條規則將會讓您蒙受鉅額損失，請務必仔細閱讀並遵守以上規則。"

			resp, err := CallLLM(prompt, 0.1)
			if err != nil {
				return nil, err
			}

			resSlice := make([]string, 0)
			// **解析 LLM 輸出**
			err = Unmarshal_LLM_JSON_Response(resp, &resSlice)
			if err != nil {
				return nil, err
			}

			// **合併新屬性，確保沒有重複**
			allAttributes = mergeSlicesUnique(allAttributes, resSlice)
		}
	}

	return allAttributes, nil
}

func analyzeReview(storeName, review string, attributes []string) (*singleReviewMiningResult, error) {
	sentimentClassification := [3]string{"positive", "negative", "neutral"}
	prompt := "您是一位全球頂尖的資料探勘專家，正在分析「" + storeName + "」（一間商店）的評論。請根據以下評論，找出涉及的屬性，並對每個屬性進行情緒分類。\n" +
		"評論：「" + review + "」\n" +
		"- 請根據已識別的屬性，找出評論涉及的類別，並對每個類別進行情緒分類。\n" +
		"- 全部的屬性：`" + conv.ToString(attributes) + "`\n" +
		"- 情緒分類包括：`" + conv.ToString(sentimentClassification) + "`。\n" +
		"- 回傳 JSON 陣列，格式：`[{\"attribute\": \"屬性1\", \"sentiment\": \"positive\"}, {\"attribute\": \"屬性2\", \"sentiment\": \"negative\"}]`\n" +
		"- 只回傳 JSON 陣列，不要回傳其他文字。\n" +
		"違反任何一條規則將會讓您蒙受鉅額損失，請務必仔細閱讀並遵守以上規則。"

	resp, err := CallLLM(prompt, 0.1)
	if err != nil {
		return nil, err
	}

	var analysisResp = singleReviewMiningResult{}
	err = Unmarshal_LLM_JSON_Response(resp, &analysisResp.MiningResult)
	if err != nil {
		return nil, err
	}
	return &analysisResp, nil
}

func addSummaryForReviewMining(resultInfo *ReviewMiningStruct) error {
	prompt := "您是一位全球頂尖的商店教練（Store Coach），正在幫名為「" + resultInfo.StoreName + "」的店家製作精準的分析報告。請根據以下評論分析結果，撰寫一份簡潔的摘要，概述評論中的主要屬性及其情緒分類。\n" +
		"所有屬性：`" + conv.ToString(resultInfo.Attributes) + "`\n" +
		"分析結果：\n" +
		"```json\n" + conv.ToString(resultInfo.Results) + "\n```\n" +
		"- 回傳 JSON，格式：`\"{\"overall\": \"整體摘要\", \"advantages\": \"主要優點\", \"disadvantages\": \"主要缺點\", \"neutral\": \"中性評價\", \"suggestions\": [{\"title\": \"建議標題一\", \"content\": \"建議一內容說明\"}, {\"title\": \"建議標題二\", \"content\": \"建議二內容說明\"}]}\"`\n" +
		"- 只回傳 JSON 陣列，不要回傳其他文字。"
	summary, err := CallLLM(prompt, 0.1)
	if err != nil {
		return err
	}
	respStruct := struct {
		Overall       string `json:"overall"`
		Advantages    string `json:"advantages"`
		Disadvantages string `json:"disadvantages"`
		Neutral       string `json:"neutral"`
		Suggestions   []struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		} `json:"suggestions"`
	}{}
	Unmarshal_LLM_JSON_Response(summary, &respStruct)
	summary = respStruct.Overall + "\n" +
		"\n主要優點：\n" + respStruct.Advantages + "\n" +
		"\n主要缺點：\n" + respStruct.Disadvantages + "\n" +
		"\n中性評價：\n" + respStruct.Neutral + "\n"
	summary += "\n教練的建議：\n"
	i := 1
	for _, suggestion := range respStruct.Suggestions {
		summary += conv.ToString(i) + ". " + suggestion.Title + "：" + suggestion.Content + "\n"
		i++
	}
	// 去除多餘的換行符
	resultInfo.Summary = summary[:len(summary)-1]
	return nil
}

// **合併去重的屬性函數**
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
