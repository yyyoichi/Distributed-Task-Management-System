# 実行結果

CLI で実際に同期を試します。

## 前提

- データノードは 2 台あり、**1 つのコマンドはどちらか 1 つのデータノードに対して実行される**。
- そのため、**データの一貫性を保持するための機構**が必要であり、それが今回**同期機構**と呼んでいるものになる。
-

## Step.1

Step.1 では同期機構は存在せず、CLI コマンドはどちらか一方にのみ実行されるのみ。

<image src="./public/step1.png" />

画像のように、作成したデータが表示されない。

- ※`create`(作成)したときデータはデータノード 1 台目に保存され、次に`list`(一覧取得)してもデータノード 2 台目にデータを見に行くため、作成分のデータを参照できない。（このあと再び`list`すると 1 台目を参照するので作成分を見ることができる。）

ここから同期機構を作成することになる。

## Step.2

Step.1 時点ではデータの一貫性を保つことができなかった。

Step.2 では同期機構の実装を加えた。

### 変更点

同期機構の実装に加えてデータノードにも修正が必要だった。

#### データノードの追加実装

1. 各データノードは操作があるたびにキーバリューストアに一意のバージョンを振っていく。
2. 同期機構からの同期内容をキーバリューストアに反映する機能。

同期機構が変更内容を把握するために、データノードでバージョンを管理する必要が出てくる。

また同期機構がリクエストするためのエンドポイントを作成した。

##### データノードのバージョンの振り方例

- create, update, delete すると、バリュー（ToDo）に 1, 2, 3..と追加する。
- ※1 つのバリューについての連番ではなく、キーバリューストア全体で一つのバージョンをカウントアップする。

0. `初期状態`

```
<!-- datanode -->
- (empty)
```

1.  cli: `create TaskA`

```
<!-- datanode -->
- ID:1, Version:1 TaskA, no-complete
```

2.  cli: `create TaskB`

```
<!-- datanode -->
- ID:1, Version:1 TaskA, no-complete
- ID:2, Version:2 TaskB, no-complete
```

3.  cli: `update 1 complete`

```
<!-- datanode -->
- ID:1, Version:3 TaskA, completed
- ID:2, Version:2 TaskB, no-complete
```

#### 同期機構の実装

同期機構の動作は主に 2 つ。

1. 同期するバージョン以上のデータを各データノードから取得。
2. 取得されたデータ(差分)を、各データノードに渡して同期する。

また、差分データは**同期機構が**管理するバージョンを付与して、それをデータノードで管理しているバージョンに上書きする。

##### 同期機構とデータノードの動作例

便宜的に 2 つのデータノードを A,B と書きます。

0. `初期状態`

```
DatanodeA
- ID:1, Version:3 TaskA, completed
- ID:2, Version:2 TaskB, no-complete
DatanodeB
- (empty)
```

1. sync: `Get differences from version 1 onwards`

```
<!-- in sync machine -->
- ID:1, Version:3 TaskA, completed
- ID:2, Version:2 TaskB, no-complete
```

2. sync: `Stamps the sync machine version and sends it to all data nodes`

```
<!-- in sync machine -->
- ID:1, Version:1 TaskA, completed
- ID:2, Version:1 TaskB, no-complete

and send to datanodes
```

3. `Datanodes that accepted synchronization`

```
DatanodeA
- ID:1, Version:1 TaskA, completed
- ID:2, Version:1 TaskB, no-complete
DatanodeB
- ID:1, Version:1 TaskA, completed
- ID:2, Version:1 TaskB, no-complete
```

4. cli: `create TaskC`

```
DatanodeA
- ID:1, Version:1 TaskA, completed
- ID:2, Version:1 TaskB, no-complete
DatanodeB
- ID:1, Version:1 TaskA, completed
- ID:2, Version:1 TaskB, no-complete
- ID:3, Version:2 TaskC, no-complete
```

### 実行結果

終了時点で同期機構によってデータの同期ができるようになる。

<image src="./public/step2_success.png" />

左が CLI アプリケーションのコマンド、右が同期機構のログ。
緑がコマンド。

- `create`の後にデータを参照`list`できいて、
- 右側`02:04:36 create`の直後に、左側で`02:03:48`に同期されていることが分かる。
