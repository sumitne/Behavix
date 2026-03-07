import { NextResponse } from "next/server";

const getApiBaseUrl = () => {
  const url = process.env.API_BASE_URL ?? process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
  return url.replace(/\/$/, "");
};

const getApiKey = () => process.env.BEHAVIX_API_KEY ?? process.env.API_KEY;

export async function GET(
  _request: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  const { id } = await params;
  const base = getApiBaseUrl();
  const key = getApiKey();
  const headers: HeadersInit = {
    "Content-Type": "application/json",
    ...(key ? { Authorization: `Bearer ${key}` } : {}),
  };
  const res = await fetch(`${base}/api/v1/customers/${encodeURIComponent(id)}`, { headers });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    return NextResponse.json(data, { status: res.status });
  }
  return NextResponse.json(data);
}
