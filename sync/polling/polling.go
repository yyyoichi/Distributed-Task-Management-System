package polling

import (
	"context"
	"fmt"
	"log"
	"yyyoichi/Distributed-Task-Management-System/pkg/store"
	"yyyoichi/Distributed-Task-Management-System/sync/api"
)

type PollingError struct {
	error
}

func NewPollingManager(urls []string) PollingManager {
	syncers := []api.SyncerInterface{}
	for _, url := range urls {
		syncers = append(syncers, api.NewSyncer(url))
	}
	pm := PollingManager{
		SyncerStore:     NewSyncerStore(syncers...),
		nextSyncVersion: 0,
	}
	return pm
}

type PollingManager struct {
	SyncerStore     SyncerStore
	nextSyncVersion int // 次回同期バージョン 同期実行後に更新される
}

func (pm *PollingManager) Polling(cxt context.Context) {
	c, cancel := context.WithCancelCause(cxt)
	defer cancel(nil)
	// 同期バージョン
	currentSyncVersion := pm.nextSyncVersion

	// TASK.1 同期するためのTodoDatasetを取得する
	// 		エラー発生時には同期を中断する
	// TASK.2 TodoDatasetの同期を実行する
	//		エラー発生時には同期を... // TODO

	// TASK.1 //
	// step.1 データノード(url)の数だけ差分検知を作成する //
	// 差分探知機チャネル
	detectorCh := pm.SyncerStore.getDifferenceDetectorCh(c, currentSyncVersion)
	// step.2 差分検知器で差分情報を取得する //
	// 差分情報チャネル
	differencesCh := lineDetector2Differences(c, detectorCh, func(dd differenceDetector) differences {
		resp := dd.Get() // 差分取得
		log.Printf("\tSyncer[%s]: GetDifferences Result is %d, Err(%s)", pm.SyncerStore.Who(dd.SyncerID), len(resp.TodoDatasets), resp.Err)
		return differences{SyncerID: dd.SyncerID, DiffResponse: resp}
	})

	// step.3 差分データセットのバージョンを書き換える //
	// 差分データセットチャネル
	todoCh := dLineDifferences2Dataset(c, differencesCh, func(d differences, produce func(store.TodoDateset)) {
		// エラー発生時ポーリングを中断する
		if d.DiffResponse.Err != nil {
			who := pm.SyncerStore.byID[d.SyncerID].Me()
			err := fmt.Errorf("PollingError[SyncerID:%d(%s)]: %s", d.SyncerID, who, d.DiffResponse.Err)
			cancel(PollingError{err})
			return
		}
		for _, todo := range d.DiffResponse.TodoDatasets {
			// 今回の同期バージョンに書き換える
			todo.Version = currentSyncVersion
			produce(todo)
		}
	})
	// step.4 差分データセット //
	todos := []store.TodoDateset{}
	for todo := range todoCh {
		todos = append(todos, todo)
	}

	// TASK.2 //
	// step.1 データノード(url)の数だけ同期実行機を作成する //
	// 同期実行機チャネル
	synchronizerCh := pm.SyncerStore.getSynchronizerCh(c, currentSyncVersion, todos)
	// step.2 同期実行機で同期結果情報を取得する //
	doneCh := lineSynchronizer2Dones(c, synchronizerCh, func(s synchronizer) dones {
		resp := s.Exec()
		return dones{SyncerID: s.SyncerID, SyncResponse: resp}
	})
	for range doneCh {
		// TODO error handling
	}

	// 同期終了。今回の同期バージョンを次回同期バージョンとする
	pm.nextSyncVersion = currentSyncVersion + 1
}
