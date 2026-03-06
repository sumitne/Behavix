/**
 * API config for backend Behavix-AI. Use server-side only for the key.
 * Set in .env.local: API_BASE_URL, BEHAVIX_API_KEY
 */
function getApiBaseUrl(): string {
  const url = process.env.API_BASE_URL ?? process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";
  return url.replace(/\/$/, "");
}

function getApiKey(): string | undefined {
  return process.env.BEHAVIX_API_KEY ?? process.env.API_KEY;
}

export { getApiBaseUrl, getApiKey };
