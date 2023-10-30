package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"yyyoichi/Distributed-Task-Management-System/store/store"
)

func TestHandler(t *testing.T) {
	tstore := store.NewStore() // テスト対象のTStoreインスタンスを生成

	// 正しいJSONデータを含むリクエストを作成
	validJSON := `{"task": "list"}`
	validReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(validJSON))
	validResp := httptest.NewRecorder()

	// ハンドラを呼び出して正常なレスポンスを得る
	tstore.Handler(validResp, validReq)

	// 正しいJSONデータの場合、200 OKのステータスコードを期待
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, validResp.Code)
	}

	// 不正なJSONデータを含むリクエストを作成
	invalidJSON := `{"wrong_field": "invalid task data"}`
	invalidReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
	invalidResp := httptest.NewRecorder()

	// ハンドラを呼び出して400 Bad Requestのレスポンスを得る
	tstore.Handler(invalidResp, invalidReq)

	// 不正なJSONデータの場合、400 Bad Requestのステータスコードを期待
	if invalidResp.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, invalidResp.Code)
	}
}
