package api

import "yyyoichi/Distributed-Task-Management-System/pkg/store"

// データ通信を行うために必要なメソッドを持つ
type SyncerInterface interface {
	// 前回更新時[latestVersion]移行の変更を取得する
	GetDifference(latestVersion int) DiffResponse
	// 変更を同期する
	Sync(nextVersion int, todo []store.TodoDateset) SyncResponse
}

type DiffResponse struct {
	Err          error
	TodoDatasets []store.TodoDateset
}

type SyncResponse struct {
	Err error
}
