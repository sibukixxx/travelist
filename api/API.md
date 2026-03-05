# Travelist API ドキュメント

最終更新: 2026-03-01

## 基本情報

- Base URL（ローカル開発）: `http://localhost:8080`
- コンテンツタイプ: `application/json`
- 認証: 現在は不要

## エンドポイント一覧

### 1. ヘルスチェック

- Method: `GET`
- Path: `/api/health`

レスポンス例（200）:

```json
{
  "status": "ok"
}
```

### 2. ユーザー登録

- Method: `POST`
- Path: `/api/users/register`

リクエスト:

```json
{
  "email": "alice@example.com"
}
```

レスポンス例（201）:

```json
{
  "id": "usr_1761998400000000000",
  "email": "alice@example.com",
  "created_at": "2026-03-01T12:00:00Z"
}
```

エラー:

- `400 Bad Request`
  - `{"error":"invalid request body"}`: JSON が不正
  - `{"error":"email is required"}`: email が空
  - `{"error":"invalid email format"}`: email 形式不正
- `409 Conflict`
  - `{"error":"email already registered"}`: 既存 email と重複
- `405 Method Not Allowed`
  - body: `method not allowed`

仕様メモ:

- email は `trim + lowercase` で正規化して保存されます。
- ユーザー保存先は現在 in-memory 実装です（プロセス再起動で消えます）。

## 実装はあるが未公開の API

以下はハンドラ/ユースケースの実装はありますが、`api/cmd/server/main.go` にルート登録されていないため、2026-03-01 時点では呼び出せません。

### 3. 旅行プラン生成（未配線）

- Method: `POST`
- Path: `/api/plans`

リクエスト:

```json
{
  "destination": "京都",
  "num_days": 2,
  "start_date": "2026-04-01",
  "preferences": {
    "interests": ["culture", "food"],
    "budget": "moderate",
    "travel_style": "balanced",
    "total_budget_yen": 30000
  },
  "constraint": {
    "max_walk_distance_m": 2000,
    "max_activities_day": 6,
    "earliest_start_time": "08:00",
    "latest_end_time": "21:00"
  }
}
```

レスポンス例（200）:

```json
{
  "itinerary": {
    "id": "itn_1761998400000000000",
    "title": "京都 2日間の旅",
    "destination": "京都",
    "start_date": "2026-04-01T00:00:00Z",
    "end_date": "2026-04-02T00:00:00Z",
    "days": [
      {
        "day_number": 1,
        "date": "2026-04-01T00:00:00Z",
        "activities": [
          {
            "order": 0,
            "place_id": "place-1",
            "start_time": "09:00",
            "end_time": "11:00",
            "duration_min": 120,
            "note": "朝一で訪問",
            "estimated_cost_yen": 500
          }
        ]
      }
    ],
    "created_at": "2026-03-01T12:00:00Z",
    "updated_at": "2026-03-01T12:00:00Z"
  },
  "violations": [],
  "budget_summary": {
    "total_cost_yen": 500,
    "daily_costs": [
      {
        "day_number": 1,
        "cost_yen": 500
      }
    ]
  }
}
```

エラー:

- `400 Bad Request`
  - `{"error":"invalid request body"}`
- `500 Internal Server Error`
  - `{"error":"invalid start_date: ..."}` など（外部依存失敗・保存失敗含む）

仕様メモ:

- `constraint.max_walk_distance_m` が `0` の場合、制約全体がデフォルト値に置き換わります。
- `violations` は制約違反があっても返却され、プラン生成自体は継続されます。
