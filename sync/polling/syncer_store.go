package polling

import (
	"context"
	"yyyoichi/Distributed-Task-Management-System/pkg/store"
	"yyyoichi/Distributed-Task-Management-System/sync/api"
)

func NewSyncerStore(syncers ...api.SyncerInterface) SyncerStore {
	ss := SyncerStore{
		byID: make(map[int]api.SyncerInterface, len(syncers)),
	}
	for id, s := range syncers {
		ss.byID[id] = s
	}
	return ss
}

// 同期機構を採番(SyncerID)して保持する
type SyncerStore struct {
	byID map[int]api.SyncerInterface
}

// Syncer情報を返す
func (ss *SyncerStore) Who(syncerID int) string {
	return ss.byID[syncerID].Me()
}

// 差分探知機チャネルを返す
func (ss *SyncerStore) getDifferenceDetectorCh(cxt context.Context, currentSyncVersion int) <-chan differenceDetector {
	return generateDifferenceDetector(cxt, ss.byID, func(k int, v api.SyncerInterface) differenceDetector {
		detector := differenceDetector{
			SyncerID: k,
			Get: func() api.DiffResponse {
				resp := v.GetDifference(currentSyncVersion)
				return api.DiffResponse{TodoDatasets: resp.TodoDatasets, Err: resp.Err}
			},
		}
		return detector
	})
}

// 同期実行機を返す
func (ss *SyncerStore) getSynchronizerCh(cxt context.Context, currentSyncVersion int, todos []store.TodoDateset) <-chan synchronizer {
	return generateSyncronizer(cxt, ss.byID, func(k int, v api.SyncerInterface) synchronizer {
		synchronizer := synchronizer{
			SyncerID: k,
			Exec: func() api.SynchronizeResponse {
				return v.Synchronize(currentSyncVersion, todos)
			}}
		return synchronizer
	})
}
