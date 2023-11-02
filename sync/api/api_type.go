package api

import "github.com/yyyoichi/Distributed-Task-Management-System/pkg/store"

// 同期通信機。データ通信を行うために必要なメソッドを持つ。
type SyncerInterface interface {
	// 差分情報をリクエストする
	// [currentSyncVersion]以降の変更を取得する
	GetDifference(currentSyncVersion int) DiffResponse
	// 同期実行をリクエストする
	// - [currentSyncVersion]今回の同期バージョン。
	// - [todos]同期内容
	Synchronize(currentSyncVersion int, todo []store.TodoDateset) SynchronizeResponse
	Me() string
}

// 差分検知レスポンス
type DiffResponse struct {
	Err          error
	TodoDatasets []store.TodoDateset
}

// 同期実行レスポンス
type SynchronizeResponse struct {
	Err error
}
