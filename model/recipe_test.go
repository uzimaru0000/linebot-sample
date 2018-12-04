package model_test

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/uzimaru0000/linebot/model"
)

func getAppID() string {
	if err := godotenv.Load("../.env"); err != nil {
		return ""
	}

	return os.Getenv("APP_ID")
}

func TestGetCategory(t *testing.T) {
	appID := getAppID()
	_, err := model.GetCategory(appID)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

func TestGetRecipe(t *testing.T) {
	appID := getAppID()
	recipe, err := model.GetRecipe(appID, "55")
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("Success.\n%v", recipe)
}

func TestTimeToCategoryID(t *testing.T) {
	id := model.TimeToCategoryID(time.Now())

	if id != "55" {
		t.Fatalf("fail: %s", id)
	}
}

func TestMatchCategory(t *testing.T) {
	appID := getAppID()
	categories, err := model.GetCategory(appID)
	if err != nil {
		t.Fatalf("%v", err)
	}

	id := model.MatchCategory(categories, "肉")

	t.Log(id)
}
