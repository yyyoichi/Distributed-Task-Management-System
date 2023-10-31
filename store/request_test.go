package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	tstore := NewStoreHandler() // テスト対象のTStoreインスタンスを生成

	// 正しいJSONデータを含むリクエストを作成
	validJSON := `{"task": "list"}`
	validReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(validJSON))
	validResp := httptest.NewRecorder()

	// ハンドラを呼び出して正常なレスポンスを得る
	tstore.commandsHandler(validResp, validReq)

	// 正しいJSONデータの場合、200 OKのステータスコードを期待
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, validResp.Code)
	}

	// 不正なJSONデータを含むリクエストを作成
	invalidJSON := `{"wrong_field": "invalid task data"}`
	invalidReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
	invalidResp := httptest.NewRecorder()

	// ハンドラを呼び出して400 Bad Requestのレスポンスを得る
	tstore.commandsHandler(invalidResp, invalidReq)

	// 不正なJSONデータの場合、400 Bad Requestのステータスコードを期待
	if invalidResp.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, invalidResp.Code)
	}
}

func TestHandler_differences(t *testing.T) {
	sh := NewStoreHandler()
	// init data
	sh.s.Create("TaskA") // ID1 version1
	sh.s.Create("TaskB") // ID2 version2
	sh.s.Update(1, true) // version3

	validJSON := `{"version": 1}`
	validReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(validJSON))
	validResp := httptest.NewRecorder()

	// ハンドラを呼び出して正常なレスポンスを得る
	sh.differencesHandler(validResp, validReq)

	// 正しいJSONデータの場合、200 OKのステータスコードを期待
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, validResp.Code)
	}
	var data []todoDataset
	if err := json.Unmarshal(validResp.Body.Bytes(), &data); err != nil {
		t.Error(err)
	}
	if len(data) != 2 {
		t.Errorf("Expected len(data) is 2, but got='%d'", len(data))
	}
}

func TestHandler_sync(t *testing.T) {
	sh := NewStoreHandler()
	// init data
	sh.s.Create("TaskA") // ID1 version1
	sh.s.Create("TaskB") // ID2 version2

	// sh.s.Update(1, true)
	validJSON := `{"NextVersion":4,"todos":[{"id":2,"task":"TaskB","completed":true,"deleted":false,"version":3}]}`
	validReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(validJSON))
	validResp := httptest.NewRecorder()

	// ハンドラを呼び出して正常なレスポンスを得る
	sh.syncHandler(validResp, validReq)

	// 正しいJSONデータの場合、200 OKのステータスコードを期待
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, validResp.Code)
	}
	todo := sh.s.ByID[2]
	if !todo.Completed {
		t.Error("Expected completed is true, but got='false'")
	}

	sh.s.Create("TaskC")
	log.Println("Create New TODO")
	todo, found := sh.s.ByID[3]
	if !found {
		t.Error("Expected to find TODO[ID:3], but it was not found")
	}
	if todo.Version != 4 {
		t.Errorf("Expected Version is 4, but got='%d'", todo.Version)
	}
}
