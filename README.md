# Distributed-Task-Management-System

TODO リストを分散キーバリューストアに保存する CLI アプリケーションを作ってみる。

## 目的

- 分散システム興味があったので実装してみたい。
- 並行処理をより実践的に使ってみたい。
- 『Go 言語による分散サービス』を買う同期にする。

## 要件:

1. ToDo の追加: ユーザーは ToDo を追加できる。
2. ToDo の表示: ユーザーは ToDo リストを閲覧できる。
3. ToDo の更新: ユーザーは ToDo の状態（完了/未完了）を変更できる。
4. ToDo の削除: ユーザーは ToDo を削除できる。
5. 分散保存: ToDo データは複数のノードに分散して保存される。

## ユースケース

### ユーザー操作:

1. タスクの追加: ユーザーは CLI を使用して新しいタスクを追加できる。例えば、add Buy groceries と入力することで、"Buy groceries"というタスクが追加される。
2. タスクの表示: ユーザーは CLI コマンド（例: list）を使用してすべてのタスクを表示できる。これにより、ToDo リスト内のすべてのタスクが表示される。
3. タスクの更新: ユーザーはタスクの状態（完了または未完了）や説明を更新できる。例えば、update 1 complete と入力することで、ID が 1 のタスクが完了したことに更新される。
4. タスクの削除: ユーザーはタスクを削除できる。例えば、delete 1 と入力することで、ID が 1 のタスクが削除される。

### 分散キーバリューストアの操作:

1. キーと値の保存: ToDo リストの各タスクは一意の ID（キー）と説明（値）を持つ。ユーザーがタスクを追加すると、分散キーバリューストアにキーと値のペアが保存される。
2. タスクの取得: ユーザーが ToDo リストを表示すると、分散キーバリューストアからすべてのキーと値のペアが取得され、ToDo リストが表示される。
3. タスクの更新: ユーザーがタスクを更新すると、分散キーバリューストア内の対応するキーと値のペアが更新される。
4. タスクの削除: ユーザーがタスクを削除すると、分散キーバリューストアから対応するキーと値のペアが削除される。

## アーキテクチャ

### 1. **クライアント** (1 台):

- ユーザーコマンドを受け付け、リクエストハンドラに送信する CLI（Command Line Interface）。
- ユーザーが ToDo リストのタスクを追加、表示、更新、削除するためのインターフェースを提供。
- ロードバランサにユーザコマンドを送信する。

### 2. **ロードバランサー** (1 台):

- リクエストを受け取り、バックエンドのデータノードに負荷を均等に分散する役割を果たす。
- 負荷分散と高可用性を提供。

### 3. **データノード** (2 台):

- キーバリューペアを保存するサーバー。
- キーは ToDo リストの各タスクの一意の ID、値はタスクの説明（例: "TaskA"）。
- データノードは ToDo リストのタスクを保存し、必要に応じて検索や更新を行う。

### 4. **データ同期機構** (1 台):

- データノード間でデータを同期するための仕組み。変更があった場合に他のノードに通知し、データの整合性を保つ。
- 同期はポーリングによって行う。

## 同期機構詳細

データノード間でデータを同期するために、ポーリングを用いる。

なお、各ステップの実行結果については、<a href="./doc/README.md">./doc/README.md</a>に記載した。

## 使い方

### 必要なもの

- docker
- docker compose
- make

## 実行方法

1. `make up`
2. `make exec`してコンテナ内で`go run .`
3. <a href="#cliコマンド">cli コマンド</a>を参考にして実行する

## ログの確認方法

- `loadbrancer` -> `make logs-loadbrancer`
- `store1` -> `make logs-store1`
- `store2` -> `make logs-store2`

## 終了方法(コンテナ削除)

`make rm`

## CLI コマンド

### "create" コマンド:

- 引数: \<task\>
- 動作: 新しい TODO を作成し、その ID を返します。
- エラー: 引数が正しくない場合、ErrSyntaxInvalidArgs とエラーメッセージを返します。

### "list" コマンド:

- 引数: なし
- 動作: 未完了の TODO と完了した TODO をリストとして返します。
- エラー: このコマンドではエラーは発生しません。

### "update" コマンド:

- 引数: \<id\> \<status\> (status は "complete" または "open")
- 動作: 指定された ID の TODO のステータスを更新します。
- エラー: 引数が正しくない場合、ErrSyntaxInvalidArgs とエラーメッセージを返します。指定された ID が存在しない場合は、ErrNoDataFound とエラーメッセージを返します。

### "delete" コマンド:

- 引数: \<id\>
- 動作: 指定された ID の TODO を削除します。
- エラー: 引数が正しくない場合、ErrSyntaxInvalidArgs とエラーメッセージを返します。指定された ID が存在しない場合は、ErrNoDataFound とエラーメッセージを返します。

### "help" コマンド:

- 引数: なし
- 動作: サポートされているコマンドの説明を返します。
- エラー: このコマンドではエラーは発生しません。
