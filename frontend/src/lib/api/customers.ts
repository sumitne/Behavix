import type { Customer, CustomerUsage } from "@/types/customer";
import { apiFetch } from "./client";

/** Return empty list on any error (backend down, 404, 501, etc.) so the page never crashes. */
export async function fetchCustomers(): Promise<Customer[]> {
  try {
    const data = await apiFetch<Customer[]>("/api/v1/customers");
    return Array.isArray(data) ? data : [];
  } catch {
    return [];
  }
}

/** Return null on 404 or network error so callers can notFound() or show fallback. */
export async function fetchCustomer(id: string): Promise<Customer | null> {
  try {
    return await apiFetch<Customer>(`/api/v1/customers/${encodeURIComponent(id)}`);
  } catch {
    return null;
  }
}

/** Return null when endpoint is missing or backend is unreachable. */
export async function fetchCustomerUsage(id: string): Promise<CustomerUsage | null> {
  try {
    return await apiFetch<CustomerUsage>(
      `/api/v1/customers/${encodeURIComponent(id)}/usage`
    );
  } catch {
    return null;
  }
}
