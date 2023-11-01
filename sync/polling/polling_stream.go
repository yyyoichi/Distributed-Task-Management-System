package polling

import (
	"context"
	"yyyoichi/Distributed-Task-Management-System/pkg/stream"
	"yyyoichi/Distributed-Task-Management-System/sync/api"
)

// 差分探知機
type DifferenceDetector struct {
	SyncerID int
	Get      func() api.DiffResponse
}

// 差分探知機をジェネレータとして送信する
func GenerateDifferenceDetector(cxt context.Context, detectors []DifferenceDetector) <-chan DifferenceDetector {
	return stream.Generator[DifferenceDetector](cxt, detectors...)
}

// 差分情報
type Differences struct {
	SyncerID     int
	DiffResponse api.DiffResponse
}

// 差分探知機から差分情報をパイプする
func LineDetector2Differences(cxt context.Context, detectorCh <-chan DifferenceDetector, fn func(DifferenceDetector) Differences) <-chan Differences {
	return stream.FunIO[DifferenceDetector, Differences](cxt, detectorCh, fn)
}

// 同期実行機
type Synchronizer struct {
	SyncerID int
	Exec     func() api.SyncResponse
}

// 同期実行機をジェネレータとして送信する
func GenerateSyncronizer(cxt context.Context, synchronizers []Synchronizer) <-chan Synchronizer {
	return stream.Generator[Synchronizer](cxt, synchronizers...)
}

// 同期結果情報
type Dones struct {
	SyncerID     int
	SyncResponse api.SyncResponse
}

// 同期実行機から同期結果情報をパイプする
func LineSynchronizer2Dones(cxt context.Context, synchronizerCh <-chan Synchronizer, fn func(Synchronizer) Dones) <-chan Dones {
	return stream.FunIO[Synchronizer, Dones](cxt, synchronizerCh, fn)
}
