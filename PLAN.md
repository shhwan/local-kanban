# Local Kanban Board - 開発計画書

## 1. プロジェクト概要

ローカル環境で動作するKanbanボードアプリケーション。
個人開発のタスク管理に特化し、AI/CLIからの操作にも対応する。

---

## 2. 要件定義

### 2.1 ボード構成

| 列 | 説明 |
|---|---|
| **TODO** | 未着手タスク |
| **DOING** | 作業中タスク |
| **DONE** | 完了タスク |

### 2.2 タスクの項目

| フィールド | 説明 | 必須 |
|---|---|---|
| タイトル | タスクの概要 | Yes |
| ラベル | FRONTEND / BACKEND / INFRA（将来追加可能） | Yes |
| 作業ログ | やったことの記録（複数エントリ、時系列） | No |
| 特記事項 | 特異事項・メモ（複数エントリ、追記のみ） | No |
| ステージ | TODO / DOING / DONE（将来追加可能） | Yes |
| 作成日時 | 自動付与 | Yes |
| 更新日時 | 自動付与 | Yes |

### 2.3 ラベルと色

| ラベル | 色 | 用途 |
|---|---|---|
| **FRONTEND** | 🔵 Blue (`#3B82F6`) | フロントエンド関連タスク |
| **BACKEND** | 🟢 Green (`#10B981`) | バックエンド関連タスク |
| **INFRA** | 🟠 Orange (`#F59E0B`) | インフラ関連タスク |

※ 初期は3種。将来追加可能（マスタテーブルで管理）

### 2.4 WIP制限

- DOING列: **ラベルごとに最大2タスク**まで
- 例: FRONTEND×2 + BACKEND×2 + INFRA×2 = 最大6タスクが同時DOING可能

### 2.5 操作方法

- ボタン操作で列を移動（TODO → DOING → DONE）
- ドラッグ&ドロップは不要
- AI/CLIからREST API経由でも操作可能

---

## 3. アーキテクチャ

### 3.1 構成図

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

### 3.2 サービス構成

| サービス | 役割 | 技術スタック |
|---|---|---|
| **frontend** | BFF。APIを呼びHTMLをレンダリング | Go, Echo, Templ, HTMX, TailwindCSS |
| **backend** | REST API (JSON) | Go, Echo, GORM |
| **db** | データ永続化 | PostgreSQL 16 |

### 3.3 設計パターン

- **BFF (Backend For Frontend)**: frontendサービスがbackend APIを呼び出し、HTMLに変換して返す
- **REST API**: backend は純粋なJSON APIとして設計。UI/CLI/AIすべてから利用可能
- **コンテナ分離**: 各サービスはDockerコンテナとして独立して動作

---

## 4. 技術スタック

### 4.1 共通

| 技術 | バージョン | 用途 |
|---|---|---|
| Go | 1.22+ | バックエンド・フロントエンド両方 |
| Docker | latest | コンテナ化 |
| Docker Compose | v2 | オーケストレーション |
| PostgreSQL | 16 | データベース |

### 4.2 Backend

| 技術 | 用途 |
|---|---|
| Echo | HTTPルーター・ミドルウェア |
| GORM | ORM (PostgreSQL接続) |

### 4.3 Frontend

| 技術 | 用途 |
|---|---|
| Echo | HTTPサーバー・ルーター |
| Templ | 型安全HTMLテンプレート (コンポーネント指向) |
| HTMX | SPAなしでリアクティブなUI |
| Tailwind CSS | ユーティリティベースCSS。ラベル色分けに活用 |

---

## 5. API設計 (Backend)

### 5.1 エンドポイント

| Method | Path | 説明 |
|---|---|---|
| `GET` | `/api/tasks` | タスク一覧取得 (クエリでステージ・ラベル絞り込み) |
| `GET` | `/api/tasks/:id` | タスク詳細取得 |
| `POST` | `/api/tasks` | タスク作成 |
| `PUT` | `/api/tasks/:id` | タスク更新 |
| `PATCH` | `/api/tasks/:id/stage` | ステージ変更 (TODO→DOING→DONE) |
| `DELETE` | `/api/tasks/:id` | タスク削除 |
| `POST` | `/api/tasks/:id/logs` | 作業ログ追加 |
| `POST` | `/api/tasks/:id/notes` | 特記事項追加 |
| `GET` | `/api/labels` | ラベル一覧取得 |
| `POST` | `/api/labels` | ラベル追加 |
| `GET` | `/api/stages` | ステータス一覧取得 |
| `POST` | `/api/stages` | ステータス追加 |

### 5.2 タスクJSON構造

```json
{
  "id": 1,
  "title": "ログイン画面のCSS調整",
  "label": {
    "id": 1,
    "name": "FRONTEND",
    "color": "#3B82F6"
  },
  "stage": {
    "id": 2,
    "name": "DOING",
    "position": 1
  },
  "notes": [
    {
      "id": 1,
      "content": "レスポンシブ対応も必要",
      "created_at": "2026-03-01T12:00:00Z"
    }
  ],
  "work_logs": [
    {
      "id": 1,
      "content": "ヘッダー部分のレイアウト修正完了",
      "created_at": "2026-03-02T10:30:00Z"
    }
  ],
  "created_at": "2026-03-01T09:00:00Z",
  "updated_at": "2026-03-02T10:30:00Z"
}
```

---

## 6. DB設計

### 6.1 テーブル

```sql
-- ラベルマスタ（将来追加可能）
CREATE TABLE labels (
    id    SERIAL      PRIMARY KEY,
    name  VARCHAR(20) NOT NULL UNIQUE,
    color VARCHAR(7)  NOT NULL          -- hex color (#3B82F6)
);

-- ステータスマスタ（将来追加可能）
CREATE TABLE stages (
    id       SERIAL      PRIMARY KEY,
    name     VARCHAR(20) NOT NULL UNIQUE,
    position INTEGER     NOT NULL        -- 列の表示順 (0: TODO, 1: DOING, 2: DONE)
);

-- タスク
CREATE TABLE tasks (
    id         SERIAL       PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    label_id   INTEGER      NOT NULL REFERENCES labels(id),
    stage_id  INTEGER      NOT NULL REFERENCES stages(id),
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- 作業ログ（追記のみ・イミュータブル）
CREATE TABLE work_logs (
    id         SERIAL    PRIMARY KEY,
    task_id    INTEGER   NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 特記事項（追記のみ・イミュータブル）
CREATE TABLE notes (
    id         SERIAL    PRIMARY KEY,
    task_id    INTEGER   NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    content    TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 初期データ
INSERT INTO labels (name, color) VALUES
    ('FRONTEND', '#3B82F6'),
    ('BACKEND',  '#10B981'),
    ('INFRA',    '#F59E0B');

INSERT INTO stages (name, position) VALUES
    ('TODO',  0),
    ('DOING', 1),
    ('DONE',  2);
```

---

## 7. ディレクトリ構成

```
local-kanban/
├── docker-compose.yml
├── PLAN.md
│
├── backend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── handler/          # HTTPハンドラ
│   │   ├── task.go
│   │   ├── label.go
│   │   └── stage.go
│   ├── model/            # データモデル (GORM)
│   │   ├── task.go
│   │   ├── work_log.go
│   │   ├── note.go
│   │   ├── label.go
│   │   └── stage.go
│   ├── repository/       # DB操作
│   │   ├── task.go
│   │   ├── label.go
│   │   └── stage.go
│   ├── service/          # ビジネスロジック (WIP制限等)
│   │   ├── task.go
│   │   ├── label.go
│   │   └── stage.go
│
├── frontend/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── main.go
│   ├── client/           # Backend APIクライアント
│   │   └── task.go
│   ├── handler/          # ページハンドラ
│   │   └── board.go
│   ├── templates/        # Templコンポーネント
│   │   ├── layout.templ
│   │   ├── board.templ
│   │   ├── task_card.templ
│   │   └── task_form.templ
│   └── static/           # CSS等
│       └── styles.css
│
└── db/
    └── init.sql          # DB初期化スクリプト
```

---

## 8. 開発フェーズ

### Phase 1: 基盤構築
- [ ] Docker Compose 環境構築 (Go + PostgreSQL)
- [ ] Backend: Echo セットアップ + DB接続
- [ ] Backend: マイグレーション実行
- [ ] Backend: CRUD API 実装

### Phase 2: フロントエンド
- [ ] Frontend: Echo + Templ セットアップ
- [ ] Frontend: Kanbanボード画面 (3列表示)
- [ ] Frontend: タスクカード (ラベル色分け)
- [ ] Frontend: ステータス変更ボタン
- [ ] Frontend: HTMX による部分更新

### Phase 3: 機能追加
- [ ] タスク作成フォーム
- [ ] 作業ログの追加・表示
- [ ] 特記事項の編集
- [ ] WIP制限 (ラベルごと2タスク) の実装とUI表示

### Phase 4: 仕上げ
- [ ] Tailwind CSS でデザイン調整
- [ ] エラーハンドリング
- [ ] ホットリロード (Air) 導入
