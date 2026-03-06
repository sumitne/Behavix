import type { Customer } from "@/types/customer";

export const mockCustomers: Customer[] = [
  { id: "cust-1", name: "Acme Corp", health_score: 72, active_users: 317, last_activity: "2025-03-06T08:00:00Z" },
  { id: "cust-2", name: "Beta Inc", health_score: 88, active_users: 124, last_activity: "2025-03-06T09:30:00Z" },
  { id: "cust-3", name: "Gamma LLC", health_score: 95, active_users: 56, last_activity: "2025-03-06T07:15:00Z" },
  { id: "cust-4", name: "Delta Co", health_score: 61, active_users: 89, last_activity: "2025-03-05T16:00:00Z" },
  { id: "cust-5", name: "Epsilon Ltd", health_score: 91, active_users: 42, last_activity: "2025-03-06T10:00:00Z" },
];
