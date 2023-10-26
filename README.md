# Distributed-Task-Management-System
TODOリストを分散キーバリューストアに保存するCLIアプリケーションを作ってみる。

## 要件:  
1. ToDoの追加: ユーザーはToDoを追加できる。
2. ToDoの表示: ユーザーはToDoリストを閲覧できる。
3. ToDoの更新: ユーザーはToDoの状態（完了/未完了）を変更できる。
4. ToDoの削除: ユーザーはToDoを削除できる。
5. 分散保存: ToDoデータは複数のノードに分散して保存される。

## ユースケース
### ユーザー操作:
1. タスクの追加: ユーザーはCLIを使用して新しいタスクを追加できる。例えば、add Buy groceriesと入力することで、"Buy groceries"というタスクが追加される。
2. タスクの表示: ユーザーはCLIコマンド（例: list）を使用してすべてのタスクを表示できる。これにより、ToDoリスト内のすべてのタスクが表示される。
3. タスクの更新: ユーザーはタスクの状態（完了または未完了）や説明を更新できる。例えば、update 1 completeと入力することで、IDが1のタスクが完了したことに更新される。
4. タスクの削除: ユーザーはタスクを削除できる。例えば、delete 1と入力することで、IDが1のタスクが削除される。

### 分散キーバリューストアの操作:
1. キーと値の保存: ToDoリストの各タスクは一意のID（キー）と説明（値）を持つ。ユーザーがタスクを追加すると、分散キーバリューストアにキーと値のペアが保存される。
2. タスクの取得: ユーザーがToDoリストを表示すると、分散キーバリューストアからすべてのキーと値のペアが取得され、ToDoリストが表示される。
3. タスクの更新: ユーザーがタスクを更新すると、分散キーバリューストア内の対応するキーと値のペアが更新される。
4. タスクの削除: ユーザーがタスクを削除すると、分散キーバリューストアから対応するキーと値のペアが削除される。

## アーキテクチャ
### 1. **クライアント** (1台):
- ユーザーコマンドを受け付け、リクエストハンドラに送信するCLI（Command Line Interface）。
- ユーザーがToDoリストのタスクを追加、表示、更新、削除するためのインターフェースを提供。

### 2. **リクエストハンドラ** (1台):
- クライアントからのリクエストを受け取り、適切なデータノードにリクエストを転送する役割を果たすサービス。
- ユーザーのToDoリストの操作（追加、表示、更新、削除）に対するリクエストを受け付け、データノードに転送。
  
### 3. **ロードバランサー** (1台):
- リクエストを受け取り、バックエンドのデータノードに負荷を均等に分散する役割を果たす。
- 負荷分散と高可用性を提供。

### 4. **データノード** (2台):
- キーバリューペアを保存するサーバー。
- キーはToDoリストの各タスクの一意のID、値はタスクの説明（例: "Buy groceries"）。
- データノードはToDoリストのタスクを保存し、必要に応じて検索や更新を行う。

### 5. **データ同期機構** (1台):
- データノード間でデータを同期するための仕組み。変更があった場合に他のノードに通知し、データの整合性を保つ。
- 同期はポーリングによって行う。

## 目的
- 分散システム興味があったので実装してみたい。
- 並行処理をより実践的に使ってみたい。
- 『Go言語による分散サービス』を買う同期にする。
