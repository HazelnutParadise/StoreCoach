package app

import (
	"fmt"
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

func Test_generateAttributesFromReviews(t *testing.T) {
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

	got, err := generateAttributesFromReviews(storeName, reviews)
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
		"味道一般，價格偏高。",
		"不是最美味，但也不難吃。",
	}

	got, err := ReviewMining(storeName, reviews)
	print(fmt.Sprint(got), "\n")
	// for _, v := range got.Results {
	// 	print(fmt.Sprint(v), "\n")
	// }
	if err != nil {
		t.Errorf("ReviewMining(%v) error: %v", reviews, err)
	}
}
