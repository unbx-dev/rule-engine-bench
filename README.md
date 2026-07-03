# Rule Engine Bench

## 検証用ルールセット

### Python

| ルール | 説明 |
|--------|------|
| `print` を使わずに `logger` を使う | ロギングは `logging` モジュールを使用する |
| `eval` を使わない | 任意コード実行の危険性があるため禁止 |
| `exec` を使わない | 任意コード実行の危険性があるため禁止 |
| `os.system` を使わない | シェルインジェクションの危険性があるため禁止 |
| `pickle` を使わない | 任意コード実行の危険性があるため禁止 |
| bare `except` を使わない | `except Exception` など具体的な例外を指定する |
| `except` 内で `pass` しない | 例外を握り潰さず適切に処理する |
| `requests.get`/`post` を直接使わない | HTTPクライアントはラッパーを経由する |
| `datetime.now` を直接使わない | テスト可能性のためにDIや抽象化を経由する |
| service層でDB sessionを `commit` しない | トランザクション管理はservice層より上位で行う |

---

### Go

| ルール | 説明 |
|--------|------|
| `panic` を `main` 以外で使わない | ライブラリ・サービスコードではエラーを返す |
| `log.Fatal` を使わない | `os.Exit` を暗黙的に呼ぶため禁止 |
| `fmt.Println` を本番コードで使わない | ロギングライブラリを使用する |
| `context.Background` を handler/service/repository で使わない | 呼び出し元から `context` を伝搬させる |
| service層で `net/http` を import しない | HTTP関連の依存はhandler層に留める |
| domain層で `database/sql` を import しない | DBの依存はrepository層に留める |
| handler層で repository を import しない | handler はservice層を経由してデータにアクセスする |
| Repository は sqlc 型を返さない | domain model に変換してから返す |
| Domain model に `json`/`db` tag を付けない | インフラの関心事をdomain層に持ち込まない |
| Fix は Replacement で返す | 自動修正はReplacement形式で提供する |

---

### TypeScript

| ルール | 説明 |
|--------|------|
| `console.log` を本番コードで使わない | ロギングライブラリを使用する |
| `any` を使わない | 型安全性を損なうため禁止 |
| `as any` を使わない | 型アサーションによる型安全性の回避を禁止 |
| `@ts-ignore` を使わない | 型エラーは適切に修正する |
| `fetch` を直接使わない | HTTPクライアントはラッパーを経由する |
| `Date.now` をdomain層で直接使わない | テスト可能性のためにDIや抽象化を経由する |
| domain層で `react` を import しない | UIの依存をdomain層に持ち込まない |
| UI componentからDB clientを import しない | データアクセスはserver層を経由する |
| Client Componentでserver-onlyモジュールを import しない | サーバー専用コードのクライアント混入を防ぐ |
| `dangerouslySetInnerHTML` を使わない | XSSの危険性があるため禁止 |
