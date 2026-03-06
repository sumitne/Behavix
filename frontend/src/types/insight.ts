/**
 * Matches backend GET /api/v1/insights response shape.
 */
export type InsightSeverity = "low" | "medium" | "high";

export interface Insight {
  id: string;
  type: string;
  title: string;
  description: string;
  severity: InsightSeverity;
  created_at: string; // ISO datetime
  metric_value?: number;
  baseline_value?: number;
  /** For display: resolved customer name (from tenant/customer context) */
  customer_name?: string;
}
