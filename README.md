# Travelist

AI を活用した旅行プラン生成 Web アプリケーション。

## 技術スタック

| レイヤー | 技術 |
|----------|------|
| Backend | Go 1.24 (net/http) |
| Frontend | React 19 + TypeScript 5.9 + Vite 7 |
| LLM | Anthropic Claude API |
| 場所検索 | Google Places API |

## API ドキュメント

- [api/API.md](api/API.md)

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

#### 方法 A: 1Password + direnv（推奨）

dotfiles の仕組みを使い、`.env` ファイルなしで秘密情報を管理する方法。

**前提条件:**
- [1Password](https://1password.com/) アカウント + デスクトップアプリ
- [1Password CLI (`op`)](https://developer.1password.com/docs/cli/) — `brew install 1password-cli`
- [direnv](https://direnv.net/) — `brew install direnv`
- `.zshrc` に `eval "$(direnv hook zsh)"` を追加済み
- 1Password アプリの **設定 → 開発者 → 「1Password CLI と連携」** を有効化

**手順:**

1. **1Password に API キーを登録**

   ```bash
   # 例: GOOGLE_PLACES_API_KEY を登録
   read -s -p "API Key: " KEY && \
   op item create \
     --vault Personal \
     --category "API Credential" \
     --title "Google Places" \
     "API key=$KEY" && \
   unset KEY
   ```

   | 環境変数名 | 1Password アイテム名 | フィールド名 |
   |-----------|---------------------|-------------|
   | `GOOGLE_PLACES_API_KEY` | `Google Places` | `API key` |
   | `LLM_API_KEY` | `Anthropic` | `API key` |

2. **`.envrc` を作成**

   ```bash
   cat <<'EOF' > .envrc
   command -v op >/dev/null || { echo "[direnv] op not found" >&2; exit 1; }
   : "${OP_VAULT:=Personal}"

   export PORT=8080
   export GOOGLE_PLACES_API_KEY="$(op read "op://${OP_VAULT}/Google Places/API key")"
   export LLM_API_KEY="$(op read "op://${OP_VAULT}/Anthropic/API key")"
   EOF
   ```

3. **direnv を許可**

   ```bash
   direnv allow
   ```

4. **動作確認**

   ```bash
   # ディレクトリに cd すると自動で環境変数がセットされる
   echo $GOOGLE_PLACES_API_KEY
   ```

> `.envrc` は `.gitignore` に含まれているため、リポジトリにはコミットされません。

#### 方法 B: `.env` ファイル（簡易）

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
| `LLM_PROVIDER` | - | `stub` | LLM プロバイダ (`stub`, `anthropic`, `gemini`) |
| `LLM_API_KEY` | **条件付き** | - | LLM API キー（`stub` 以外で必須） |

### API キーの取得方法

#### Google Places API キー

1. [Google Cloud Console](https://console.cloud.google.com/) にアクセス
2. プロジェクトを作成（または既存のプロジェクトを選択）
3. 「API とサービス」 → 「ライブラリ」 → **「Places API (New)」** を検索して有効化
4. 「API とサービス」 → 「認証情報」 → 「認証情報を作成」 → 「API キー」
5. 作成されたキーをコピーし `.env` に設定
6. **推奨**: 「キーを制限」から Places API のみに制限を設定

> 料金: 月 $200 分の無料枠あり。詳細は [料金ページ](https://developers.google.com/maps/billing-and-pricing/pricing) を参照。

#### LLM API キー

`LLM_PROVIDER` に応じて、対応する API キーを `LLM_API_KEY` に設定してください。

| プロバイダ | 取得先 | キー形式 |
|-----------|--------|---------|
| `anthropic` | [Anthropic Console](https://console.anthropic.com/) → API Keys → Create Key | `sk-ant-...` |
| `gemini` | [Google AI Studio](https://aistudio.google.com/) → Get API key → Create API key | `AIzaSy...` |
| `stub` | (不要) | - |

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
