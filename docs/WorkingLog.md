# Working Log

## 2026-01-14

### プロジェクト開始: komikan-go

マンガ購入管理Nostr Botプロジェクトを開始。

#### プロジェクト設定

- **リポジトリ**: https://github.com/utenadev/komikan-go
- **ライセンス**: MIT
- **言語**: Go (Golang)
- **ターゲット**: ラズパイ3 (ARM64)

#### 実装した機能

1. **プロジェクト雛形**
   - ディレクトリ構造: cmd/, internal/, data/
   - go.mod 初期化
   - Taskfile.yml 作成（go-taskによるビルド自動化）

2. **Nostrライブラリ統合**
   - `nbd-wtf/go-nostr` 採用
   - リレー接続実装
   - イベント投稿実装
   - bot用鍵ペア生成ツール (cmd/genkey)

3. **データベース層 (BadgerDB)**
   - Pure Go KVS、ラズパイ最適
   - ISBNベースのマンガ保存
   - シリーズ索引の実装

4. **楽天ブックスAPIクライアント**
   - ISBN検索（要修正）
   - タイトル検索
   - ソート対応（発売日降順）

5. **マンガ管理機能**
   - 巻数抽出（正規表現）
   - 最新刊チェック
   - 新刊通知

6. **CLIツール**
   ```bash
   komikan-cli -isbn <ISBN>     # ISBNで登録
   komikan-cli -list             # 登録済み一覧
   komikan-cli -latest <タイトル> # 最新刊チェック
   ```

7. **Bot**
   - Nostr接続
   - 定期新刊チェック
   - タイムラインへの通知

#### 技術的選択

| 項目 | 選択 | 理由 |
|------|------|------|
| Nostrライブラリ | nbd-wtf/go-nostr | 2025年もアクティブ更新 |
| データベース | BadgerDB v4 | Pure Go、ARM対応 |
| タスクランナー | go-task | クロスプラットフォーム |
| 設定ファイル | YAML | 人間が読み書きしやすい |

#### ARM64ビルド

```bash
task build-arm64
```

生成物：
- bin/komikan-cli-arm64
- bin/komikan-bot-arm64

#### 既知の問題

1. **楽天API ISBN検索の問題**
   - isbnjanパラメータで指定ISBNと異なる本が返ってくる
   - タイトル検索は正常動作
   - 要調査/対応

#### Gitコミット

- `fbcad1e` Initial commit
- `0abd6b8` feat: integrate Nostr library and BadgerDB
- `9e4d740` feat: add volume extraction and latest volume check
- `4e4059b` feat: implement periodic new release checking in bot

#### 使用したAPI/サービス

- 楽天ブックスAPI (Application ID: 1056697313521518198)
- Nostrリレー:
  - wss://relay.damus.io
  - wss://nos.lol
  - wss://relay-jp.nostr.wirednet.jp

#### BotのNostr公開鍵

```
npub150lyta0e3vje8s8pv76dgq0fn32d86jhuuzu4f8x5ymxwlggjp7s3aukk2
```

#### 次回の課題

- 楽天API ISBN検索の問題解決
- 作品名からのマンガ登録機能
- ラズパイ3へのデプロイ
- テストデータの登録と動作確認
