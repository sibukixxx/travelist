# Travelist

AI を活用した旅行プラン生成 Web アプリケーション。

## 技術スタック

| レイヤー | 技術 |
|----------|------|
| Backend | Go 1.24 (net/http) |
| Frontend | React 19 + TypeScript 5.9 + Vite 7 |
| LLM | Anthropic Claude API |
| 場所検索 | Google Places API |

## セットアップ

### 前提条件

- Go 1.24+
- Node.js 22+
- npm
- Docker & Docker Compose（Docker で起動する場合）

### 1. リポジトリをクローン

```bash
git clone https://github.com/sibukixxx/travelist.git
cd travelist
```

### 2. 環境変数を設定

```bash
cp .env.example .env
```

`.env` を開き、各 API キーを設定してください。取得方法は下記「環境変数一覧」を参照。

### 3. 起動

#### Docker Compose（推奨）

```bash
make dev
# API: http://localhost:8080
# Frontend: http://localhost:5173
```

#### ローカル起動（Docker なし）

ターミナルを2つ開いて実行：

```bash
# ターミナル 1: API サーバー
make dev-api

# ターミナル 2: フロントエンド
make dev-frontend
```

## 環境変数一覧

| 変数名 | 必須 | デフォルト | 説明 |
|--------|------|-----------|------|
| `PORT` | - | `8080` | API サーバーのポート番号 |
| `STATIC_DIR` | - | (なし) | フロントエンド静的ファイルのパス（本番のみ） |
| `GOOGLE_PLACES_API_KEY` | **必須** | - | Google Places API キー |
| `LLM_API_KEY` | **必須** | - | Anthropic Claude API キー |

### API キーの取得方法

#### Google Places API キー

1. [Google Cloud Console](https://console.cloud.google.com/) にアクセス
2. プロジェクトを作成（または既存のプロジェクトを選択）
3. 「API とサービス」 → 「ライブラリ」 → **「Places API (New)」** を検索して有効化
4. 「API とサービス」 → 「認証情報」 → 「認証情報を作成」 → 「API キー」
5. 作成されたキーをコピーし `.env` に設定
6. **推奨**: 「キーを制限」から Places API のみに制限を設定

> 料金: 月 $200 分の無料枠あり。詳細は [料金ページ](https://developers.google.com/maps/billing-and-pricing/pricing) を参照。

#### Anthropic Claude API キー

1. [Anthropic Console](https://console.anthropic.com/) にアクセス
2. アカウントを作成 / ログイン
3. 「API Keys」 → 「Create Key」で新しいキーを発行
4. 発行されたキー（`sk-ant-...` 形式）をコピーし `.env` に設定

> 料金: 従量課金制。詳細は [料金ページ](https://www.anthropic.com/pricing) を参照。

## 開発コマンド

```bash
make dev              # Docker Compose で起動
make dev-api          # Go API サーバーのみ起動
make dev-frontend     # フロントエンド開発サーバーのみ起動

make test             # 全テスト実行
make test-api         # Go テストのみ
make test-frontend    # フロントエンドテストのみ

make lint             # 全 lint 実行
make lint-api         # Go vet のみ
make lint-frontend    # ESLint のみ

make build            # Docker イメージビルド
make clean            # ビルド成果物を削除
```

## プロジェクト構成

```
travelist/
├── api/                  # Go バックエンド
│   ├── cmd/server/       # エントリポイント
│   └── internal/
│       ├── domain/       # ドメインモデル
│       ├── handler/      # HTTP ハンドラ
│       ├── infra/        # 外部サービスクライアント
│       └── usecase/      # ユースケース
├── frontend/             # React フロントエンド
│   └── src/
│       ├── api/          # API クライアント
│       ├── components/   # React コンポーネント
│       ├── routes/       # ルート定義
│       └── types/        # TypeScript 型定義
├── .env.example          # 環境変数テンプレート
├── docker-compose.yml    # 開発用 Docker Compose
├── docker-compose.prod.yml # 本番用 Docker Compose
└── Makefile              # 開発コマンド
```
