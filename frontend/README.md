# Travelist Frontend

AI を活用した旅行プラン生成 Web アプリケーションのフロントエンド。

## 技術スタック

| カテゴリ | ライブラリ | バージョン |
|----------|-----------|-----------|
| UI | React | 19.x |
| 言語 | TypeScript | 5.9 |
| ビルド | Vite | 7.x |
| ルーティング | TanStack React Router | 1.x |
| データ取得 | TanStack React Query | 5.x |
| フォーム | React Hook Form | 7.x |
| PWA | vite-plugin-pwa | 1.x |

## セットアップ

### 前提条件

- Node.js 22+
- pnpm

### インストール

```bash
cd frontend
pnpm install
```

### 開発サーバー起動

```bash
pnpm dev
# http://localhost:5173 で起動
```

バックエンド API（Go、ポート 8080）が起動している必要があります。開発時は Vite のプロキシ設定により `/api` へのリクエストが `http://localhost:8080` に転送されます。

### Docker で起動

```bash
docker build -t travelist-frontend .
docker run -p 5173:5173 travelist-frontend
```

## pnpm スクリプト

| コマンド | 説明 |
|---------|------|
| `pnpm dev` | 開発サーバー起動（HMR 有効） |
| `pnpm build` | TypeScript 型チェック + 本番ビルド |
| `pnpm lint` | ESLint によるコード検査 |
| `pnpm preview` | 本番ビルドのプレビュー |
| `pnpm test` | Vitest によるテスト実行 |
| `pnpm test:watch` | Vitest のウォッチモード（変更検知で自動再実行） |

## ディレクトリ構成

```
src/
├── api/           # API クライアント（バックエンドとの通信）
├── components/    # React コンポーネント
│   ├── HomePage.tsx      # 生成された旅程の表示
│   ├── PlanForm.tsx      # 旅行プラン入力フォーム
│   └── RootLayout.tsx    # ルートレイアウト
├── routes/        # ルート定義
├── types/         # TypeScript 型定義
├── index.css      # グローバルスタイル
└── main.tsx       # エントリポイント
```

## アーキテクチャ

### データフロー

1. **PlanForm** でユーザーが旅行条件を入力（目的地、日数、予算、旅行スタイルなど）
2. React Query の mutation で `POST /api/plans` にリクエスト送信
3. バックエンド（Go + LLM）が旅程を生成
4. **HomePage** に生成された旅程を日別・時間帯別に表示

### API 連携

- エンドポイント: `/api/plans`（POST）、`/api/health`（GET）
- 開発時は Vite プロキシで `localhost:8080` に転送
- 本番では Go サーバーが静的ファイルも配信

### PWA

- Service Worker による自動更新
- API レスポンスの NetworkFirst キャッシュ（TTL: 5 分）
- オフライン対応の静的アセットキャッシュ

## 環境変数

フロントエンド自体に環境変数は不要です。バックエンド側で以下が必要です：

| 変数名 | 説明 |
|--------|------|
| `GOOGLE_PLACES_API_KEY` | Google Places API キー |
| `LLM_API_KEY` | LLM API キー（Anthropic Claude） |
