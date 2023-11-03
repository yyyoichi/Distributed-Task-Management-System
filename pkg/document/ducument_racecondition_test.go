package document

import (
	"strconv"
	"sync"
	"testing"
)

func TestTDocument_ConcurrentAccess(t *testing.T) {
	tDocument := NewTDocument()

	// ゴルーチンの数
	numGoroutines := 100

	// ゴルーチンごとにTODOを作成する関数
	createTodo := func(index int) {
		tDocument.Create("Task " + strconv.Itoa(index))
	}

	// ゴルーチンの終了を待つためのWaitGroup
	var wg sync.WaitGroup

	// ゴルーチンを起動してTODOを作成する
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			createTodo(index)
		}(i)
	}

	// すべてのゴルーチンの終了を待つ
	wg.Wait()

	// すべてのTODOが正しく作成されたことを確認
	for i := 0; i < numGoroutines; i++ {
		task := "Task " + strconv.Itoa(i)
		id := tDocument.FindIDByTask(task)
		if id == -1 {
			t.Errorf("Expected to find TODO with task: '%s', but it was not found", task)
		}
	}
}

func TestTDocument_ConcurrentUpdate(t *testing.T) {
	tDocument := NewTDocument()

	// 初期状態でTODOを作成
	id := tDocument.Create("Task 1")

	// ゴルーチンの数
	numGoroutines := 100

	// ゴルーチンごとにTODOを更新する関数
	updateTodo := func() {
		tDocument.Update(id, true)
	}

	// ゴルーチンの終了を待つためのWaitGroup
	var wg sync.WaitGroup

	// ゴルーチンを起動してTODOを更新する
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			updateTodo()
		}()
	}

	// すべてのゴルーチンの終了を待つ
	wg.Wait()

	// TODOが正しく更新されたことを確認
	todo, found := tDocument.ByID[id]
	if !found || !todo.Completed {
		t.Errorf("Expected completed: true, but got: false")
	}
}

func TestTDocument_ConcurrentDelete(t *testing.T) {
	tDocument := NewTDocument()

	// 初期状態でTODOを作成
	id := tDocument.Create("Task 1")

	// ゴルーチンの数
	numGoroutines := 100

	// ゴルーチンごとにTODOを削除する関数
	deleteTodo := func() {
		tDocument.Delete(id)
	}

	// ゴルーチンの終了を待つためのWaitGroup
	var wg sync.WaitGroup

	// ゴルーチンを起動してTODOを削除する
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			deleteTodo()
		}()
	}

	// すべてのゴルーチンの終了を待つ
	wg.Wait()

	// TODOが正しく削除されたことを確認
	todo, found := tDocument.ByID[id]
	if !found {
		t.Errorf("Expected to find TODO[ID:1], but it was not found")
	}
	if !todo.Deleted {
		t.Errorf("Expected deleted: true, but got: false")
	}
}

// FindIDByTask は指定されたtaskを持つTODOのIDを返します。存在しない場合は-1を返します。
func (dc *TDocument) FindIDByTask(task string) int {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	for id, todo := range dc.ByID {
		if todo.Task == task {
			return id
		}
	}
	return -1
}
