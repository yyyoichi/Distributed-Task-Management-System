package polling

import (
	"yyyoichi/Distributed-Task-Management-System/pkg/store"
	"yyyoichi/Distributed-Task-Management-System/sync/api"
)

func NewSyncerStore(urls []string) SyncerStore {
	ss := SyncerStore{
		byID: make(map[int]api.SyncerInterface, len(urls)),
	}
	for id, url := range urls {
		ss.byID[id] = api.NewSyncer(url)
	}
	return ss
}

// 同期機構を採番(SyncerID)して保持する
type SyncerStore struct {
	byID map[int]api.SyncerInterface
}

// [syncerID]の同期機構を返す
func (ss *SyncerStore) GetAt(syncerID int) api.SyncerInterface {
	return ss.byID[syncerID]
}

// 差分探知機を返す
func (ss *SyncerStore) GetDifferenceDetectors(latestSyncVersion int) []DifferenceDetector {
	dd := make([]DifferenceDetector, len(ss.byID))
	for id, syncer := range ss.byID {
		dd = append(dd, DifferenceDetector{SyncerID: id, Get: func() api.DiffResponse {
			return syncer.GetDifference(latestSyncVersion)
		}})
	}
	return dd
}

// 同期実行機を返す
func (ss *SyncerStore) GetSynchronizers(nextVersion int, todos []store.TodoDateset) []Synchronizer {
	sn := make([]Synchronizer, len(ss.byID))
	for id, syncer := range ss.byID {
		sn = append(sn, Synchronizer{SyncerID: id, Exec: func() api.SyncResponse {
			return syncer.Sync(nextVersion, todos)
		}})
	}
	return sn
}
