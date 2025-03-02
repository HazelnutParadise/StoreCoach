package app

import (
	"reflect"
	"testing"

	"github.com/joho/godotenv"
)

func Test_splitIntoChunks(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	length := 3
	want := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
		{10},
	}
	got := splitIntoChunks(slice, length)
	t.Log(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("splitIntoChunks(%v, %d) = %v; want %v", slice, length, got, want)
	}
}

func Test_generateAspectsFromReviews(t *testing.T) {
	godotenv.Load("../StoreCoach.env")
	storeName := "好旺來餐廳"
	reviews := []string{
		"這家店的餐點很好吃，價格也很實惠。",
		"環境不錯，服務也很好。",
		"這家店的態度很差，也不好吃。",
		"so disgusting",
		"drinks are too expensive",
		"價格偏高，但是餐點很美味。",
		"不衛生，價格也不便宜。",
	}

	got, err := generateAspectsFromReviews(storeName, reviews)
	if err != nil {
		t.Errorf("generateTopics(%v) error: %v", reviews, err)
	}
	t.Log(got)

}

func TestReviewMining(t *testing.T) {
	godotenv.Load("../StoreCoach.env")
	storeName := "好旺來餐廳"
	reviews := []string{
		"這家店的餐點很好吃，價格也很實惠。",
		"環境不錯，服務也很好。",
		"這家店的態度很差，也不好吃。",
		"so disgusting",
		"drinks are too expensive",
		"價格偏高，但是餐點很美味。",
		"不衛生，價格也不便宜。",
	}

	got, err := ReviewMining(storeName, reviews)
	t.Log(got)
	if err != nil {
		t.Errorf("ReviewMining(%v) error: %v", reviews, err)
	}
	t.Log(got)
}
