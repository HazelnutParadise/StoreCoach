package app

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestAnalyzeReview(t *testing.T) {
	godotenv.Load("../StoreCoach.env")
	// 測試單個評論分析
	storeName := "星巴克"
	productName := "咖啡"
	review := "這家星巴克的咖啡很香，服務態度也很好，但是價格有點貴，環境很舒適，值得推薦。"
	attributes := []string{"咖啡品質", "服務態度", "價格", "環境", "整體評價"}

	result, err := analyzeReview(storeName, productName, review, attributes)
	if err != nil {
		t.Errorf("分析評論時發生錯誤: %v", err)
		return
	}
	t.Logf("=== 評論分析結果 ===")
	t.Logf("評論內容: %s", result.ReviewContent)
	t.Logf("評論評分: %d", result.ReviewRating)
	t.Logf("分析結果:")

	for _, mining := range result.MiningResults {
		t.Logf("  屬性: %s | 情感: %s | 分數: %d",
			mining.Attribute, mining.Sentiment, mining.Score)
	}
}

func TestGenerateAttributesFromReviews(t *testing.T) {
	godotenv.Load("../StoreCoach.env")
	// 測試屬性生成
	storeName := "麥當勞"
	productName := "漢堡"
	reviews := []string{
		"這個漢堡很好吃，肉很juicy，但是薯條有點鹹",
		"服務很快，但是環境有點吵，價格還算合理",
		"食物品質不錯，但是包裝有點草率，整體還可以",
		"等待時間太長，但是味道還不錯，下次還會來",
		"價格便宜，分量足夠，但是口味偏重",
	}

	attributes, err := generateAttributesFromReviews(storeName, productName, reviews)
	if err != nil {
		t.Errorf("生成屬性時發生錯誤: %v", err)
		return
	}
	t.Logf("=== 生成的屬性 ===")
	for i, attr := range attributes {
		t.Logf("%d. %s", i+1, attr)
	}
	t.Logf("總共生成了 %d 個屬性", len(attributes))
}

func TestReviewMiningComplete(t *testing.T) {
	godotenv.Load("../StoreCoach.env")
	// 測試完整的評論挖掘流程
	storeName := "台積電餐廳"
	productName := ""
	reviews := []string{
		"餐點種類豐富，價格實惠，但是排隊時間有點長",
		"菜色新鮮，服務人員態度親切，用餐環境乾淨舒適",
		"份量足夠，口味不錯，但是有些菜品會太鹹",
		"價格便宜，CP值很高，但是尖峰時間很難找位子",
		"食物品質穩定，選擇多樣化，整體體驗很好",
	}
	ratings := []uint8{4, 5, 3, 4, 5}

	result, err := ReviewMining(storeName, productName, reviews, ratings)
	if err != nil {
		t.Errorf("評論挖掘時發生錯誤: %v", err)
		return
	}
	t.Logf("=== 完整評論挖掘結果 ===")
	t.Logf("商店名稱: %s", result.StoreName)
	t.Logf("產品名稱: %s", result.ProductName)

	t.Logf("\n生成的屬性 (%d個):", len(result.Attributes))
	for i, attr := range result.Attributes {
		t.Logf("%d. %s", i+1, attr)
	}

	t.Logf("\n屬性平均分數:")
	for attr, score := range result.AverageAttributeScores {
		t.Logf("  %s: %.2f", attr, score)
	}

	t.Logf("\n線性迴歸結果:")
	for attr, regress := range result.SimpleLinearRegress {
		t.Logf("  %s:", attr)
		t.Logf("    斜率: %.4f", regress.Slope)
		t.Logf("    截距: %.4f", regress.Intercept)
		t.Logf("    R²: %.4f", regress.R2)
		t.Logf("    P值: %.4f", regress.PValue)
	}
	t.Logf("\nT檢定結果:")
	for attr, ttest := range result.TTest {
		t.Logf("  %s:", attr)
		t.Logf("    T值: %.4f", ttest.TValue)
		t.Logf("    P值: %.4f", ttest.PValue)
		t.Logf("    自由度: %d", ttest.Df)
	}

	t.Logf("\n摘要:\n%s", result.Summary)
}

func TestLLMAttributeConsistency(t *testing.T) {
	godotenv.Load("../StoreCoach.env")
	// 測試 LLM 返回的屬性一致性
	storeName := "手搖飲料店"
	productName := "珍珠奶茶"
	review := "這杯珍珠奶茶的茶味很濃，珍珠Q彈，甜度剛好，但是價格有點高。"

	// 測試多次調用看結果是否一致
	t.Logf("=== LLM 屬性一致性測試 ===")
	t.Logf("評論: %s\n", review)

	predefinedAttributes := []string{"茶味", "珍珠口感", "甜度", "價格", "整體評價"}

	for i := 0; i < 3; i++ {
		t.Logf("第 %d 次調用:", i+1)
		result, err := analyzeReview(storeName, productName, review, predefinedAttributes)
		if err != nil {
			t.Logf("  錯誤: %v", err)
			continue
		}

		t.Logf("  評分: %d", result.ReviewRating)
		for _, mining := range result.MiningResults {
			t.Logf("    屬性: %-10s | 情感: %-8s | 分數: %d",
				mining.Attribute, mining.Sentiment, mining.Score)
		}
		t.Logf("")
	}
}
