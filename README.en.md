# Rule Engine Bench

A benchmark repository for evaluating the detection accuracy of an automated architecture-policy rule engine (Unbx Code Quarantine). The CI workflow opens inline PR review comments for every rule violation it finds; this repo measures how many violations it catches versus how many it misses.

## Rule Set

### Python

| Rule | Description |
|------|-------------|
| Use `logger`, not `print` | All logging must use the `logging` module |
| Do not use `eval` | Arbitrary code execution risk |
| Do not use `exec` | Arbitrary code execution risk |
| Do not use `os.system` | Shell injection risk |
| Do not use `pickle` | Arbitrary code execution risk |
| Do not use bare `except` | Specify a concrete exception type (e.g. `except Exception`) |
| Do not `pass` inside `except` | Handle exceptions properly; do not swallow them |
| Do not call `requests.get`/`post` directly | HTTP clients must go through a wrapper |
| Do not call `datetime.now` directly | Use DI or abstraction for testability |
| Do not `commit` the DB session in the service layer | Transaction management belongs above the service layer |

---

### Go

| Rule | Description |
|------|-------------|
| Do not use `panic` outside `main` | Return errors in library/service code |
| Do not use `log.Fatal` | Prohibited because it implicitly calls `os.Exit` |
| Do not use `fmt.Println` in production code | Use a logging library |
| Do not use `context.Background` in handler/service/repository | Propagate context from the caller |
| Do not import `net/http` in the service layer | Keep HTTP dependencies in the handler layer |
| Do not import `database/sql` in the domain layer | Keep DB dependencies in the repository layer |
| Do not import repository in the handler layer | Handlers must access data through the service layer |
| Repository must not return sqlc types | Convert to domain models before returning |
| Do not add `json`/`db` tags to domain models | Do not bring infrastructure concerns into the domain layer |
| Return fixes as `Replacement` | Auto-fixes must be provided in Replacement format |

---

### TypeScript

| Rule | Description |
|------|-------------|
| Do not use `console.log` in production code | Use a logging library |
| Do not use `any` | Prohibits any type that breaks type safety |
| Do not use `as any` | Prohibits bypassing the type system via assertion |
| Do not use `@ts-ignore` | Fix type errors properly |
| Do not call `fetch` directly | HTTP clients must go through a wrapper |
| Do not use `Date.now` in the domain layer | Use DI or abstraction for testability |
| Do not import `react` in the domain layer | Do not bring UI dependencies into the domain layer |
| Do not import a DB client from a UI component | Access data through the server layer |
| Do not import `server-only` modules in Client Components | Prevent server-only code from leaking to the client |
| Do not use `dangerouslySetInnerHTML` | XSS risk |

---

## Repository Structure

```
rule-engine-bench/
├── go/
│   ├── domain/user.go        # domain layer — intentional violations
│   ├── handler/user.go       # handler layer — intentional violations
│   ├── repository/user.go    # repository layer — intentional violations
│   ├── service/user.go       # service layer — intentional violations
│   └── go.mod
├── python/
│   └── bad_example.py        # intentional violations of all 10 Python rules
├── typescript/
│   ├── bad_example.ts        # intentional violations (any, fetch, console.log, …)
│   ├── domain/user.ts        # domain layer — react import, Date.now
│   └── components/UserCard.tsx  # Client Component — DB client, server-only, dangerouslySetInnerHTML
├── BENCHMARK_RESULTS.md      # accuracy measurements (see below)
└── README.en.md              # this file
```

Each source file contains intentional policy violations, annotated with `❌` comments that identify which rule is being tested.

---

## Benchmark Results Summary

Full results are in [`BENCHMARK_RESULTS.md`](./BENCHMARK_RESULTS.md). The table below shows the final scan (Run 5, 2026-07-03 14:34 UTC), which covers the restructured file layout.

| Language | Ground Truth | TP | FP | FN | Precision | Recall | F1 |
|:--------:|:------------:|:--:|:--:|:--:|:---------:|:------:|:--:|
| Go | 13 | 10 | 2 | 3 | 83.3% | 76.9% | 80.0% |
| Python | 14 | 9 | 0 | 5 | 100.0% | 64.3% | 78.3% |
| TypeScript | 19 | 6 | 0 | 13 | 100.0% | 31.6% | 48.1% |
| **Total** | **46** | **25** | **2** | **21** | **92.6%** | **54.3%** | **68.4%** |

### Key observations

- **High precision (~93%), low recall (~54%).** The engine rarely reports a false alarm, but it misses nearly half of all violations.
- **TypeScript coverage is weakest** (F1 48%): `any`, `as any`, and `@ts-ignore` rules are never triggered.
- **Semantic rules are not detected**: bare `except`, `pass` in `except`, `db.commit` in service layer, returning sqlc types, and domain-model struct tags are consistently missed.
- **Layer detection has false positives**: `net/http` in a handler file and `database/sql` in a repository file were incorrectly flagged as domain/service violations.

---

## CI Workflow

The repository uses GitHub Actions with the `unbx-dev/unbx-cli` action:

```yaml
on:
  pull_request:
    types: [opened, synchronize, reopened]
```

On every push to the PR branch the engine scans the diff, generates inline review comments for each detected violation, and posts a summary comment reporting the total violation count.

---

## How to Add Rules

Rules are managed in the Unbx platform and referenced via `UNBX_ACCESS_TOKEN` / `UNBX_SECRET_TOKEN` repository secrets. To extend the benchmark:

1. Add a new source file (or extend an existing one) with intentional violations of the new rule.
2. Annotate each violation with an `❌` comment describing which rule it violates.
3. Update the rule table in this README.
4. Record the new ground truth count in `BENCHMARK_RESULTS.md` after the CI runs.
