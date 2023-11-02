package api

import "yyyoichi/Distributed-Task-Management-System/pkg/store"

// データ通信を行うために必要なメソッドを持つ
type SyncerInterface interface {
	// 前回更新時[latestVersion]移行の変更を取得する
	GetDifference(latestVersion int) DiffResponse
	// 変更を同期する
	// - [currentVersion]今回の同期バージョン。
	// - [todos]同期内容
	Sync(currentVersion int, todo []store.TodoDateset) SyncResponse
	Me() string
}

type DiffResponse struct {
	Err          error
	TodoDatasets []store.TodoDateset
}

type SyncResponse struct {
	Err error
}
