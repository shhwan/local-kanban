# Local Kanban Board

ローカル環境で動作するKanbanボード。
生成AIの作業を追うための、kanbanボード。
SkillやCLAUDE.mdにて、このボードでタスク管理をするような指示は必要。

## アーキテクチャ

```
┌───────────────────────────────────────────────────┐
│  Docker Compose                                   │
│                                                   │
│  ┌───────────────────┐    ┌────────────────────┐  │
│  │  frontend          │    │  backend (API)     │  │
│  │  Go + Templ + HTMX │───▶│  Go + Echo         │  │
│  │  BFF パターン       │    │  JSON REST API     │  │
│  │  port: 3000        │    │  port: 8080        │  │
│  └───────────────────┘    └─────────┬──────────┘  │
│                                     │             │
│                            ┌────────▼──────────┐  │
│                            │  PostgreSQL       │  │
│                            │  port: 5432       │  │
│                            └───────────────────┘  │
│                                                   │
└───────────────────────────────────────────────────┘
```

## 必要なもの

- Docker / Docker Compose v2
- Go 1.25（ローカルビルド時のみ）
- [templ](https://templ.guide/)（ローカルビルド時のみ）

## 起動

```bash
make up
```

ブラウザで http://localhost:3000 を開く。

## 停止

```bash
make down
```

## Makeコマンド

| コマンド | 説明 |
|---|---|
| `make up` | Docker Compose でビルド＆起動 |
| `make down` | 停止 |
| `make build` | backend / frontend をローカルビルド |
| `make generate` | templ テンプレートのコード生成 |
| `make test` | テスト実行 |
| `make vet` | 静的解析 |
| `make fmt` | コードフォーマット |
| `make clean` | ビルド成果物の削除 |

## API

### エンドポイント

| Method | Path | 説明 |
|---|---|---|
| `GET` | `/api/tasks` | タスク一覧取得 |
| `GET` | `/api/tasks/:id` | タスク詳細取得 |
| `POST` | `/api/tasks` | タスク作成 |
| `PUT` | `/api/tasks/:id` | タスク更新 |
| `PATCH` | `/api/tasks/:id/stage` | ステージ変更 |
| `DELETE` | `/api/tasks/:id` | タスク削除 |
| `POST` | `/api/tasks/:id/logs` | 作業ログ追加 |
| `POST` | `/api/tasks/:id/notes` | 特記事項追加 |
| `GET` | `/api/labels` | ラベル一覧取得 |
| `POST` | `/api/labels` | ラベル追加 |
| `GET` | `/api/stages` | ステージ一覧取得 |
| `POST` | `/api/stages` | ステージ追加 |

### 使用例

```bash
# タスク作成
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "READMEを書く", "label_id": 1, "stage_id": 1}'

# ステージ変更（TODO → DOING）
curl -X PATCH http://localhost:8080/api/tasks/1/stage \
  -H "Content-Type: application/json" \
  -d '{"stage_id": 2}'

# 作業ログ追加
curl -X POST http://localhost:8080/api/tasks/1/logs \
  -H "Content-Type: application/json" \
  -d '{"content": "下書き完了"}'
```

## WIP制限

DOING列はラベルごとに最大2タスクまで。制限を超えるとAPIが `409 Conflict` を返す。
