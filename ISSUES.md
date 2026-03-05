# 現状の課題

更新日: 2026-03-01

## 現在も未解決の課題

## 1. `/api/plans` がAPIサーバーに配線されていない
- 内容:
  フロントエンドは `POST /api/plans` を呼び出しているが、`api/cmd/server/main.go` でルート登録されていない。
- 影響:
  プラン生成機能がHTTP経由では利用できない（404）。
- 根拠:
  - フロント: `frontend/src/api/client.ts` は `/api/plans` を呼び出し
  - API起動: `api/cmd/server/main.go` は `/api/health` と `/api/users/register` のみ登録
- 対応案:
  - `PlanHandler` の依存実装（Places/LLM/ItineraryRepository）を用意し、`/api/plans` を `main.go` で登録する

## 2. ユーザ登録がインメモリ保存のみ（再起動で消える）
- 内容:
  `POST /api/users/register` は実装済みだが、保存先が in-memory repository。
- 影響:
  サーバ再起動で登録ユーザが消失し、本番利用できない。
- 根拠:
  - `api/internal/infra/repo/user_repo.go` が `map` ベースの実装
- 対応案:
  - PostgreSQL など永続ストレージに `users` テーブルを追加
  - `email` のユニーク制約をDB側でも担保

## 3. `make lint` の API 側が Go キャッシュ権限に依存する
- 内容:
  `test-api` は `GOCACHE=/tmp/go-build-cache` 指定だが、`lint-api` は未指定。
- 影響:
  環境によって `go vet` がキャッシュ権限エラーで失敗する。
- 根拠:
  - `Makefile` の `lint-api`: `cd api && go vet ./...`
- 対応案:
  - `lint-api` も `GOCACHE=/tmp/go-build-cache` を指定する

## 4. フロントエンド README の npm スクリプト表が実装と不一致
- 内容:
  `frontend/package.json` には `test` / `test:watch` があるが、`frontend/README.md` の表に未記載。
- 影響:
  テストコマンドの発見性が低い。
- 対応案:
  - README の npm スクリプト表を `package.json` と同期する

## 使いやすさを向上させる次の課題（2026-03-01 追加）

## 5. APIエラーがそのまま表示され、ユーザーに意味が伝わりにくい
- 内容:
  フロントの `fetchJSON` は `API error 409: {"error":"..."}` 形式で例外化し、その文字列をフォームで直接表示している。
- 影響:
  「何を直せばよいか」が分かりづらく、離脱につながる。
- 根拠:
  - `frontend/src/api/client.ts` の `throw new Error(...)`
  - `frontend/src/components/PlanForm.tsx` / `UserRegistrationForm.tsx` の `mutation.error.message` 直接表示
- 対応案:
  - APIエラーを `status` / `code` / `message` に正規化する共通エラーパーサを導入
  - `409`（重複メール）や `400`（入力不備）に対して、ユーザー向け文言を個別表示

## 6. プラン生成フォームの入力支援が弱く、誤入力を防ぎにくい
- 内容:
  `開始日` は過去日付の制約がなく、`興味・関心` はカンマ区切り自由入力のみ。
- 影響:
  入力ミスや解釈ブレが起きやすく、生成品質のばらつきが大きい。
- 根拠:
  - `frontend/src/components/PlanForm.tsx` の `start_date` に `min` 制約なし
  - `interests` が単一テキスト入力（タグUIなし）
- 対応案:
  - `start_date` に `min=today` を設定し、過去日を選べないようにする
  - `interests` をタグ入力（候補サジェスト付き）に変更する

## 7. 生成結果の再利用導線（保存・共有・再編集）がなく、試行錯誤しづらい
- 内容:
  生成結果は画面表示のみで、保存・共有・再編集の導線がない。
- 影響:
  比較検討や家族/友人との共有が難しく、実利用時の利便性が低い。
- 根拠:
  - `frontend/src/components/HomePage.tsx` は表示のみ（アクションボタンなし）
  - API側も `GET /api/itineraries` 系未提供（`api/API.md` 記載なし）
- 対応案:
  - MVPとして「JSONダウンロード」「再入力して再生成」ボタンを追加
  - 永続化後に `一覧→詳細→再編集` 導線を段階導入

## 8. フォーム/結果のE2E観点テストが不足し、UX劣化を検知しにくい
- 内容:
  フロントのテストはレンダリング中心で、送信エラー時表示や再送フローの検証が薄い。
- 影響:
  小さなUI変更で使い勝手が悪化しても、CIで検知できない。
- 根拠:
  - `frontend/src/components/HomePage.test.tsx` は初期表示確認が中心
- 対応案:
  - `PlanForm` / `UserRegistrationForm` に対して
    - バリデーションエラー表示
    - API失敗時メッセージ
    - 成功後の状態遷移
    をテスト追加する

## 追加したほうがいい機能（提案）

## A. メール確認付き登録（ダブルオプトイン）
- 目的:
  不正登録・タイプミスを減らし、通知の到達率を上げる。
- 最小実装:
  - `users` テーブルに `email_verified_at`, `verification_token`, `token_expires_at` を追加
  - `/api/users/verify?token=...` を追加

## B. パスワードレスログイン（マジックリンク）
- 目的:
  パスワード管理コストを下げつつ認証を提供する。
- 最小実装:
  - `/api/auth/magic-link/request`, `/api/auth/magic-link/verify` を追加
  - 有効期限付きワンタイムトークンを発行

## C. 旅程の保存・再編集機能
- 目的:
  生成結果を後で見直し、再生成・共有しやすくする。
- 最小実装:
  - `itineraries` の永続化
  - `GET /api/itineraries`, `GET /api/itineraries/:id`, `PATCH /api/itineraries/:id`

## D. 通知配信（登録完了・旅程完成）
- 目的:
  ユーザの再訪率を上げる。
- 最小実装:
  - 登録完了メール
  - 旅程生成完了メール
  - 配信失敗時の再試行キュー

## メール機能の提案（無料で始められるサービス）

## 1. Resend
- 向いている用途:
  開発者向けトランザクションメール（登録確認、マジックリンク）
- 無料枠の目安:
  月3,000通 / 日100通
- 特徴:
  APIがシンプル、実装が速い

## 2. Brevo
- 向いている用途:
  通知 + 将来のマーケ配信も見据える場合
- 無料枠の目安:
  日300通
- 特徴:
  マーケ機能・テンプレート管理を同じ基盤で扱える

## 3. Mailgun
- 向いている用途:
  API中心の配信基盤、Webhook連携
- 無料枠の目安:
  日100通
- 特徴:
  ログ/イベント解析やWebhookが扱いやすい

## 4. Amazon SES
- 向いている用途:
  将来的に大量配信・AWS統合を前提にする場合
- 無料枠の目安:
  AWSアカウント作成時期・プランにより異なる（新規アカウント向けにはAWSクレジット方式）
- 特徴:
  従量課金が安価で、スケール時のコスト効率が高い

## 解消済み（旧ISSUESからの更新）

## A. `generate_plan` のユニットテスト未整備は解消済み
- 根拠:
  - `api/internal/usecase/generate_plan_test.go` が存在

## B. `npm test` スクリプト未定義問題は解消済み
- 根拠:
  - `frontend/package.json` に `"test": "vitest run"` が存在
