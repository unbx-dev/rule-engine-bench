# ルールエンジン ベンチマーク結果

**計測日:** 2026-07-03 〜 2026-07-07（最新実行: 2026-07-07）  
**PR:** #1 — RuleEngine検証用PR (benchmark → main)  
**エンジン:** Unbx Code Quarantine (`unbx-dev/unbx-cli@v0.0.10`)

---

## 計測方法

CIエンジンが GitHub PR のインラインレビューコメントとして投稿した違反検知結果を集計した。  
グランドトゥルース（正解）は、各ソースファイルを 30 ルール（言語ごとに 10 ルール）に照らして手動監査して作成した。

**指標の定義:**
- **TP（真陽性）:** エンジンが正しく検知した違反
- **FP（偽陽性）:** エンジンが誤ってフラグを立てた箇所（実際には違反でない）
- **FN（偽陰性）:** 実際には違反だがエンジンが見逃した箇所
- **適合率 (Precision)** = TP / (TP + FP)
- **再現率 (Recall)** = TP / (TP + FN)
- **F1スコア** = 2 × 適合率 × 再現率 / (適合率 + 再現率)

---

## スキャン実行サマリー

| 実行 | 日時 (UTC) | スキャン対象ファイル | 正解数 | TP | FP | FN | 適合率 | 再現率 | F1 |
|:---:|:----------:|:------------------|:------:|:--:|:--:|:--:|:------:|:------:|:--:|
| 1 | 2026-07-03 02:29 | `bad_example.go` | 15 | 11 | 0 | 4 | 100.0% | 73.3% | 84.6% |
| 2 | 2026-07-03 02:36 | `bad_example.go`, `bad_example.py` | 29 | 14 | 0 | 15 | 100.0% | 48.3% | 65.1% |
| 3 | 2026-07-03 12:31 | `bad_example.go` | 15 | 5 | 0 | 10 | 100.0% | 33.3% | 50.0% |
| 4 | 2026-07-03 13:55 | `bad_example.go`, `bad_example.py`, `bad_example.ts` | 50 | 29 | 0 | 21 | 100.0% | 58.0% | 73.4% |
| 5 | 2026-07-03 14:34 | `go/*`, `python/*`, `typescript/*` | 46 | 25 | 2 | 21 | 92.6% | 54.3% | 68.4% |
| 6 | 2026-07-04 08:30 | `go/*`, `python/*`, `typescript/*` | 46 | 22 | 6 | 24 | 78.6% | 47.8% | 59.4% |
| 7 | 2026-07-05 03:48 | `go/*`, `python/*`, `typescript/*` | 46 | 23 | 1 | 23 | 95.8% | 50.0% | 65.7% |
| 8 | 2026-07-06 06:22 | `go/*`, `python/*`, `typescript/*` | 46 | 24 | 4 | 22 | 85.7% | 52.2% | 64.9% |
| 9 | 2026-07-07 07:10 | `go/*`, `python/*`, `typescript/*` | 46 | 17 | 0 | 29 | 100.0% | 37.0% | 54.0% |

> **ファイル構成の変化:** 実行 1〜4 はリポジトリルートのフラットなテストファイルをスキャン。実行 5（コミット `84fbcd4` トリガー）から言語別サブディレクトリ構成に移行。実行 6〜9 も同じ構成を継続。`typescript/components/UserCard.tsx` は実行 5〜9 でスキャン対象外。その 4 件の違反は正解数 46 に含まれ FN としてカウントする。

---

## グランドトゥルース：ファイルごとの違反内容

### Go — `bad_example.go`（フラット構成、実行 1〜4）

| # | ルール | 件数 |
|---|--------|:----:|
| 1 | `main` 以外での `panic` 使用 | 2（service、handler 各 1） |
| 2 | `log.Fatal` 使用 | 3（service 1、main 2） |
| 3 | 本番コードでの `fmt.Println` 使用 | 2（service、handler 各 1） |
| 4 | handler/service 層での `context.Background` 使用 | 2（service、handler 各 1） |
| 5 | service 層での `net/http` インポート | 1 |
| 6 | domain 層での `database/sql` インポート | 1 |
| 7 | handler 層での repository 直接インポート | 1 |
| 8 | Repository が sqlc 型をそのまま返す | 1 |
| 9 | Domain model に `json`/`db` struct タグを付与 | 1 |
| 10 | `Fix` フィールドを `Replacement` でラップしていない | 1 |
| **合計** | | **15** |

### Go — 再編成後のファイル群（実行 5）

| ファイル | 違反内容 | 件数 |
|---------|---------|:----:|
| `go/domain/user.go` | `database/sql` インポート；json/db struct タグ | 2 |
| `go/handler/user.go` | repository インポート；`context.Background`；`panic`；`fmt.Println` | 4 |
| `go/repository/user.go` | sqlc 型をそのまま返す；`Fix` を `Replacement` で返さない | 2 |
| `go/service/user.go` | `net/http` インポート；`context.Background`；`log.Fatal`；`panic`；`fmt.Println` | 5 |
| **合計** | | **13** |

### Python — `bad_example.py` / `python/bad_example.py`（実行 2、4、5）

| # | ルール | 件数 |
|---|--------|:----:|
| 1 | `logger` の代わりに `print` を使用 | 3（9行目、15行目、77行目） |
| 2 | `eval` 使用 | 1（14行目） |
| 3 | `exec` 使用 | 1（26行目） |
| 4 | `os.system` 使用 | 1（21行目） |
| 5 | `pickle` 使用 | 2（dump 32行目、load 38行目） |
| 6 | `requests.get` 直接呼び出し | 1（43行目） |
| 7 | `requests.post` 直接呼び出し | 1（49行目） |
| 8 | `datetime.now` 直接呼び出し | 1（55行目） |
| 9 | bare `except:`（例外型の指定なし） | 1（63行目） |
| 10 | `except` ブロック内で `pass` | 1（64行目） |
| 11 | service 層での `db.commit()` | 1（76行目） |
| **合計** | | **14** |

### TypeScript — `bad_example.ts`（フラット構成、実行 4）

| # | ルール | 件数 |
|---|--------|:----:|
| 1 | `@ts-ignore` | 2（1行目、25行目） |
| 2 | `any` 型アノテーション | 6（各所） |
| 3 | `as any` キャスト | 2 |
| 4 | 本番コードでの `console.log` | 4（9、29、56、65行目） |
| 5 | `fetch` 直接呼び出し | 2 |
| 6 | domain 層での `Date.now` 使用 | 1 |
| 7 | domain 層への `react` インポート | 1 |
| 8 | UI コンポーネントへの DB クライアントインポート | 1 |
| 9 | Client Component への `server-only` インポート | 1 |
| 10 | `dangerouslySetInnerHTML` 使用 | 1 |
| **合計** | | **21** |

### TypeScript — 再編成後のファイル群（実行 5）

| ファイル | 違反内容 | 件数 |
|---------|---------|:----:|
| `typescript/bad_example.ts` | `@ts-ignore`×2；`any`×5；`as any`×2；`console.log`×2；`fetch`×2 | 13 |
| `typescript/domain/user.ts` | `react` インポート；`Date.now` | 2 |
| `typescript/components/UserCard.tsx` | DB クライアントインポート；`server-only`；`dangerouslySetInnerHTML`；`any` | 4 |
| **合計** | | **19** |

---

## 実行別詳細分析

### 実行 1 — 2026-07-03 02:29 UTC

**スキャン対象:** `bad_example.go`  
**エンジン報告:** "Detected 11 architecture policy violation(s)"

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `main` 以外での `panic`（×2） | ✅ | TP |
| `log.Fatal`（×3） | ✅ | TP |
| 本番コードでの `fmt.Println`（×2） | ✅ | TP |
| handler/service での `context.Background`（×2） | ✅ | TP |
| service 層の `net/http` インポート | ✅ | TP |
| domain 層の `database/sql` インポート | ✅ | TP |
| handler が repository を直接インポート | ❌ | FN |
| Repository が sqlc 型をそのまま返す | ❌ | FN |
| Domain model に json/db タグ | ❌ | FN |
| `Fix` を `Replacement` で返さない | ❌ | FN |

**TP=11, FP=0, FN=4 → 適合率=100.0%, 再現率=73.3%, F1=84.6%**

---

### 実行 2 — 2026-07-03 02:36 UTC

**スキャン対象:** `bad_example.go`、`bad_example.py`  
**エンジン報告:** "Detected 14 architecture policy violation(s)"

Go 検知結果：実行 1 と同一（TP=11、FP=0、FN=4）。

| Python ルール | 検知 | 判定 |
|--------------|:----:|------|
| `print`（×3） | ✅ | TP |
| `eval` | ❌ | FN |
| `exec` | ❌ | FN |
| `os.system` | ❌ | FN |
| `pickle`（×2） | ❌ | FN |
| `requests.get` | ❌ | FN |
| `requests.post` | ❌ | FN |
| `datetime.now` | ❌ | FN |
| bare `except` | ❌ | FN |
| `except` 内で `pass` | ❌ | FN |
| service 層で `db.commit` | ❌ | FN |

**TP=14, FP=0, FN=15 → 適合率=100.0%, 再現率=48.3%, F1=65.1%**

---

### 実行 3 — 2026-07-03 12:31 UTC

**スキャン対象:** `bad_example.go`（部分スキャン）  
**エンジン報告:** "Detected 5 architecture policy violation(s)"

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `log.Fatal`（×3） | ✅ | TP |
| `fmt.Println`（×2） | ✅ | TP |
| `panic`（×2） | ❌ | FN |
| `context.Background`（×2） | ❌ | FN |
| service 層の `net/http` | ❌ | FN |
| domain 層の `database/sql` | ❌ | FN |
| handler が repository を直接インポート | ❌ | FN |
| Repository が sqlc 型を返す | ❌ | FN |
| Domain model に json/db タグ | ❌ | FN |
| `Fix` を `Replacement` で返さない | ❌ | FN |

**TP=5, FP=0, FN=10 → 適合率=100.0%, 再現率=33.3%, F1=50.0%**

> ロギング関連の 2 ルールしか検知されておらず、ルール設定が一時的に変更または縮退していた可能性がある。

---

### 実行 4 — 2026-07-03 13:55 UTC

**スキャン対象:** `bad_example.go`、`bad_example.py`、`bad_example.ts`  
**エンジン報告:** "Detected 29 architecture policy violation(s)"

Go（bad_example.go）：実行 1 と同一（TP=11、FP=0、FN=4）。

Python（bad_example.py）：

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `print`（×3） | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `requests.get/post`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `datetime.now` | ✅ | TP |
| bare `except` | ❌ | FN |
| `except` 内で `pass` | ❌ | FN |
| service 層で `db.commit` | ❌ | FN |

TypeScript（bad_example.ts）：

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `console.log`（×4） | ✅ | TP |
| `fetch` 直接呼び出し（×2） | ✅ | TP |
| domain 層での `Date.now` | ✅ | TP |
| domain 層への `react` インポート | ✅ | TP |
| Client Component への `server-only` | ✅ | TP |
| `@ts-ignore`（×2） | ❌ | FN |
| `any` 型（×6） | ❌ | FN |
| `as any`（×2） | ❌ | FN |
| UI への DB クライアントインポート | ❌ | FN |
| `dangerouslySetInnerHTML` | ❌ | FN |

**TP=29, FP=0, FN=21 → 適合率=100.0%, 再現率=58.0%, F1=73.4%**

---

### 実行 5 — 2026-07-03 14:34 UTC

**スキャン対象:** `go/domain/user.go`、`go/handler/user.go`、`go/repository/user.go`、`go/service/user.go`、`python/bad_example.py`、`typescript/bad_example.ts`、`typescript/domain/user.ts`  
**エンジン報告:** "Detected 27 architecture policy violation(s)"

> ※ `typescript/components/UserCard.tsx` は今回のスキャン対象に含まれていない。

#### Go（ファイル別）

| ファイル | 検知内容 | TP | FP | FN | 備考 |
|---------|---------|:--:|:--:|:--:|------|
| `go/domain/user.go` | `database/sql` インポート | 1 | 0 | 1 | json/db タグを見逃し |
| `go/handler/user.go` | `panic`、`fmt.Println`、`context.Background`、repository インポート、`net/http` インポート | 4 | 1 | 0 | `net/http` は handler 層では正当 → **FP** |
| `go/repository/user.go` | `database/sql` インポート | 0 | 1 | 2 | `database/sql` は repository 層では正当 → **FP**；sqlc 返却型・Fix フォーマットを見逃し |
| `go/service/user.go` | `net/http`、`context.Background`、`log.Fatal`、`panic`、`fmt.Println` | 5 | 0 | 0 | 完全検知 |

#### Python（`python/bad_example.py`）

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `print`（×3） | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `requests.get/post`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `datetime.now` | ✅ | TP |
| bare `except` | ❌ | FN |
| `except` 内で `pass` | ❌ | FN |
| service 層で `db.commit` | ❌ | FN |

#### TypeScript

| ファイル | 検知内容 | TP | FP | FN |
|---------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `console.log`×2、`fetch`×2 | 4 | 0 | 9（`@ts-ignore`×2、`any`×5、`as any`×2） |
| `typescript/domain/user.ts` | `Date.now`、`react` インポート | 2 | 0 | 0 |
| `typescript/components/UserCard.tsx` | *(スキャン対象外)* | 0 | 0 | 4 |

**TP=25, FP=2, FN=21 → 適合率=92.6%, 再現率=54.3%, F1=68.4%**

---

### 実行 6 — 2026-07-04 08:30 UTC

**スキャン対象:** 実行 5 と同一  
**エンジン報告:** "Detected 28 architecture policy violation(s)"

> 一部のアラートコメントでルールタイトル（角括弧内）と本文の内容が不一致。判定はコメント本文とファイルパスに基づく。

#### Go（ファイル別）

| ファイル | 検知内容 | TP | FP | FN | 備考 |
|---------|---------|:--:|:--:|:--:|------|
| `go/domain/user.go` | `database/sql` インポート | 1 | 0 | 1 | json/db タグを見逃し |
| `go/handler/user.go` | `panic`、`fmt.Println`、`context.Background`、repository インポート、`net/http` インポート | 4 | 1 | 0 | `net/http` は handler 層では正当 → **FP** |
| `go/repository/user.go` | `database/sql` インポート | 0 | 1 | 2 | `database/sql` は repository 層では正当 → **FP**；sqlc 返却型・Fix フォーマット見逃し |
| `go/service/user.go` | `panic`、`log.Fatal`、`fmt.Println`、`net/http` | 4 | 1 | 1 | `database/sql` は service 層に存在しない → **FP**；`context.Background` 見逃し |

#### Python（`python/bad_example.py`）

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `print`（×3） | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `requests.get/post`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `datetime.now` | ✅ | TP |
| service 層での `db.commit` | ✅ | **TP（全実行で初めての検知）** |
| bare `except` | ❌ | FN |
| `except` 内で `pass` | ❌ | FN |

#### TypeScript

| ファイル | 検知内容 | TP | FP | FN |
|---------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `any`×2、`Date.now`×2 | 2 | 2 | 11 — `Date.now` はこのファイルに存在しない → **FP×2** |
| `typescript/domain/user.ts` | `react` インポート、「UI から DB クライアントインポート」 | 1 | 1 | 1 — `react` インポートは TP；domain は UI ではない → **FP** |
| `typescript/components/UserCard.tsx` | *(スキャン対象外)* | 0 | 0 | 4 |

**TP=22, FP=6, FN=24 → 適合率=78.6%, 再現率=47.8%, F1=59.4%**

---

### 実行 7 — 2026-07-05 03:48 UTC

**スキャン対象:** 実行 5 と同一  
**エンジン報告:** "Detected 24 architecture policy violation(s)"

#### Go（ファイル別）

| ファイル | 検知内容 | TP | FP | FN | 備考 |
|---------|---------|:--:|:--:|:--:|------|
| `go/domain/user.go` | *(なし)* | 0 | 0 | 2 | database/sql・json/db タグ共に見逃し |
| `go/handler/user.go` | `panic`、`fmt.Println`、`context.Background`、`net/http` インポート | 3 | 1 | 1 | `net/http` は handler 層では正当 → **FP**；repository インポート見逃し |
| `go/repository/user.go` | *(なし)* | 0 | 0 | 2 | sqlc 返却型・Fix フォーマット見逃し |
| `go/service/user.go` | `net/http`、`panic`、`log.Fatal`、`fmt.Println`、`context.Background` | 5 | 0 | 0 | 完全検知 |

#### Python（`python/bad_example.py`）

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `print`（×3） | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `requests.get/post`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `datetime.now` | ✅ | TP |
| service 層での `db.commit` | ❌ | FN（実行 6 からのリグレッション） |
| bare `except` | ❌ | FN |
| `except` 内で `pass` | ❌ | FN |

#### TypeScript

| ファイル | 検知内容 | TP | FP | FN |
|---------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `console.log`×2、`fetch`×2 | 4 | 0 | 9（`@ts-ignore`×2、`any`×5、`as any`×2） |
| `typescript/domain/user.ts` | `Date.now`、`react` インポート | 2 | 0 | 0 |
| `typescript/components/UserCard.tsx` | *(スキャン対象外)* | 0 | 0 | 4 |

**TP=23, FP=1, FN=23 → 適合率=95.8%, 再現率=50.0%, F1=65.7%**

---

### 実行 8 — 2026-07-06 06:22 UTC

**スキャン対象:** 実行 5 と同一  
**エンジン報告:** "Detected 28 architecture policy violation(s)"

#### Go（ファイル別）

| ファイル | 検知内容 | TP | FP | FN | 備考 |
|---------|---------|:--:|:--:|:--:|------|
| `go/domain/user.go` | `database/sql` インポート | 1 | 0 | 1 | json/db タグ見逃し |
| `go/handler/user.go` | `panic`、`fmt.Println`、`context.Background`、repository インポート、`net/http` インポート | 4 | 1 | 0 | `net/http` は handler 層では正当 → **FP** |
| `go/repository/user.go` | `database/sql` インポート | 0 | 1 | 2 | `database/sql` は repository 層では正当 → **FP** |
| `go/service/user.go` | `panic`、`log.Fatal`、`fmt.Println`、`context.Background`、`net/http` | 5 | 0 | 0 | 完全検知 |

#### Python（`python/bad_example.py`）

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `print`（×3） | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `datetime.now` | ✅ | TP |
| `requests.get/post` | ❌ | FN（今回はまったく未検知） |
| service 層での `db.commit` | ❌ | FN |
| bare `except` | ❌ | FN |
| `except` 内で `pass` | ❌ | FN |

#### TypeScript

| ファイル | 検知内容 | TP | FP | FN |
|---------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `console.log`×4（正解は 2 件）、`fetch`×2 | 4 | 2 | 9 — `console.log` を 2 件過剰検知 → **FP×2** |
| `typescript/domain/user.ts` | `Date.now`、`react` インポート | 2 | 0 | 0 |
| `typescript/components/UserCard.tsx` | *(スキャン対象外)* | 0 | 0 | 4 |

**TP=24, FP=4, FN=22 → 適合率=85.7%, 再現率=52.2%, F1=64.9%**

---

### 実行 9 — 2026-07-07 07:10 UTC

**スキャン対象:** 実行 5 と同一  
**エンジン報告:** "Detected 17 architecture policy violation(s)"

> ルールメッセージが新フォーマットに移行。具体的な修正提案付きの説明文が採用された（例：「slog.Info(...) を使ってください」「logging.getLogger(__name__).info(...) を使ってください」）。より厳選されたルールセットで再現率を犠牲にして適合率を高めた可能性がある。

#### Go（ファイル別）

| ファイル | 検知内容 | TP | FP | FN | 備考 |
|---------|---------|:--:|:--:|:--:|------|
| `go/domain/user.go` | *(なし)* | 0 | 0 | 2 | |
| `go/handler/user.go` | `fmt.Println` | 1 | 0 | 3 | `panic`・`context.Background`・repository インポート見逃し |
| `go/repository/user.go` | *(なし)* | 0 | 0 | 2 | |
| `go/service/user.go` | `log.Fatal`、`fmt.Println` | 2 | 0 | 3 | `net/http`・`context.Background`・`panic` 見逃し |

#### Python（`python/bad_example.py`）

| ルール | 検知 | 判定 |
|--------|:----:|------|
| `print`（×3） | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `requests.get/post`（2 件中 1 件） | ✅/❌ | TP 1、FN 1 |
| `datetime.now` | ✅ | TP |
| service 層での `db.commit` | ❌ | FN |
| bare `except` | ❌ | FN |
| `except` 内で `pass` | ❌ | FN |

#### TypeScript

| ファイル | 検知内容 | TP | FP | FN |
|---------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `console.log`×2、`fetch`×2 | 4 | 0 | 9（`@ts-ignore`×2、`any`×5、`as any`×2） |
| `typescript/domain/user.ts` | `Date.now` | 1 | 0 | 1（`react` インポート見逃し） |
| `typescript/components/UserCard.tsx` | *(スキャン対象外)* | 0 | 0 | 4 |

**TP=17, FP=0, FN=29 → 適合率=100.0%, 再現率=37.0%, F1=54.0%**

---

## 言語別サマリー

### 実行 5（2026-07-03）

| 言語 | 正解数 | TP | FP | FN | 適合率 | 再現率 | F1 |
|:----:|:------:|:--:|:--:|:--:|:------:|:------:|:--:|
| Go | 13 | 10 | 2 | 3 | 83.3% | 76.9% | 80.0% |
| Python | 14 | 9 | 0 | 5 | 100.0% | 64.3% | 78.3% |
| TypeScript | 19 | 6 | 0 | 13 | 100.0% | 31.6% | 48.1% |
| **合計** | **46** | **25** | **2** | **21** | **92.6%** | **54.3%** | **68.4%** |

### 実行 6（2026-07-04）

| 言語 | 正解数 | TP | FP | FN | 適合率 | 再現率 | F1 |
|:----:|:------:|:--:|:--:|:--:|:------:|:------:|:--:|
| Go | 13 | 9 | 3 | 4 | 75.0% | 69.2% | 72.0% |
| Python | 14 | 10 | 0 | 4 | 100.0% | 71.4% | 83.3% |
| TypeScript | 19 | 3 | 3 | 16 | 50.0% | 15.8% | 24.0% |
| **合計** | **46** | **22** | **6** | **24** | **78.6%** | **47.8%** | **59.4%** |

### 実行 7（2026-07-05）

| 言語 | 正解数 | TP | FP | FN | 適合率 | 再現率 | F1 |
|:----:|:------:|:--:|:--:|:--:|:------:|:------:|:--:|
| Go | 13 | 8 | 1 | 5 | 88.9% | 61.5% | 72.7% |
| Python | 14 | 9 | 0 | 5 | 100.0% | 64.3% | 78.3% |
| TypeScript | 19 | 6 | 0 | 13 | 100.0% | 31.6% | 48.1% |
| **合計** | **46** | **23** | **1** | **23** | **95.8%** | **50.0%** | **65.7%** |

### 実行 8（2026-07-06）

| 言語 | 正解数 | TP | FP | FN | 適合率 | 再現率 | F1 |
|:----:|:------:|:--:|:--:|:--:|:------:|:------:|:--:|
| Go | 13 | 10 | 2 | 3 | 83.3% | 76.9% | 80.0% |
| Python | 14 | 8 | 0 | 6 | 100.0% | 57.1% | 72.7% |
| TypeScript | 19 | 6 | 2 | 13 | 75.0% | 31.6% | 44.4% |
| **合計** | **46** | **24** | **4** | **22** | **85.7%** | **52.2%** | **64.9%** |

### 実行 9（2026-07-07）

| 言語 | 正解数 | TP | FP | FN | 適合率 | 再現率 | F1 |
|:----:|:------:|:--:|:--:|:--:|:------:|:------:|:--:|
| Go | 13 | 3 | 0 | 10 | 100.0% | 23.1% | 37.5% |
| Python | 14 | 9 | 0 | 5 | 100.0% | 64.3% | 78.3% |
| TypeScript | 19 | 5 | 0 | 14 | 100.0% | 26.3% | 41.7% |
| **合計** | **46** | **17** | **0** | **29** | **100.0%** | **37.0%** | **54.0%** |

---

## 偽陽性（FP）分析

実行 5〜8（再編成レイアウト）の FP はすべてレイヤー判定の誤りが根本原因。実行ごとに FP の種類が変化した。

| 実行 | 箇所 | フラグが立ったルール | FP である理由 |
|:---:|------|------------------|-------------|
| 5〜8 | `go/handler/user.go` | 「service 層で `net/http` をインポートしない」 | handler 層であり `net/http` のインポートは正当 |
| 5〜8 | `go/repository/user.go` | 「domain 層で `database/sql` をインポートしない」 | repository 層であり `database/sql` のインポートは正当 |
| 6 | `go/service/user.go` | 「domain 層で `database/sql` をインポートしない」 | service 層は `database/sql` を一切インポートしていない |
| 6 | `typescript/bad_example.ts` | 「domain 層で `Date.now` を直接使用しない」 | `Date.now` はこのファイルに存在せず `domain/user.ts` に存在する |
| 6 | `typescript/domain/user.ts` | 「UI コンポーネントから DB クライアントをインポートしない」 | domain はUI コンポーネントではない |
| 8 | `typescript/bad_example.ts` | 「本番コードで `console.log` を使わない」×2 余分 | ファイル内の `console.log` は 2 件のみ；4 件フラグが立った |

**根本原因:** エンジンはパッケージパスではなくルールテキストやファイル内容の文字列マッチングでレイヤーを推定しており、正しく命名されたレイヤーディレクトリで誤分類が生じる。実行 6 ではファイル間の混同、実行 8 では出現回数の過剰カウントも観測された。実行 9 は厳選ルールセットに移行することで FP をゼロにした。

---

## 偽陰性（FN）分析：全実行を通じて未検知のルール

### Go
| ルール | 最初に見逃した実行 |
|--------|:----------------:|
| Domain model への `json`/`db` struct タグ | 実行 1 |
| Repository が domain model に変換せず sqlc 型を返す | 実行 1 |
| `Fix` フィールドを `Replacement` でラップしない | 実行 1 |

### Python
| ルール | 状況 |
|--------|------|
| bare `except:`（例外型の指定なし） | 全実行（実行 2〜9）で未検知 |
| `except` ブロック内での `pass` | 全実行（実行 2〜9）で未検知 |
| service 層での `db.commit()` | **実行 6 で初検知**；他の全実行では FN |

### TypeScript
| ルール | 状況 |
|--------|------|
| `any` 型アノテーション | 実行 6 のみ部分検知（5 件中 2 件）；他の全実行では FN |
| `as any` キャスト | 全実行（実行 4〜9）で未検知 |
| `@ts-ignore` コメント | 全実行（実行 4〜9）で未検知 |
| UI コンポーネントへの DB クライアントインポート | スキャン済みファイル内では全実行で未検知 |
| `dangerouslySetInnerHTML` | 全実行（実行 4〜9）で未検知 |

### TypeScript（実行 5〜9 でスキャン未実施）
| ルール | ファイル |
|--------|---------|
| UI コンポーネントへの DB クライアントインポート | `typescript/components/UserCard.tsx` |
| Client Component への `server-only` インポート | `typescript/components/UserCard.tsx` |
| `dangerouslySetInnerHTML` | `typescript/components/UserCard.tsx` |
| `any` 型（コンポーネント内） | `typescript/components/UserCard.tsx` |

---

## 部分検知のルール（N 件中 1 件のみ検知）

| ルール | ファイル | 検知件数 | 見逃し件数 |
|--------|---------|:--------:|:---------:|
| `pickle` | `python/bad_example.py` | 1 | 1 |
| `requests.get/post` | `python/bad_example.py` | 1 | 1 |

---

## 主な知見

1. **適合率は不安定、再現率は一貫して低い。** 実行 5〜9 の適合率は 78.6%（実行 6）〜100%（実行 9）、再現率は 37.0%（実行 9）〜54.3%（実行 5）の範囲で変動した。実行ごとに適合率と再現率のトレードオフが異なる。

2. **TypeScript のカバレッジが一貫して最も弱い。** TypeScript F1 は 24.0%（実行 6）〜48.1%（実行 5・7）の範囲。`any`・`as any`・`@ts-ignore` の 3 ルールは安定して検知されない。

3. **レイヤー判定が表層的。** `net/http` を handler でフラグ、`database/sql` を repository でフラグという同じ 2 件の FP が実行 5〜8 すべてで再現した。実行 6 ではファイル間の混同や存在しないインポートへのフラグも追加で発生した。実行 9 は厳選ルールセットで FP をゼロにした。

4. **意味論的ルールはほぼ検知できない。** コードフローやアーキテクチャ理解を要するルール（bare except、except-pass、sqlc 型の返却、json/db タグ、Fix/Replacement フォーマット）はほぼ未検知。`db.commit` は実行 6 で 1 度だけ検知された。

5. **実行 3 でのリグレッション（当日内）。** 実行 1 が同ファイルから 15 件中 11 件を検知したのに対し、実行 3 は 5 件のみ。02:36〜12:31 UTC の間にルール設定が一時的に縮退した可能性がある。

6. **実行 9 はエンジン設定の転換点。** 再編成レイアウトで初めて FP がゼロになったが、再現率は 37.0%（F1=54.0%）と最低水準。ルールメッセージが具体的な修正提案を含む新フォーマットに移行しており、より厳選されたルールセットが適用された可能性がある。

7. **Go のカバレッジが最も安定。** 実行 5〜8 の Go F1 は 72.0%〜80.0%。service 層の検知は安定しているが、domain・repository 層は一貫して弱い。実行 9 では Go 再現率が 23.1% まで低下した。
