// ❌ @ts-ignoreを使う
// @ts-ignore
const _unused: any = null; // ❌ anyを使う

// ❌ anyを使う
function processData(data: any): any {
  // ❌ console.logを本番コードで使う
  console.log("処理開始:", data);
  return data;
}

// ❌ as anyを使う
function castValue(value: unknown): string {
  return (value as any).toString();
}

// ❌ fetchを直接使う
async function fetchUser(userId: number) {
  // ❌ @ts-ignoreを使う
  // @ts-ignore
  const response = await fetch(`/api/users/${userId}`);
  const data: any = await response.json(); // ❌ anyを使う
  console.log("ユーザー取得:", data); // ❌ console.logを使う
  return data;
}

// ❌ fetchを直接使う / as any / Promise<any>
async function postEvent(payload: object): Promise<any> {
  const res = await fetch("/api/events", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  return res.json() as any; // ❌ as anyを使う
}

export { processData, castValue, fetchUser, postEvent };
