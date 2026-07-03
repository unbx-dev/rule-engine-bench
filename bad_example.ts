// @ts-ignore
import React from "react"; // domain層でreactをimport
import { createClient } from "@supabase/supabase-js"; // UI componentからDB clientをimport
import "server-only"; // Client ComponentでServer-onlyモジュールをimport

// anyを使う
function processData(data: any): any {
  // console.logを本番コードで使う
  console.log("処理開始:", data);
  return data;
}

// as anyを使う
function castValue(value: unknown): string {
  return (value as any).toString();
}

// domain層でDate.nowを直接使う
function getCurrentTimestamp(): number {
  return Date.now();
}

// fetchを直接使う
async function fetchUser(userId: number) {
  // @ts-ignoreを使う
  // @ts-ignore
  const response = await fetch(`/api/users/${userId}`);
  const data: any = await response.json(); // anyを使う
  console.log("ユーザー取得:", data); // console.logを使う
  return data;
}

async function postEvent(payload: object): Promise<any> {
  // fetchを直接使う
  const res = await fetch("/api/events", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });
  return res.json() as any; // as anyを使う
}

// DB clientをUI componentで直接使う
const supabase = createClient(
  process.env.NEXT_PUBLIC_SUPABASE_URL!,
  process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
);

// "use client" なClient Component
// (コメントで示す: 本来は先頭に "use client" が付く想定)
function UserCard({ userId }: { userId: number }) {
  const [html, setHtml] = React.useState("<b>ロード中...</b>");

  React.useEffect(() => {
    // DB clientをClient Componentで直接使う
    supabase
      .from("users")
      .select("*")
      .eq("id", userId)
      .then(({ data }: { data: any }) => {
        console.log("DBから取得:", data); // console.logを使う
        setHtml(`<b>${data?.[0]?.name}</b>`);
      });
  }, [userId]);

  // dangerouslySetInnerHTMLを使う
  return <div dangerouslySetInnerHTML={{ __html: html }} />;
}

// domain層でreactをimportしてロジックを書く
function useDomainLogic() {
  const [value, setValue] = React.useState<any>(null); // anyを使う
  const timestamp = getCurrentTimestamp(); // Date.nowをdomain層で使う

  console.log("domain logic timestamp:", timestamp); // console.logを使う

  return { value, setValue, timestamp };
}

export { processData, castValue, fetchUser, postEvent, UserCard, useDomainLogic };
