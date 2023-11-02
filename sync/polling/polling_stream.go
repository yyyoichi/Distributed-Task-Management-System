package polling

import (
	"context"
	"yyyoichi/Distributed-Task-Management-System/pkg/store"
	"yyyoichi/Distributed-Task-Management-System/pkg/stream"
	"yyyoichi/Distributed-Task-Management-System/sync/api"
)

// 差分探知機
type differenceDetector struct {
	SyncerID int
	Get      func() api.DiffResponse
}

// 差分探知機をジェネレータとして送信する
func generateDifferenceDetector(cxt context.Context, detectors []differenceDetector) <-chan differenceDetector {
	return stream.Generator[differenceDetector](cxt, detectors...)
}

// 差分情報
type differences struct {
	SyncerID     int
	DiffResponse api.DiffResponse
}

// 差分探知機から差分情報をパイプする
func lineDetector2Differences(cxt context.Context, detectorCh <-chan differenceDetector, fn func(differenceDetector) differences) <-chan differences {
	return stream.FunIO[differenceDetector, differences](cxt, detectorCh, fn)
}

// 差分情報から複数のDatasetを送信する
func dLineDifferences2Dataset(cxt context.Context, inCh <-chan differences, fn func(d differences, produce func(store.TodoDateset))) <-chan store.TodoDateset {
	return stream.Demulti[differences, store.TodoDateset](cxt, inCh, fn)
}

// 同期実行機
type synchronizer struct {
	SyncerID int
	Exec     func() api.SyncResponse
}

// 同期実行機をジェネレータとして送信する
func generateSyncronizer(cxt context.Context, synchronizers []synchronizer) <-chan synchronizer {
	return stream.Generator[synchronizer](cxt, synchronizers...)
}

// 同期結果情報
type dones struct {
	SyncerID     int
	SyncResponse api.SyncResponse
}

// 同期実行機から同期結果情報をパイプする
func lineSynchronizer2Dones(cxt context.Context, synchronizerCh <-chan synchronizer, fn func(synchronizer) dones) <-chan dones {
	return stream.FunIO[synchronizer, dones](cxt, synchronizerCh, fn)
}
