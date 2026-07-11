# Rule Engine Benchmark Results

**Date measured:** 2026-07-03 — 2026-07-07 (latest run: 2026-07-07)  
**PR:** [#1 — RuleEngine検証用PR](https://github.com) (benchmark → main)  
**Engine:** Unbx Code Quarantine (`unbx-dev/unbx-cli@v0.0.10`)

---

## Methodology

Violations were extracted from GitHub PR inline review comments left by the CI engine.
Ground truth was established by manually auditing each source file against the 30 rules defined in the rule set (10 per language).

**Metrics:**
- **TP (True Positive):** violation correctly flagged by the engine
- **FP (False Positive):** engine flagged a location that does not violate the rule
- **FN (False Negative):** a real violation that the engine missed
- **Precision** = TP / (TP + FP)
- **Recall** = TP / (TP + FN)
- **F1** = 2 × Precision × Recall / (Precision + Recall)

---

## Scan Run Summary

| Run | Date (UTC) | Files Scanned | GT | TP | FP | FN | Precision | Recall | F1 |
|:---:|:----------:|:--------------|:--:|:--:|:--:|:--:|:---------:|:------:|:--:|
| 1 | 2026-07-03 02:29 | `bad_example.go` | 15 | 11 | 0 | 4 | 100.0% | 73.3% | 84.6% |
| 2 | 2026-07-03 02:36 | `bad_example.go`, `bad_example.py` | 29 | 14 | 0 | 15 | 100.0% | 48.3% | 65.1% |
| 3 | 2026-07-03 12:31 | `bad_example.go` | 15 | 5 | 0 | 10 | 100.0% | 33.3% | 50.0% |
| 4 | 2026-07-03 13:55 | `bad_example.go`, `bad_example.py`, `bad_example.ts` | 50 | 29 | 0 | 21 | 100.0% | 58.0% | 73.4% |
| 5 | 2026-07-03 14:34 | `go/*`, `python/*`, `typescript/*` | 46 | 25 | 2 | 21 | 92.6% | 54.3% | 68.4% |
| 6 | 2026-07-04 08:30 | `go/*`, `python/*`, `typescript/*` | 46 | 22 | 6 | 24 | 78.6% | 47.8% | 59.4% |
| 7 | 2026-07-05 03:48 | `go/*`, `python/*`, `typescript/*` | 46 | 23 | 1 | 23 | 95.8% | 50.0% | 65.7% |
| 8 | 2026-07-06 06:22 | `go/*`, `python/*`, `typescript/*` | 46 | 24 | 4 | 22 | 85.7% | 52.2% | 64.9% |
| 9 | 2026-07-07 07:10 | `go/*`, `python/*`, `typescript/*` | 46 | 17 | 0 | 29 | 100.0% | 37.0% | 54.0% |

> **File structure note:** Runs 1–4 scanned flat test files at the repository root. Run 5 (triggered by commit `84fbcd4`) scanned the restructured layout with per-language subdirectories and split files. Runs 6–9 used the same restructured layout. `typescript/components/UserCard.tsx` was not scanned in Runs 5–9; its 4 violations are counted as FN in the GT of 46.

---

## Ground Truth: Violations per File

### Go — `bad_example.go` (flat, Runs 1–4)

| # | Rule | Instances |
|---|------|-----------|
| 1 | `panic` outside `main` | 2 (service, handler) |
| 2 | `log.Fatal` | 3 (service, main×2) |
| 3 | `fmt.Println` in production code | 2 (service, handler) |
| 4 | `context.Background` in handler/service | 2 (service, handler) |
| 5 | `net/http` import in service layer | 1 |
| 6 | `database/sql` import in domain layer | 1 |
| 7 | handler imports repository directly | 1 |
| 8 | Repository returns sqlc type | 1 |
| 9 | Domain model has `json`/`db` struct tags | 1 |
| 10 | `Fix` field not wrapped in `Replacement` | 1 |
| **Total** | | **15** |

### Go — Restructured files (Run 5)

| File | Violations | Count |
|------|-----------|-------|
| `go/domain/user.go` | `database/sql` import; json/db struct tags | 2 |
| `go/handler/user.go` | repository import; `context.Background`; `panic`; `fmt.Println` | 4 |
| `go/repository/user.go` | returns sqlc type; `Fix` not in `Replacement` | 2 |
| `go/service/user.go` | `net/http` import; `context.Background`; `log.Fatal`; `panic`; `fmt.Println` | 5 |
| **Total** | | **13** |

### Python — `bad_example.py` / `python/bad_example.py` (Runs 2, 4, 5)

| # | Rule | Instances |
|---|------|-----------|
| 1 | `print` instead of `logger` | 3 (lines 9, 15, 77) |
| 2 | `eval` | 1 (line 14) |
| 3 | `exec` | 1 (line 26) |
| 4 | `os.system` | 1 (line 21) |
| 5 | `pickle` | 2 (dump line 32, load line 38) |
| 6 | `requests.get` direct call | 1 (line 43) |
| 7 | `requests.post` direct call | 1 (line 49) |
| 8 | `datetime.now` direct call | 1 (line 55) |
| 9 | bare `except:` | 1 (line 63) |
| 10 | `pass` in `except` block | 1 (line 64) |
| 11 | `db.commit()` in service layer | 1 (line 76) |
| **Total** | | **14** |

### TypeScript — `bad_example.ts` (flat, Run 4)

| # | Rule | Instances |
|---|------|-----------|
| 1 | `@ts-ignore` | 2 (lines 1, 25) |
| 2 | `any` type annotation | 6 (various) |
| 3 | `as any` cast | 2 |
| 4 | `console.log` | 4 (lines 9, 29, 56, 65) |
| 5 | direct `fetch` call | 2 |
| 6 | `Date.now` in domain layer | 1 |
| 7 | `react` import in domain layer | 1 |
| 8 | DB client import in UI component | 1 |
| 9 | `server-only` in Client Component | 1 |
| 10 | `dangerouslySetInnerHTML` | 1 |
| **Total** | | **21** |

### TypeScript — Restructured files (Run 5)

| File | Violations | Count |
|------|-----------|-------|
| `typescript/bad_example.ts` | `@ts-ignore`×2; `any`×5; `as any`×2; `console.log`×2; `fetch`×2 | 13 |
| `typescript/domain/user.ts` | `react` import; `Date.now` | 2 |
| `typescript/components/UserCard.tsx` | DB client import; `server-only`; `dangerouslySetInnerHTML`; `any` | 4 |
| **Total** | | **19** |

---

## Detailed Run Analysis

### Run 1 — 2026-07-03 02:29 UTC

**Scope:** `bad_example.go`  
**Engine report:** "Detected 11 architecture policy violation(s)"

| Rule | Detected | Status |
|------|:--------:|--------|
| `panic` outside `main` (×2) | ✅ | TP |
| `log.Fatal` (×3) | ✅ | TP |
| `fmt.Println` in production (×2) | ✅ | TP |
| `context.Background` in handler/service (×2) | ✅ | TP |
| `net/http` in service layer | ✅ | TP |
| `database/sql` in domain layer | ✅ | TP |
| handler imports repository | ❌ | FN |
| Repository returns sqlc type | ❌ | FN |
| Domain model has json/db tags | ❌ | FN |
| `Fix` not in `Replacement` | ❌ | FN |

**TP=11, FP=0, FN=4 → Precision=100.0%, Recall=73.3%, F1=84.6%**

---

### Run 2 — 2026-07-03 02:36 UTC

**Scope:** `bad_example.go`, `bad_example.py`  
**Engine report:** "Detected 14 architecture policy violation(s)"

Go detection: identical to Run 1 (11 TP, 0 FP, 4 FN).

| Python Rule | Detected | Status |
|-------------|:--------:|--------|
| `print` ×3 | ✅ | TP |
| `eval` | ❌ | FN |
| `exec` | ❌ | FN |
| `os.system` | ❌ | FN |
| `pickle` ×2 | ❌ | FN |
| `requests.get` | ❌ | FN |
| `requests.post` | ❌ | FN |
| `datetime.now` | ❌ | FN |
| bare `except` | ❌ | FN |
| `pass` in `except` | ❌ | FN |
| `db.commit` in service | ❌ | FN |

**TP=14, FP=0, FN=15 → Precision=100.0%, Recall=48.3%, F1=65.1%**

---

### Run 3 — 2026-07-03 12:31 UTC

**Scope:** `bad_example.go` (partial scan)  
**Engine report:** "Detected 5 architecture policy violation(s)"

| Rule | Detected | Status |
|------|:--------:|--------|
| `log.Fatal` ×3 | ✅ | TP |
| `fmt.Println` ×2 | ✅ | TP |
| `panic` ×2 | ❌ | FN |
| `context.Background` ×2 | ❌ | FN |
| `net/http` in service | ❌ | FN |
| `database/sql` in domain | ❌ | FN |
| handler imports repository | ❌ | FN |
| Repository returns sqlc type | ❌ | FN |
| Domain model json/db tags | ❌ | FN |
| `Fix` not in `Replacement` | ❌ | FN |

**TP=5, FP=0, FN=10 → Precision=100.0%, Recall=33.3%, F1=50.0%**

> This run appears to have used a different or experimental rule configuration, detecting only logging-related rules.

---

### Run 4 — 2026-07-03 13:55 UTC

**Scope:** `bad_example.go`, `bad_example.py`, `bad_example.ts`  
**Engine report:** "Detected 29 architecture policy violation(s)"

Go (bad_example.go): identical to Run 1 (11 TP, 0 FP, 4 FN).

Python (bad_example.py):

| Rule | Detected | Status |
|------|:--------:|--------|
| `print` ×3 | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `requests.get/post` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `datetime.now` | ✅ | TP |
| bare `except` | ❌ | FN |
| `pass` in `except` | ❌ | FN |
| `db.commit` in service | ❌ | FN |

TypeScript (bad_example.ts):

| Rule | Detected | Status |
|------|:--------:|--------|
| `console.log` ×4 | ✅ | TP |
| direct `fetch` ×2 | ✅ | TP |
| `Date.now` in domain | ✅ | TP |
| `react` import in domain | ✅ | TP |
| `server-only` in Client | ✅ | TP |
| `@ts-ignore` ×2 | ❌ | FN |
| `any` type ×6 | ❌ | FN |
| `as any` ×2 | ❌ | FN |
| DB client in UI | ❌ | FN |
| `dangerouslySetInnerHTML` | ❌ | FN |

**TP=29, FP=0, FN=21 → Precision=100.0%, Recall=58.0%, F1=73.4%**

---

### Run 5 — 2026-07-03 14:34 UTC

**Scope:** `go/domain/user.go`, `go/handler/user.go`, `go/repository/user.go`, `go/service/user.go`, `python/bad_example.py`, `typescript/bad_example.ts`, `typescript/domain/user.ts`  
**Engine report:** "Detected 27 architecture policy violation(s)"

> Note: `typescript/components/UserCard.tsx` was **not scanned** in this run.

#### Go (per file)

| File | Detected | TP | FP | FN | Notes |
|------|---------|:--:|:--:|:--:|-------|
| `go/domain/user.go` | `database/sql` import | 1 | 0 | 1 | json/db tags missed |
| `go/handler/user.go` | `panic`, `fmt.Println`, `context.Background`, repository import, `net/http` import | 4 | 1 | 0 | `net/http` in handler is **expected** — FP |
| `go/repository/user.go` | `database/sql` import | 0 | 1 | 2 | `database/sql` in repository is **expected** — FP; sqlc return type and Fix format both missed |
| `go/service/user.go` | `net/http`, `context.Background`, `log.Fatal`, `panic`, `fmt.Println` | 5 | 0 | 0 | Perfect detection |

#### Python (`python/bad_example.py`)

| Rule | Detected | Status |
|------|:--------:|--------|
| `print` ×3 | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `requests.get/post` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `datetime.now` | ✅ | TP |
| bare `except` | ❌ | FN |
| `pass` in `except` | ❌ | FN |
| `db.commit` in service | ❌ | FN |

#### TypeScript

| File | Detected | TP | FP | FN |
|------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `console.log`×2, `fetch`×2 | 4 | 0 | 9 (`@ts-ignore`×2, `any`×5, `as any`×2) |
| `typescript/domain/user.ts` | `Date.now`, `react` import | 2 | 0 | 0 |
| `typescript/components/UserCard.tsx` | *(not scanned)* | 0 | 0 | 4 |

**TP=25, FP=2, FN=21 → Precision=92.6%, Recall=54.3%, F1=68.4%**

---

### Run 6 — 2026-07-04 08:30 UTC

**Scope:** `go/domain/user.go`, `go/handler/user.go`, `go/repository/user.go`, `go/service/user.go`, `python/bad_example.py`, `typescript/bad_example.ts`, `typescript/domain/user.ts`  
**Engine report:** "Detected 28 architecture policy violation(s)"

> Note: `typescript/components/UserCard.tsx` was **not scanned**. Several alert comments had mismatched title/body text (the rule title in brackets did not match the violation description in the body); judgements below are based on the violation body text and the file the comment was attached to.

#### Go (per file)

| File | Detected | TP | FP | FN | Notes |
|------|---------|:--:|:--:|:--:|-------|
| `go/domain/user.go` | `database/sql` import | 1 | 0 | 1 | json/db tags missed |
| `go/handler/user.go` | `panic`, `fmt.Println`, `context.Background`, repository import, `net/http` import | 4 | 1 | 0 | `net/http` in handler is **expected** — FP |
| `go/repository/user.go` | `database/sql` import | 0 | 1 | 2 | `database/sql` in repository is **expected** — FP; sqlc return type and Fix format both missed |
| `go/service/user.go` | `panic`, `log.Fatal`, `fmt.Println`, `net/http` | 4 | 1 | 1 | `database/sql` flagged in service (not imported there) — FP; `context.Background` missed |

#### Python (`python/bad_example.py`)

| Rule | Detected | Status |
|------|:--------:|--------|
| `print` ×3 | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `requests.get/post` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `datetime.now` | ✅ | TP |
| `db.commit` in service | ✅ | **TP (first detection across all runs)** |
| bare `except` | ❌ | FN |
| `pass` in `except` | ❌ | FN |

#### TypeScript

| File | Detected | TP | FP | FN |
|------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `any`×2, `Date.now`×2 | 2 | 2 | 11 (`@ts-ignore`×2, `any`×3, `as any`×2, `console.log`×2, `fetch`×2) — `Date.now` not in this file → FP×2 |
| `typescript/domain/user.ts` | `react` import, DB client in domain | 1 | 1 | 1 | `react` import TP; "DB client in UI" flagged on domain → FP; `Date.now` missed |
| `typescript/components/UserCard.tsx` | *(not scanned)* | 0 | 0 | 4 |

**TP=22, FP=6, FN=24 → Precision=78.6%, Recall=47.8%, F1=59.4%**

---

### Run 7 — 2026-07-05 03:48 UTC

**Scope:** same as Run 6  
**Engine report:** "Detected 24 architecture policy violation(s)"

#### Go (per file)

| File | Detected | TP | FP | FN | Notes |
|------|---------|:--:|:--:|:--:|-------|
| `go/domain/user.go` | *(nothing)* | 0 | 0 | 2 | database/sql and json/db tags both missed |
| `go/handler/user.go` | `panic`, `fmt.Println`, `context.Background`, `net/http` import | 3 | 1 | 1 | `net/http` in handler is **expected** — FP; repository import missed |
| `go/repository/user.go` | *(nothing)* | 0 | 0 | 2 | sqlc return type and Fix format both missed |
| `go/service/user.go` | `net/http`, `panic`, `log.Fatal`, `fmt.Println`, `context.Background` | 5 | 0 | 0 | Perfect detection |

#### Python (`python/bad_example.py`)

| Rule | Detected | Status |
|------|:--------:|--------|
| `print` ×3 | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `requests.get/post` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `datetime.now` | ✅ | TP |
| `db.commit` in service | ❌ | FN (regression vs Run 6) |
| bare `except` | ❌ | FN |
| `pass` in `except` | ❌ | FN |

#### TypeScript

| File | Detected | TP | FP | FN |
|------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `console.log`×2, `fetch`×2 | 4 | 0 | 9 (`@ts-ignore`×2, `any`×5, `as any`×2) |
| `typescript/domain/user.ts` | `Date.now`, `react` import | 2 | 0 | 0 |
| `typescript/components/UserCard.tsx` | *(not scanned)* | 0 | 0 | 4 |

**TP=23, FP=1, FN=23 → Precision=95.8%, Recall=50.0%, F1=65.7%**

---

### Run 8 — 2026-07-06 06:22 UTC

**Scope:** same as Run 6  
**Engine report:** "Detected 28 architecture policy violation(s)"

#### Go (per file)

| File | Detected | TP | FP | FN | Notes |
|------|---------|:--:|:--:|:--:|-------|
| `go/domain/user.go` | `database/sql` import | 1 | 0 | 1 | json/db tags missed |
| `go/handler/user.go` | `panic`, `fmt.Println`, `context.Background`, repository import, `net/http` import | 4 | 1 | 0 | `net/http` in handler is **expected** — FP |
| `go/repository/user.go` | `database/sql` import | 0 | 1 | 2 | `database/sql` in repository is **expected** — FP; sqlc return type and Fix format missed |
| `go/service/user.go` | `panic`, `log.Fatal`, `fmt.Println`, `context.Background`, `net/http` | 5 | 0 | 0 | Perfect detection |

#### Python (`python/bad_example.py`)

| Rule | Detected | Status |
|------|:--------:|--------|
| `print` ×3 | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `datetime.now` | ✅ | TP |
| `requests.get/post` | ❌ | FN (regression — not detected at all) |
| `db.commit` in service | ❌ | FN |
| bare `except` | ❌ | FN |
| `pass` in `except` | ❌ | FN |

#### TypeScript

| File | Detected | TP | FP | FN |
|------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `console.log`×4 (GT=2), `fetch`×2 | 4 | 2 | 9 — `console.log` over-detected (2 FPs); `@ts-ignore`×2, `any`×5, `as any`×2 missed |
| `typescript/domain/user.ts` | `Date.now`, `react` import | 2 | 0 | 0 |
| `typescript/components/UserCard.tsx` | *(not scanned)* | 0 | 0 | 4 |

**TP=24, FP=4, FN=22 → Precision=85.7%, Recall=52.2%, F1=64.9%**

---

### Run 9 — 2026-07-07 07:10 UTC

**Scope:** same as Run 6  
**Engine report:** "Detected 17 architecture policy violation(s)"

> Note: Rule messages in this run adopted a new format — more descriptive with concrete code suggestions (e.g., "use slog.Info(...) instead", "use logging.getLogger(__name__).info(...)"). This appears to be a tighter, higher-confidence rule set that trades recall for precision.

#### Go (per file)

| File | Detected | TP | FP | FN | Notes |
|------|---------|:--:|:--:|:--:|-------|
| `go/domain/user.go` | *(nothing)* | 0 | 0 | 2 | |
| `go/handler/user.go` | `fmt.Println` | 1 | 0 | 3 | `panic`, `context.Background`, repository import all missed |
| `go/repository/user.go` | *(nothing)* | 0 | 0 | 2 | |
| `go/service/user.go` | `log.Fatal`, `fmt.Println` | 2 | 0 | 3 | `net/http`, `context.Background`, `panic` all missed |

#### Python (`python/bad_example.py`)

| Rule | Detected | Status |
|------|:--------:|--------|
| `print` ×3 | ✅ | TP |
| `eval` | ✅ | TP |
| `exec` | ✅ | TP |
| `os.system` | ✅ | TP |
| `pickle` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `requests.get/post` (1 of 2) | ✅/❌ | 1 TP, 1 FN |
| `datetime.now` | ✅ | TP |
| `db.commit` in service | ❌ | FN |
| bare `except` | ❌ | FN |
| `pass` in `except` | ❌ | FN |

#### TypeScript

| File | Detected | TP | FP | FN |
|------|---------|:--:|:--:|:--|
| `typescript/bad_example.ts` | `console.log`×2, `fetch`×2 | 4 | 0 | 9 (`@ts-ignore`×2, `any`×5, `as any`×2) |
| `typescript/domain/user.ts` | `Date.now` | 1 | 0 | 1 (`react` import missed) |
| `typescript/components/UserCard.tsx` | *(not scanned)* | 0 | 0 | 4 |

**TP=17, FP=0, FN=29 → Precision=100.0%, Recall=37.0%, F1=54.0%**

---

## Per-Language Summaries

### Run 5 (2026-07-03)

| Language | GT | TP | FP | FN | Precision | Recall | F1 |
|:--------:|:--:|:--:|:--:|:--:|:---------:|:------:|:--:|
| Go | 13 | 10 | 2 | 3 | 83.3% | 76.9% | 80.0% |
| Python | 14 | 9 | 0 | 5 | 100.0% | 64.3% | 78.3% |
| TypeScript | 19 | 6 | 0 | 13 | 100.0% | 31.6% | 48.1% |
| **Total** | **46** | **25** | **2** | **21** | **92.6%** | **54.3%** | **68.4%** |

### Run 6 (2026-07-04)

| Language | GT | TP | FP | FN | Precision | Recall | F1 |
|:--------:|:--:|:--:|:--:|:--:|:---------:|:------:|:--:|
| Go | 13 | 9 | 3 | 4 | 75.0% | 69.2% | 72.0% |
| Python | 14 | 10 | 0 | 4 | 100.0% | 71.4% | 83.3% |
| TypeScript | 19 | 3 | 3 | 16 | 50.0% | 15.8% | 24.0% |
| **Total** | **46** | **22** | **6** | **24** | **78.6%** | **47.8%** | **59.4%** |

### Run 7 (2026-07-05)

| Language | GT | TP | FP | FN | Precision | Recall | F1 |
|:--------:|:--:|:--:|:--:|:--:|:---------:|:------:|:--:|
| Go | 13 | 8 | 1 | 5 | 88.9% | 61.5% | 72.7% |
| Python | 14 | 9 | 0 | 5 | 100.0% | 64.3% | 78.3% |
| TypeScript | 19 | 6 | 0 | 13 | 100.0% | 31.6% | 48.1% |
| **Total** | **46** | **23** | **1** | **23** | **95.8%** | **50.0%** | **65.7%** |

### Run 8 (2026-07-06)

| Language | GT | TP | FP | FN | Precision | Recall | F1 |
|:--------:|:--:|:--:|:--:|:--:|:---------:|:------:|:--:|
| Go | 13 | 10 | 2 | 3 | 83.3% | 76.9% | 80.0% |
| Python | 14 | 8 | 0 | 6 | 100.0% | 57.1% | 72.7% |
| TypeScript | 19 | 6 | 2 | 13 | 75.0% | 31.6% | 44.4% |
| **Total** | **46** | **24** | **4** | **22** | **85.7%** | **52.2%** | **64.9%** |

### Run 9 (2026-07-07)

| Language | GT | TP | FP | FN | Precision | Recall | F1 |
|:--------:|:--:|:--:|:--:|:--:|:---------:|:------:|:--:|
| Go | 13 | 3 | 0 | 10 | 100.0% | 23.1% | 37.5% |
| Python | 14 | 9 | 0 | 5 | 100.0% | 64.3% | 78.3% |
| TypeScript | 19 | 5 | 0 | 14 | 100.0% | 26.3% | 41.7% |
| **Total** | **46** | **17** | **0** | **29** | **100.0%** | **37.0%** | **54.0%** |

---

## False Positive Analysis

FPs across all restructured-layout runs (Runs 5–8) share the same root cause — layer-detection errors — but the set of FP types grew over time.

| Run | Location | Flagged Rule | Reason It Is a FP |
|:---:|----------|-------------|-------------------|
| 5–8 | `go/handler/user.go` | "Do not import net/http in the service layer" | `handler/user.go` is the handler layer; `net/http` is expected there |
| 5–8 | `go/repository/user.go` | "Do not import database/sql in the domain layer" | `repository/user.go` is the repository layer; `database/sql` is expected there |
| 6 | `go/service/user.go` | "Do not import database/sql in the domain layer" | `service/user.go` does not import `database/sql` at all |
| 6 | `typescript/bad_example.ts` | "Do not use Date.now in the domain layer" | `Date.now` is not present in `bad_example.ts`; it lives in `domain/user.ts` |
| 6 | `typescript/domain/user.ts` | "Do not import DB client from UI component" | `domain/user.ts` is not a UI component |
| 8 | `typescript/bad_example.ts` | "Do not use console.log" ×2 extra | Only 2 `console.log` calls exist in the file; 4 were flagged |

**Root cause:** The engine infers architectural layer from rule text or file content rather than the package path, leading to misclassification in correctly-named layer directories. Run 6 also shows evidence of cross-file confusion (flagging a rule on the wrong file). Run 8 shows over-counting of occurrences. Run 9 eliminated all FPs by adopting a narrower, higher-confidence rule set.

---

## False Negative Analysis: Consistently Missed Rules

Rules that were **never detected** across all nine runs:

### Go
| Rule | Missed Since |
|------|-------------|
| Domain model has `json`/`db` struct tags | Run 1 |
| Repository returns sqlc type instead of domain model | Run 1 |
| `Fix` field not wrapped in `Replacement` | Run 1 |

### Python
| Rule | Status |
|------|--------|
| bare `except:` (no exception type) | Never detected (Runs 2–9) |
| `pass` in `except` block | Never detected (Runs 2–9) |
| `db.commit()` in service layer | Detected **once** in Run 6; FN in all other runs |

### TypeScript
| Rule | Status |
|------|--------|
| `any` type annotation | Detected partially (×2 of 5) in Run 6 only; FN in all other runs |
| `as any` cast | Never detected (Runs 4–9) |
| `@ts-ignore` comment | Never detected (Runs 4–9) |
| DB client (`createClient`) import in UI component | Never detected in scanned files (Runs 4–9) |
| `dangerouslySetInnerHTML` | Never detected (Runs 4–9) |

### TypeScript (file never scanned — Runs 5–9)
| Rule | File |
|------|------|
| DB client import in UI | `typescript/components/UserCard.tsx` |
| `server-only` in Client Component | `typescript/components/UserCard.tsx` |
| `dangerouslySetInnerHTML` | `typescript/components/UserCard.tsx` |
| `any` type in component | `typescript/components/UserCard.tsx` |

---

## Partially Detected Rules (1 of N occurrences)

| Rule | File | Detected | Missed |
|------|------|:--------:|:------:|
| `pickle` | `python/bad_example.py` | 1 | 1 |
| `requests.get/post` | `python/bad_example.py` | 1 | 1 |

---

## Key Findings

1. **High precision, unstable recall.** Across Runs 5–9, precision ranged from 78.6% (Run 6) to 100% (Run 9), while recall ranged from 37.0% (Run 9) to 54.3% (Run 5). The engine trades off precision against recall differently in each run.

2. **TypeScript coverage is consistently weakest.** F1 for TypeScript ranged from 24.0% (Run 6) to 48.1% (Runs 5, 7). `any`, `as any`, and `@ts-ignore` are never reliably detected. Runs 5, 7, and 9 achieved 0 TypeScript FPs; Runs 6 and 8 introduced cross-file confusion and over-counting.

3. **Layer detection is shallow.** The same two FPs — `net/http` in handler and `database/sql` in repository — recurred in every restructured-layout run (5–8). Run 6 added three new FP types (wrong file, wrong layer type, non-existent import), while Run 9 eliminated all FPs by adopting a narrower rule set.

4. **Semantic rules are rarely detected.** Rules requiring code-flow or architectural understanding (bare `except`, `except`-`pass`, `db.commit` in service, Repository returning sqlc types, json/db struct tags, Fix/Replacement format) are almost never caught. `db.commit` was detected **only once** in Run 6; all others have never been detected.

5. **Run 3 regression (intra-day).** Run 3 detected only 5 of 15 violations from the same file that Run 1 detected 11 of 15, suggesting the rule set was temporarily degraded between 02:36 UTC and 12:31 UTC on 2026-07-03.

6. **Run 9 marks a configuration shift.** Run 9 produced the first 0-FP result in the restructured-layout series but had the lowest recall (37.0%, F1=54.0%). Rule messages adopted a new descriptive format with code suggestions, indicating the engine may be running on a tighter, more curated rule set.

7. **Go coverage is most consistent.** Go F1 across Runs 5–8 ranged 72.0%–80.0%. Service layer detection is reliable; domain and repository layer detection remains weak. Run 9 dropped Go recall to 23.1% due to the narrower rule set.
