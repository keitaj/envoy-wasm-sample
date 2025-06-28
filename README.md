# Envoy WASM Authentication Filter Sample

Go言語で実装したEnvoy Proxy用のWASM認証フィルタのサンプルです。

## 概要

このプロジェクトは、Envoy ProxyのWASMフィルタ機能を使用して、HTTPリクエストの認証を行うサンプル実装です。

### アーキテクチャ

```
[Client] → [Envoy Proxy + WASM Filter] → [Backend Service]
```

- **Envoy Proxy**: リバースプロキシとして動作し、WASMフィルタを実行
- **WASM Filter**: Go言語で実装された認証フィルタ
- **Backend Service**: シンプルなHTTPサーバー

### 機能

- Bearerトークンによる認証
- ヘルスチェックエンドポイントの認証スキップ
- 認証成功時のユーザー情報ヘッダー追加
- レスポンスへのフィルタ識別ヘッダー追加

## 必要な環境

- Docker
- Docker Compose
- Go 1.24以降

## セットアップ手順

1. リポジトリをクローン
```bash
git clone <repository-url>
cd envoy-wasm-sample
```

2. WASMフィルタをビルドしてサービスを起動
```bash
make run
```

これにより以下の処理が実行されます：
- Go言語のWASMフィルタをビルド
- Dockerイメージのビルド
- EnvoyとBackendサービスの起動

## 動作確認手順

### 個別のテスト実行

1. **ヘルスチェック（認証不要）**
```bash
curl -i http://localhost:10000/health
```

2. **認証なしでのアクセス（401エラー）**
```bash
curl -i http://localhost:10000/api/data
```

3. **無効なトークンでのアクセス（401エラー）**
```bash
curl -i -H "Authorization: Bearer invalid-token" http://localhost:10000/api/data
```

4. **有効なユーザートークンでのアクセス（成功）**
```bash
curl -i -H "Authorization: Bearer secret-token-123" http://localhost:10000/api/data
```

5. **有効な管理者トークンでのアクセス（成功）**
```bash
curl -i -H "Authorization: Bearer admin-token-456" http://localhost:10000/api/data
```

### 一括テスト実行

```bash
make test
```

## システム挙動の解説

### 認証フロー

1. **リクエスト受信**
   - クライアントからのHTTPリクエストがEnvoyに到着
   - WASMフィルタの`OnHttpRequestHeaders`が呼び出される

2. **パスチェック**
   - `/health`エンドポイントは認証をスキップ
   - その他のパスは認証処理を実行

3. **認証処理**
   - Authorizationヘッダーの存在確認
   - Bearerトークンの検証
   - 有効なトークン:
     - `secret-token-123`: userロール
     - `admin-token-456`: adminロール

4. **認証成功時**
   - `x-auth-user`ヘッダーにユーザー情報を追加
   - リクエストをバックエンドサービスに転送

5. **認証失敗時**
   - 401 Unauthorizedレスポンスを返却
   - JSONフォーマットでエラーメッセージを返す

### レスポンス処理

すべてのレスポンスに`x-wasm-filter: go-auth`ヘッダーが追加され、WASMフィルタが動作していることを確認できます。

### ログ出力

WASMフィルタは以下のログを出力します：
- プラグイン開始時: "plugin started"
- リクエスト処理時: パス情報
- 認証成功/失敗時: 詳細情報

これらのログはEnvoyのログとして確認できます：
```bash
docker logs envoy
```

## ディレクトリ構成

```
.
├── Makefile              # ビルド・実行・テスト用のコマンド定義
├── docker-compose.yaml   # Docker Compose設定
├── envoy.yaml           # Envoy設定ファイル
├── filter.wasm          # ビルドされたWASMファイル（自動生成）
├── backend/             # バックエンドサービス
│   ├── Dockerfile
│   └── main.go
└── wasm-filter/         # WASMフィルタのソースコード
    ├── go.mod
    ├── go.sum
    └── main.go
```

## トラブルシューティング

### WASMフィルタが読み込まれない場合

1. Envoyのログを確認
```bash
docker logs envoy
```

2. Go言語のバージョンを確認（1.24以降が必要）
```bash
go version
```

3. WASMファイルが正しくビルドされているか確認
```bash
ls -la filter.wasm
```

### クリーンアップ

サービスを停止し、ビルドファイルを削除：
```bash
make clean
```

## 技術詳細

- **Envoy**: v1.34-latest
- **Go**: 1.24（wasip1/wasm target）
- **proxy-wasm-go-sdk**: Envoy用WASMプラグイン開発SDK

WASMフィルタは`init()`関数でVMコンテキストを設定し、リクエストごとにHTTPコンテキストを作成して処理を行います。