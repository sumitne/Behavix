/**
 * Client-side fetch: calls Next.js API routes (same origin). No API key in browser.
 */
export async function clientFetch<T>(path: string): Promise<T> {
  const res = await fetch(path);
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    const err = new Error(typeof data?.error === "string" ? data.error : `Request failed: ${res.status}`);
    (err as Error & { status?: number }).status = res.status;
    throw err;
  }
  return data as T;
}
