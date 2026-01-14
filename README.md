# komikan-go

マンガ購入管理 Nostr Bot

## 概要

自分が購入しているマンガの情報を管理し、楽天ブックスAPIで新刊情報を取得してNostrタイムラインに通知するボット。

## 機能

- [x] マンガの手動登録（ISBN）
- [x] 最新巻の管理
- [x] 楽天ブックスAPIによる新刊情報取得
- [x] Nostrタイムラインへの通知
- [x] タイトルからの最新刊チェック
- [x] 巻数抽出（正規表現）
- [x] 定期新刊チェック（bot）
- [x] ARM64対応（ラズパイ3/4/5）

## セットアップ

### 依存関係

```bash
# インストール
go mod download

# または go-task を使用
task deps
```

### 設定ファイル

```bash
# 設定ファイルの作成
cp config.example.yaml config.yaml
```

`config.yaml` を編集：

```yaml
nostr:
  secret_key: "nsec..."  # Nostr秘密鍵
  relays:
    - "wss://relay.damus.io"
    - "wss://nos.lol"
    - "wss://relay-jp.nostr.wirednet.jp"

rakuten:
  application_id: "YOUR_APP_ID"  # 楽天アプリケーションID

database:
  path: "data/komikan.db"

bot:
  check_interval: "1h"
  announce_new_releases: true
```

### Nostr鍵ペアの生成

```bash
# 鍵ペア生成ツール
go run ./cmd/genkey

# またはビルドして実行
go build -o bin/genkey ./cmd/genkey
./bin/genkey
```

出力された `nsec` を `config.yaml` に設定します。

### ビルド

```bash
# 通常ビルド
go build -o bin/komikan-cli ./cmd/cli
go build -o bin/komikan-bot ./cmd/bot

# または go-task を使用
task build
```

## 使い方

### CLIツール

```bash
# マンガを登録（ISBN）
RAKUTEN_APP_ID=your_app_id ./bin/komikan-cli -isbn 9784088847207

# 登録済みマンガの一覧
./bin/komikan-cli -list

# 最新刊をチェック
RAKUTEN_APP_ID=your_app_id ./bin/komikan-cli -latest ダンダダン
RAKUTEN_APP_ID=your_app_id ./bin/komikan-cli -latest ワンピース
```

### Botの実行

```bash
# 設定ファイルを指定
./bin/komikan-bot -config config.yaml

# バージョン表示
./bin/komikan-bot -version
```

Botは以下の動作を行います：
1. Nostrリレーに接続
2. 登録済みマンガの最新刊を定期チェック
3. 新刊が見つかったらNostrタイムラインに通知

### ラズパイ3での動作

```bash
# ARM64向けビルド
task build-arm64

# または手動ビルド
GOOS=linux GOARCH=arm64 go build -o bin/komikan-bot-arm64 ./cmd/bot

# ラズパイ3へ転送して実行
scp bin/komikan-bot-arm64 pi@raspberrypi:~/
ssh pi@raspberrypi ./komikan-bot-arm64 -config config.yaml
```

## プロジェクト構成

```
komikan-go/
├── cmd/
│   ├── bot/           # メインのbot
│   ├── cli/           # CLIツール
│   └── genkey/        # Nostr鍵ペア生成ツール
├── internal/
│   ├── api/           # 楽天ブックスAPI
│   ├── config/        # 設定管理
│   ├── db/            # BadgerDBデータベース
│   ├── manga/         # マンガ管理・新刊チェック
│   └── nostr/         # Nostrクライアント
├── data/              # データベースファイル
├── docs/              # ドキュメント
└── Taskfile.yml       # go-task タスク定義
```

## 技術スタック

| 項目 | 技術 |
|------|------|
| 言語 | Go (Golang) |
| Nostrライブラリ | nbd-wtf/go-nostr |
| データベース | BadgerDB v4 (Pure Go KVS) |
| 設定ファイル | YAML |
| タスクランナー | go-task |

## 既知の問題

- 楽天ブックスAPIのISBN検索が不安定（タイトル検索は正常動作）

## ドキュメント

- [WorkingLog.md](docs/WorkingLog.md) - 開発履歴
- [PLAN.md](docs/PLAN.md) - プロジェクト計画・課題

## ライセンス

MIT
