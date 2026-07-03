"use client";

import React from "react";
import { createClient } from "@supabase/supabase-js"; // ❌ UI componentからDB clientをimport
import "server-only"; // ❌ Client ComponentでServer-onlyモジュールをimport

// ❌ DB clientをClient Componentで直接初期化
const supabase = createClient(
  process.env.NEXT_PUBLIC_SUPABASE_URL!,
  process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
);

function UserCard({ userId }: { userId: number }) {
  const [name, setName] = React.useState("");

  React.useEffect(() => {
    supabase
      .from("users")
      .select("name")
      .eq("id", userId)
      .single()
      .then(({ data }: { data: any }) => setName(data?.name ?? ""));
  }, [userId]);

  // ❌ dangerouslySetInnerHTMLを使う
  return <div dangerouslySetInnerHTML={{ __html: `<b>${name}</b>` }} />;
}

export { UserCard };
