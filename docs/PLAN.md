# Project Plan

## 概要

komikan-go はマンガ購入情報を管理し、新刊情報をNostrタイムラインに通知するボットです。

## 実装済み機能 ✅

- [x] プロジェクト雛形
- [x] Nostrライブラリ統合 (nbd-wtf/go-nostr)
- [x] BadgerDBデータベース層
- [x] 楽天ブックスAPIクライアント
- [x] 巻数抽出機能（正規表現）
- [x] 最新刊チェック機能
- [x] CLIツール (`-isbn`, `-list`, `-latest`)
- [x] Bot定期新刊チェック
- [x] Nostr通知機能
- [x] ARM64クロスコンパイル
- [x] go-taskビルド自動化

## 既知の問題 🔧

### 優先度: 高

1. **楽天API ISBN検索の問題**
   - 現象: isbnjanパラメータで指定したISBNと異なる本が返ってくる
   - 影響: ISBNによるマンガ登録が正しく動作しない
   - 代替策: タイトル検索は正常動作
   - 対応方針:
     - ISBN検索を廃止し、タイトル+著者検索に移行
     - または、手動登録機能の強化

### 優先度: 中

2. **作品名からのマンガ登録機能不足**
   - 現状: ISBN登録のみ（しかし上記問題で動作不良）
   - 必要: タイトルから選択して登録する機能
   - 実装案: `komikan-cli -add <タイトル>` で検索結果を表示して選択

3. **シリーズ管理の不備**
   - 現状: AddToSeries関数があるが未使用
   - 必要: ISBN登録時に自動的にシリーズ索引を更新

## 今後の実装予定

### 短期 (v0.1.0)

1. **ISBN検索問題の解決**
   - [ ] タイトル検索に切り替え
   - [ ] 検索結果から選択して登録

2. **登録機能の改善**
   - [ ] タイトル検索インターフェース
   - [ ] 複数候補からの選択
   - [ ] シリーズ情報の自動保存

3. **基本動作の確認**
   - [ ] テストデータの登録
   - [ ] 新刊チェックの動作確認
   - [ ] Nostr通知の確認

### 中期 (v0.2.0)

1. **Bot機能の拡張**
   - [ ] 既存マンガの登録解除機能
   - [ ] 作者名によるフィルタリング
   - [ ] 複数作品の同時管理

2. **CLI改善**
   - [ ] 対話的な登録インターフェース
   - [ ] バーコードスキャン対応（ラズパイ+カメラ）
   - [ ] 設定ファイルのホットリロード

3. **通知の改善**
   - [ ] 新刊通知のテンプレート選択
   - [ ] 通知時間の設定
   - [ ] 既通知の管理（重複通知回避）

### 長期 (v1.0.0)

1. **Web UI**
   - [ ] Webダッシュボード
   - [ ] 登録/管理画面
   - [ ] Nostrタイムライン表示

2. **データソース拡張**
   - [ ] Amazon Product Advertising API
   - [ ] 国立国会図書館API
   - [ ] 出版社API

3. **高度な機能**
   - [ ] 読書履歴の記録
   - [ ] レンタル期限の通知
   - [ ] 作者の新刊チェック

## ラズパイ3デプロイ

### 準備

- [ ] ARM64バイナリのビルド (`task build-arm64`)
- [ ] config.yaml の作成
- [ ] 楽天APIキーの取得
- [ ] Nostr鍵ペアの生成

### デプロイ手順

```bash
# ラズパイ3にて
git clone https://github.com/utenadev/komikan-go.git
cd komikan-go
task build-arm64  # またはクロスコンパイル済みバイナリを転送
cp config.example.yaml config.yaml
# config.yaml を編集
./bin/komikan-bot-arm64 -config config.yaml
```

### サービス化

- [ ] systemd unit file の作成
- [ ] 自動起動の設定
- [ ] ログローテーション

## 技術的負債

- [ ] テストコードの追加
- [ ] エラーハンドリングの強化
- [ ] ログ出力の改善
- [ ] 設定バリデーション
- [ ] ヘルプメッセージの多言語対応

## 依存関係

### Goライブラリ

- `github.com/nbd-wtf/go-nostr` v0.52.3
- `github.com/dgraph-io/badger/v4` v4.9.0
- `gopkg.in/yaml.v3`

### 外部サービス

- 楽天ウェブサービス: https://webservice.rakuten.co.jp/
- Nostrリレー:
  - wss://relay.damus.io
  - wss://nos.lol
  - wss://relay-jp.nostr.wirednet.jp

## 開発ワークフロー

1. 機能実装
2. `task build` でビルド
3. ローカルテスト
4. `task build-arm64` でラズパイ用ビルド
5. git commit & push
6. GitHub Actions CI/CD（未実装）

## リリース計画

- **v0.1.0**: ISBN検索問題解決、基本機能完了
- **v0.2.0**: Bot機能拡張、UI改善
- **v0.5.0**: Web UI追加
- **v1.0.0**: 本格運用版

## 参考資料

- [楽天ブックスAPIドキュメント](https://webservice.rakuten.co.jp/documentation/books-book-search)
- [Nostrプロトコル仕様](https://github.com/nostr-protocol/nips)
- [BadgerDBドキュメント](https://dgraph.io/docs/badger/)
- [go-taskドキュメント](https://taskfile.dev/)
