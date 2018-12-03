package model_test

import (
	"testing"
	"time"

	"github.com/uzimaru0000/linebot/model"
)

const appID = ""

func TestGetCategory(t *testing.T) {
	categories, err := model.GetCategory(appID)
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("Success.\n%v", categories)
}

func TestGetRecipe(t *testing.T) {
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
