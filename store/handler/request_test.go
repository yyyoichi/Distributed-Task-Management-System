package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yyyoichi/Distributed-Task-Management-System/pkg/document"
)

func TestHandler(t *testing.T) {
	tstore := NewStoreHandler() // テスト対象のTStoreインスタンスを生成

	// 正しいJSONデータを含むリクエストを作成
	validJSON := `{"task": "list"}`
	validReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(validJSON))
	validResp := httptest.NewRecorder()

	// ハンドラを呼び出して正常なレスポンスを得る
	tstore.CommandsHandler(validResp, validReq)

	// 正しいJSONデータの場合、200 OKのステータスコードを期待
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, validResp.Code)
	}

	// 不正なJSONデータを含むリクエストを作成
	invalidJSON := `{"wrong_field": "invalid task data"}`
	invalidReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
	invalidResp := httptest.NewRecorder()

	// ハンドラを呼び出して400 Bad Requestのレスポンスを得る
	tstore.CommandsHandler(invalidResp, invalidReq)

	// 不正なJSONデータの場合、400 Bad Requestのステータスコードを期待
	if invalidResp.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, but got %d", http.StatusBadRequest, invalidResp.Code)
	}
}

func TestHandler_differences(t *testing.T) {
	sh := NewStoreHandler()
	// init data
	sh.dc.Create("TaskA") // ID1 version1
	sh.dc.Create("TaskB") // ID2 version2
	sh.dc.Update(1, true) // version3

	validJSON := `{"version": 1}`
	validReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(validJSON))
	validResp := httptest.NewRecorder()

	// ハンドラを呼び出して正常なレスポンスを得る
	sh.DifferencesHandler(validResp, validReq)

	// 正しいJSONデータの場合、200 OKのステータスコードを期待
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, validResp.Code)
	}
	var data []document.TodoDataset
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
	sh.dc.Create("TaskA") // ID1 version1
	sh.dc.Create("TaskB") // ID2 version2
	// now
	// ID1 version1 TaskA no-complete
	// ID2 version2 TaskB no-complete

	// sh.dc.Update(1, true)
	validJSON := `{"Version":1,"todos":[
		{"id":1,"task":"TaskA","completed":true,"deleted":false,"version":1},
		{"id":2,"task":"TaskB","completed":true,"deleted":false,"version":1}
		]}`
	validReq, _ := http.NewRequest("POST", "/", bytes.NewBufferString(validJSON))
	validResp := httptest.NewRecorder()

	// ハンドラを呼び出して正常なレスポンスを得る
	sh.SynchronizeHandler(validResp, validReq)

	// now
	// ID1 version1 TaskA completed
	// ID2 version1 TaskB completed

	// 正しいJSONデータの場合、200 OKのステータスコードを期待
	if validResp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, validResp.Code)
	}
	todo := sh.dc.ByID[2]
	if !todo.Completed {
		t.Error("Expected completed is true, but got='false'")
	}

	sh.dc.Create("TaskC")
	// now
	// ID1 version1 TaskA completed
	// ID2 version1 TaskB completed
	// ID3 version2 TaskC no-complete

	todo, found := sh.dc.ByID[3]
	if !found {
		t.Error("Expected to find TODO[ID:3], but it was not found")
	}
	if todo.Version != 2 {
		t.Errorf("Expected Version is 2, but got='%d'", todo.Version)
	}
}
