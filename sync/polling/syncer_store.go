package polling

import (
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

// 差分探知機を返す
func (ss *SyncerStore) getDifferenceDetectors(latestSyncVersion int) []differenceDetector {
	dd := make([]differenceDetector, len(ss.byID))
	for id, syncer := range ss.byID {
		dd = append(dd, differenceDetector{SyncerID: id, Get: func() api.DiffResponse {
			return syncer.GetDifference(latestSyncVersion)
		}})
	}
	return dd
}

// 同期実行機を返す
func (ss *SyncerStore) getSynchronizers(nextVersion int, todos []store.TodoDateset) []synchronizer {
	sn := make([]synchronizer, len(ss.byID))
	for id, syncer := range ss.byID {
		sn = append(sn, synchronizer{SyncerID: id, Exec: func() api.SyncResponse {
			return syncer.Sync(nextVersion, todos)
		}})
	}
	return sn
}
