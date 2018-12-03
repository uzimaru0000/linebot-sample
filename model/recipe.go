package model

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/poccariswet/bot/template"
)

type Category struct {
	ID     int    `json:"categoryId"`
	Name   string `json:"categoryName"`
	Parent *int   `json:"parentCategoryId"`
}

type Categories struct {
	Large  []Category `json:"large"`
	Medium []Category `json:"medium"`
	Small  []Category `json:"small"`
}

type Recipe struct {
	Title        string `json:"recipeTitle"`
	RecipeURL    string `json:"recipeUrl"`
	FoodImageURL string `json:"foodImageUrl"`
	RecipeCost   string `json:"recipeCost"`
}

func send(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetCategory(token string) (*Categories, error) {
	endPoint := "https://app.rakuten.co.jp/services/api/Recipe/CategoryList/20170426"

	value := &url.Values{}
	value.Add("applicationId", token)
	value.Add("format", "json")

	raw, err := send(endPoint + "?" + value.Encode())
	if err != nil {
		return nil, err
	}

	var cate struct {
		Result *Categories
	}
	json.Unmarshal(raw, &cate)

	return cate.Result, nil
}

func GetRecipe(token string, id string) ([]Recipe, error) {
	endPoint := "https://app.rakuten.co.jp/services/api/Recipe/CategoryRanking/20170426"

	value := &url.Values{}
	value.Add("applicationId", token)
	value.Add("categoryId", id)
	value.Add("format", "json")

	raw, err := send(endPoint + "?" + value.Encode())
	if err != nil {
		return nil, err
	}

	var recipe struct {
		Result []Recipe
	}
	json.Unmarshal(raw, &recipe)

	return recipe.Result, nil
}

func TimeToCategoryID(date time.Time) string {
	month := date.Month()

	if 3 <= month && 5 >= month {
		return "52"
	} else if 6 <= month && 8 >= month {
		return "53"
	} else if 9 <= month && 11 >= month {
		return "54"
	} else if 12 <= month || 2 >= month {
		return "55"
	}

	return ""
}

func (r *Recipe) RecipeTemplate() *template.Buttons {
	btn := template.NewButtons()
	btn.Title = r.Title
	btn.SubTitle = r.RecipeCost
	btn.ImagePath = r.FoodImageURL
	btn.AddButtons(linebot.NewURIAction("レシピへ", r.RecipeURL))

	return btn
}
