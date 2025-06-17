package app

import (
	"math"
	"math/rand"
	"time"

	"github.com/HazelnutParadise/Go-Utils/conv"
	"github.com/HazelnutParadise/Go-Utils/sliceutil"
	"github.com/HazelnutParadise/insyra"
	"github.com/HazelnutParadise/insyra/stats"
)

type ReviewMiningStruct struct {
	DataUUID               string                                                    `bson:"dataUUID"`
	StoreName              string                                                    `json:"storeName" bson:"storeName"`
	ProductName            string                                                    `json:"productName" bson:"productName"`
	Attributes             []string                                                  `json:"attributes" bson:"attributes"`
	Results                []SingleReviewMiningResult                                `json:"results" bson:"results"`
	AverageAttributeScores map[string]float64                                        `json:"averageAttributeScores" bson:"averageAttributeScores"`
	SimpleLinearRegress    map[string]ReviewMiningAttributeSimpleLinearRegressResult `json:"simpleLinearRegress" bson:"simpleLinearRegress"`
	TTest                  map[string]ReviewMiningAttributeTTestResult               `json:"tTest" bson:"tTest"`
	Summary                string                                                    `json:"summary" bson:"summary"`
	Timestamp              int64                                                     `json:"timestamp" bson:"timestamp"`
}

type SingleReviewMiningResult struct {
	ReviewContent string `json:"reviewContent" bson:"reviewContent"`
	ReviewRating  uint8  `json:"reviewRating" bson:"reviewRating"`
	MiningResults []struct {
		Attribute string `json:"attribute" bson:"attribute"`
		Sentiment string `json:"sentiment" bson:"sentiment"`
		Score     uint8  `json:"score" bson:"score"`
	} `json:"miningResults" bson:"miningResults"`
}

type ReviewMiningAttributeSimpleLinearRegressResult struct {
	Slope         float64    `json:"slope" bson:"slope"`
	Intercept     float64    `json:"intercept" bson:"intercept"`
	R2            float64    `json:"r2" bson:"r2"`
	PValue        float64    `json:"pValue" bson:"pValue"`
	StandardError float64    `json:"standardError" bson:"standardError"`
	CI            [2]float64 `json:"ci" bson:"ci"`
}

type ReviewMiningAttributeTTestResult struct {
	TValue float64 `json:"tValue" bson:"tValue"`
	PValue float64 `json:"pValue" bson:"pValue"`
	Df     int     `json:"df" bson:"df"`
}

func ReviewMining(storeName string, productName string, reviews []string, ratings []uint8) (reviewMiningResultInfo *ReviewMiningStruct, err error) {
	// **將評論分塊並生成屬性**
	attributes, err := generateAttributesFromReviews(storeName, productName, reviews)
	if err != nil {
		return nil, err
	}

	// **分析評論**
	var results []SingleReviewMiningResult
	for i, review := range reviews {
		result, err := analyzeReview(storeName, productName, review, attributes)
		if err != nil {
			return nil, err
		}
		result.ReviewContent = review
		// 順便存評分
		if ratings != nil && len(ratings) == len(reviews) {
			result.ReviewRating = ratings[i]
		}

		results = append(results, *result)
	}

	// **計算屬性平均分數**
	averageAttributeScores := calculateAttributeAverageScores(results)

	// **針對每個屬性簡單線性迴歸**
	simpleLinearRegressResults := simpleLinearRegressEachAttributeScoreAndRatings(attributes, results)

	// todo：多元迴歸分析
	// multipleLinearRegressAttributeScoresAndRatings(attributes, results)

	// **進行 T 檢定**
	ttest := reviewsAttributeTTest(attributes, results)

	reviewMiningStruct := ReviewMiningStruct{
		StoreName:              storeName,
		ProductName:            productName,
		Attributes:             attributes,
		Results:                results,
		AverageAttributeScores: averageAttributeScores,
		SimpleLinearRegress:    simpleLinearRegressResults,
		TTest:                  ttest,
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

func generateAttributesFromReviews(storeName string, productName string, reviews []string) ([]string, error) {
	var allAttributes []string
	// **迭代 3 次以提升精準度**
	for range 3 {
		shuffledReviews := randomSort(reviews)         // **隨機排序評論**
		chunks := splitIntoChunks(shuffledReviews, 50) // **分塊處理評論**
		for _, chunk := range chunks {
			var chunkedReviewsStr string
			for _, review := range chunk {
				chunkedReviewsStr += "  - 「" + review + "」\n"
			}
			beginningPrompt := "您是全球頂尖的資料探勘專家，正在幫名為「" + storeName + "」的店家或機構"
			if productName != "" {
				beginningPrompt += "的「" + productName + "」（一項產品或服務）"
			}
			prompt := beginningPrompt + "製作精準的分析報告。讓我們一步一步思考來捕捉評論中的各個屬性:\n" +
				"1. 首先，仔細閱讀下列評論，並從中歸納出能夠完整反映評論主旨的各項具體屬性（注意：單一評論可能同時涉及多個屬性）:\n" +
				chunkedReviewsStr +
				"2. 屬性的種類包括但不限於：\n" +
				"  - 屬性功能(例如：材質、顏色、價格、美感、保證、速度等)\n" +
				"  - 利益與用途(例如：安全、便利、易用性、銳利、耐用、效率、效能、性價比等)\n" +
				"  - 品牌個性(例如：浪漫、品味、歡樂、典雅、創新等)\n" +
				"  - 形象(與特定品牌相關的形象或印象)\n" +
				"3. 務必使用精確、專業且具體的詞彙描述每個屬性，確保捕捉評論中的所有關鍵意涵。\n" +
				"4. 屬性必須具備分類上的意義，不能只針對個別評論而設定。\n" +
				"5. 僅選取評論中有明確提及的屬性，切勿自行添加評論中未出現的屬性，並且僅回傳屬性名稱，不附加括號內說明或註解。\n" +
				"6. 參考以下已識別的屬性，並根據評論內容進行調整、增加或合併，以確保結果準確反映評論：\n" +
				"  `" + conv.ToString(allAttributes) + "`\n" +
				"7. 最後，以 JSON 陣列格式回傳結果，格式：`[\"屬性1\", \"屬性2\", \"屬性3\"]`\n" +
				"8. 只回傳 JSON 陣列，不要回傳其他文字。\n" +
				"違反任何一條規則將會讓您蒙受鉅額損失，請務必仔細閱讀並遵守以上規則。"

			resp, err := CallLLM(prompt, 0.1, LLM_Local)
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

func analyzeReview(storeName, productName string, review string, attributes []string) (*SingleReviewMiningResult, error) {
	sentimentClassification := [3]string{"positive", "negative", "neutral"}
	beginningPrompt := "您是一位全球頂尖的資料探勘專家，正在分析「" + storeName + "」（一間商店或機構）"
	if productName != "" {
		beginningPrompt += "的「" + productName + "」（一項產品或服務）"
	}
	prompt := beginningPrompt + "的評論。\n" +
		"評論：「" + review + "」\n\n" +
		"請根據以上評論，找出涉及的屬性，並確實完成以下兩項任務：\n" +
		"1. 對每個屬性進行情緒分類。\n" +
		"	- 請根據已識別的屬性，找出評論涉及的類別，並對每個類別進行情緒分類。\n" +
		"	- 全部的屬性：`" + conv.ToString(attributes) + "`\n" +
		"	- 情緒分類包括：`" + conv.ToString(sentimentClassification) + "`。\n" +
		"2. 對每個評論涉及的屬性給出0至10分的評分\n" +
		"	- 評分標準：\n" +
		"		- 0分：完全無法接受，災難性的體驗或完全缺乏該屬性\n" +
		"		- 1分：極差，接近最糟，強烈不建議\n" +
		"		- 2分：非常差，幾乎無法接受，極度不滿意\n" +
		"		- 3分：差，體驗不佳，有多處需要改善\n" +
		"		- 4分：普通偏差，有明顯缺點，不太想再使用\n" +
		"		- 5分：一般，剛好及格，不特別好也不差\n" +
		"		- 6分：還可以，略低於期望，有改進空間\n" +
		"		- 7分：普通偏好，有明顯優點也有些缺點\n" +
		"		- 8分：很好，有些小地方可以更好，但整體令人滿意\n" +
		"		- 9分：極好，幾乎沒有缺點，會想再次購買或使用\n" +
		"		- 10分：完美無缺，超出期望，令人驚艷，強烈推薦\n\n" +
		"- 你只能從現有的屬性中選擇，不能添加新的屬性。\n" +
		"- 特別留意評論是否為反諷\n" +
		"- 回傳 JSON 陣列，格式：`[{\"attribute\": \"屬性1\", \"sentiment\": \"positive\", \"score\": 6}, {\"attribute\": \"屬性2\", \"sentiment\": \"negative\", \"score\": 4}]`\n" +
		"- 只回傳 JSON 陣列，不要回傳其他文字。\n" +
		"違反任何一條規則將會讓您蒙受鉅額損失，請務必仔細閱讀並遵守以上規則。"

	resp, err := CallLLM(prompt, 0.1, LLM_Local)
	if err != nil {
		return nil, err
	}

	var analysisResp = SingleReviewMiningResult{}
	err = Unmarshal_LLM_JSON_Response(resp, &analysisResp.MiningResults)
	if err != nil {
		return nil, err
	}
	return &analysisResp, nil
}

func addSummaryForReviewMining(resultInfo *ReviewMiningStruct) error {
	beginningPrompt := "您是一位全球頂尖的商店教練（Store Coach），正在幫名為「" + resultInfo.StoreName + "」的店家或機構"
	if resultInfo.ProductName != "" {
		beginningPrompt += "的「" + resultInfo.ProductName + "」（一項產品或服務）"
	}
	prompt := beginningPrompt + "製作精準的分析報告。請根據以下評論分析結果，撰寫一份簡潔的摘要，概述評論中的主要屬性及其情緒分類。\n" +
		"所有屬性：`" + conv.ToString(resultInfo.Attributes) + "`\n" +
		"分析結果：\n" +
		"```json\n" + conv.ToString(resultInfo.Results) + "\n```\n" +
		"- 請根據以上分析結果，撰寫一份簡潔的摘要，包括以下內容：\n" +
		"  1. overall：對評論的整體情緒進行總結。\n" +
		"  2. advantages：列舉評論中提到的主要優點。\n" +
		"  3. disadvantages：列舉評論中提到的主要缺點。\n" +
		"  4. neutral：列舉評論中提到的中性評價。\n" +
		"  5. 教練的建議（suggestions）：根據分析結果，提出至少3點建議（最好5點以上），幫助店家或機構改進。\n" +
		"- 回傳 JSON，格式：`\"{\"overall\": \"整體摘要內容\", \"advantages\": \"主要優點內容\", \"disadvantages\": \"主要缺點內容\", \"neutral\": \"中性評價內容\", \"suggestions\": [{\"title\": \"建議標題一\", \"content\": \"建議一內容說明\"}, {\"title\": \"建議標題二\", \"content\": \"建議二內容說明\"}]}\"`\n" +
		"- 只回傳 JSON 陣列，不要回傳其他文字。" +
		"- 確實填好每一欄，不要有任何留空。" +
		"- 在相應欄位中填寫內容就好，不要抄題" +
		"違反任何一條規則將會讓您蒙受鉅額損失，請務必仔細閱讀並遵守以上規則。"
	summary, err := CallLLM(prompt, 0.1, LLM_Gemini)
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

func calculateAttributeAverageScores(reviewMiningResults []SingleReviewMiningResult) map[string]float64 {
	scoreDT := insyra.NewDataTable()
	for _, result := range reviewMiningResults {
		for _, miningResult := range result.MiningResults {
			scoreDT.AppendRowsByColName(map[string]any{miningResult.Attribute: miningResult.Score})
		}
	}
	attributeScores := make(map[string]float64)
	_, numCols := scoreDT.Size()
	for i := range numCols {
		col := scoreDT.GetColByNumber(i)
		attributeName := col.GetName()
		attributeScores[attributeName] = col.Mean()
	}
	return attributeScores
}

func simpleLinearRegressEachAttributeScoreAndRatings(attributes []string, reviewMiningResults []SingleReviewMiningResult) (attributeRegressResultsMap map[string]ReviewMiningAttributeSimpleLinearRegressResult) {
	attributeRegressResultsMap = make(map[string]ReviewMiningAttributeSimpleLinearRegressResult)
	attributeScores := make(map[string]*insyra.DataList)
	ratingsFromReviewsMentionedAttribute := make(map[string]*insyra.DataList)
	for _, attribute := range attributes {
		attributeScores[attribute] = insyra.NewDataList()
		ratingsFromReviewsMentionedAttribute[attribute] = insyra.NewDataList()
	}
	for _, result := range reviewMiningResults {
		for _, miningResult := range result.MiningResults {
			attributeScores[miningResult.Attribute].Append(miningResult.Score)
			ratingsFromReviewsMentionedAttribute[miningResult.Attribute].Append(result.ReviewRating)
		}
	}
	// **對每個屬性進行簡單線性迴歸**
	for attribute, scores := range attributeScores {
		ratings := ratingsFromReviewsMentionedAttribute[attribute]
		if scores.Len() < 2 || ratings.Len() < 2 {
			continue // 略過這組數據
		}
		result := stats.LinearRegression(scores, ratings)
		if result == nil {
			continue // 如果線性迴歸失敗，跳過這個屬性
		}
		attributeRegressResultsMap[attribute] = ReviewMiningAttributeSimpleLinearRegressResult{
			Slope:         result.Slope,
			Intercept:     result.Intercept,
			R2:            result.RSquared,
			PValue:        result.PValue,
			StandardError: result.StandardError,
			CI:            result.ConfidenceIntervalSlope,
		}
	}
	return attributeRegressResultsMap
}

func multipleLinearRegressAttributeScoresAndRatings(attributes []string, reviewMiningResults []SingleReviewMiningResult) *stats.LinearRegressionResult {
	// todo:多元線性回歸
	ratingDL := insyra.NewDataList()
	attributeDLs := make(map[string][2]*insyra.DataList)
	for _, attribute := range attributes {
		attributeScoresDL := insyra.NewDataList()
		mentionedDL := insyra.NewDataList()
		attributeDLs[attribute] = [2]*insyra.DataList{attributeScoresDL, mentionedDL}
	}

	// 準備變項
	for _, result := range reviewMiningResults {
		attributesToCollect := make([]string, len(attributes))
		copy(attributesToCollect, attributes)
		for _, miningResult := range result.MiningResults {
			attributeDLs[miningResult.Attribute][0].Append(miningResult.Score)
			attributeDLs[miningResult.Attribute][1].Append(1)
			attributesToCollect = sliceutil.RemoveAll(attributesToCollect, miningResult.Attribute)
		}
		for _, attribute := range attributesToCollect {
			attributeDLs[attribute][0].Append(5) // 未提及的屬性分數記為 5
			attributeDLs[attribute][1].Append(0) // 未提及的屬性記為 0
		}
		ratingDL.Append(result.ReviewRating)
	}
	// todo
	// stats.LinearRegression()
	return nil
}

func reviewsAttributeTTest(attributes []string, miningResults []SingleReviewMiningResult) map[string]ReviewMiningAttributeTTestResult {
	count0ratings := 0
	for _, result := range miningResults {
		if result.ReviewRating == 0 {
			count0ratings++
		}
	}
	if count0ratings == len(miningResults) {
		return nil
	}
	var statResults = map[string]ReviewMiningAttributeTTestResult{}
	for _, attribute := range attributes {
		group0 := insyra.NewDataList()
		group1 := insyra.NewDataList()
	NEXT_RESULT:
		for _, result := range miningResults {
			for _, miningResult := range result.MiningResults {
				if miningResult.Attribute == attribute {
					group1.Append(result.ReviewRating)
					continue NEXT_RESULT
				}
			}
			group0.Append(result.ReviewRating)
		}
		var varEqual = false
		if group0.Len() < 2 || group1.Len() < 2 {
			continue // 略過這組數據
		}
		ftest := stats.FTestForVarianceEquality(group0, group1)
		if ftest != nil {
			if ftest.PValue <= 0.05 {
				varEqual = false
			} else {
				varEqual = true
			}
		}
		result := stats.TwoSampleTTest(group0, group1, varEqual, 0.95)
		if result == nil {
			continue
		}
		if math.IsNaN(result.Statistic) || math.IsNaN(result.PValue) {
			continue
		}
		statResults[attribute] = ReviewMiningAttributeTTestResult{
			TValue: result.Statistic,
			PValue: result.PValue,
			Df:     int(*result.DF),
		}
	}
	return statResults
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
