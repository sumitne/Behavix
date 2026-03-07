import { NextResponse } from "next/server";

const getApiBaseUrl = () => {
  const url = process.env.API_BASE_URL ?? process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
  return url.replace(/\/$/, "");
};

const getApiKey = () => process.env.BEHAVIX_API_KEY ?? process.env.API_KEY;

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const limit = searchParams.get("limit") ?? "20";
  const base = getApiBaseUrl();
  const key = getApiKey();
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(key ? { Authorization: `Bearer ${key}` } : {}),
  };
  const res = await fetch(`${base}/api/v1/insights?limit=${limit}`, { headers });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    return NextResponse.json(data, { status: res.status });
  }
  return NextResponse.json(data);
}
