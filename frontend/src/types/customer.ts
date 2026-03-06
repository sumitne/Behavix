/**
 * Matches backend GET /api/v1/customers and GET /api/v1/customers/{id} response shapes.
 */
export interface Customer {
  id: string;
  name: string;
  health_score: number; // 0–100
  active_users: number;
  last_activity: string; // ISO datetime
  email?: string;
}

export interface CustomerUsage {
  customer_id: string;
  period: string;
  active_users: number;
  events_count: number;
  feature_adoption?: Record<string, number>; // feature_id -> adoption %
}
