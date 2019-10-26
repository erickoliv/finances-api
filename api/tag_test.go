package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mocket "github.com/Selvatico/go-mocket"
	"github.com/ericktm/olivsoft-golang-api/model"
	"github.com/stretchr/testify/assert"
)

var tagOne = map[string]interface{}{"uuid": "3272d69a-e38d-45c6-94fb-4a0dd9e69385", "name": "a test name", "description": "a test description", "owner": user.String()}
var tagTwo = map[string]interface{}{"uuid": "3272d69a-e38d-45c6-94fb-4a0dd9e61235", "name": "a two name", "description": "", "owner": user.String()}
var singleResult = []map[string]interface{}{tagOne}
var allRows = []map[string]interface{}{tagOne, tagTwo}
var count = []map[string]interface{}{{"count(*)": len(allRows)}}

func setupTagsDatabase() {
	mocket.Catcher.NewMock().WithArgs(tagOne["uuid"].(string), tagOne["owner"].(string)).WithReply(singleResult)

	mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "public"."tags"  WHERE "public"."tags"."deleted_at" IS NULL`).WithReply(count)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "public"."tags"  WHERE "public"."tags"."deleted_at" IS NULL`).WithReply(allRows)
}

func TestCreateTag(t *testing.T) {
	tag := model.Tag{
		Name:        "a simple tag",
		Description: "a tag description",
	}

	str, _ := json.Marshal(tag)

	req, _ := http.NewRequest("POST", "/api/tags", bytes.NewReader(str))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetTags(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/tags", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"total":2`)
}

func TestGetTag(t *testing.T) {
	url := fmt.Sprintf("/api/tags/%s", "3272d69a-e38d-45c6-94fb-4a0dd9e69385")
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), tagOne["name"])
}

func TestUpdateTag(t *testing.T) {
	tag := model.Tag{
		Name:        "updated name",
		Description: "a new description",
	}

	str, _ := json.Marshal(tag)
	url := fmt.Sprintf("/api/tags/%s", tagOne["uuid"].(string))

	req, _ := http.NewRequest("PUT", url, bytes.NewReader(str))
	req.Header.Add("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteTag(t *testing.T) {
	url := fmt.Sprintf("/api/tags/%s", tagOne["uuid"].(string))
	req, _ := http.NewRequest("DELETE", url, nil)
	w := httptest.NewRecorder()
	mocket.Catcher.NewMock().
		WithQuery(`UPDATE "public"."tags" SET "deleted_at"=?  WHERE "public"."tags"."deleted_at" IS NULL AND ((uuid = ? AND owner = ?))`).
		WithRowsNum(1)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
