import type { Insight } from "@/types/insight";
import { apiFetch } from "./client";

/**
 * Raw insight from backend (id may be UUID string, created_at ISO string).
 */
export type InsightResponse = Insight;

export async function fetchInsights(limit = 20): Promise<Insight[]> {
  const query = new URLSearchParams({ limit: String(limit) });
  const data = await apiFetch<InsightResponse[]>(
    `/api/v1/insights?${query.toString()}`
  );
  return Array.isArray(data) ? data : [];
}
