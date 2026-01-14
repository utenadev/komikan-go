# komikan-go

マンガ購入管理 Nostr Bot

## 概要

自分が購入しているマンガの情報を管理し、新刊情報などを Nostr タイムラインに通知するボット。

## 機能

- [ ] マンガの手動登録
- [ ] 最新巻の管理
- [ ] 楽天ブックスAPIによる新刊情報取得
- [ ] Nostrタイムラインへの通知

## セットアップ

```bash
# 依存関係のインストール
go mod download

# 設定ファイルの作成
cp config.example.yaml config.yaml
# config.yaml を編集してAPIキーなどを設定

# ビルド
go build -o bin/komikan ./cmd/bot

# 実行
./bin/komikan
```

## ラズパイ3での動作

```bash
# ARM向けビルド
GOARM=7 GOARCH=arm go build -o bin/komikan-arm ./cmd/bot
```

## ライセンス

MIT
