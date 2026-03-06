import Link from "next/link";
import { fetchCustomers } from "@/lib/api/customers";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString(undefined, {
    dateStyle: "short",
    timeStyle: "short",
  });
}

export default async function CustomersPage() {
  let customers = await fetchCustomers();

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold text-foreground">Customers</h2>
        <p className="text-sm text-muted-foreground">
          Health scores, active users, and last activity.
        </p>
      </div>
      <div className="rounded-lg border border-border bg-card">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border bg-muted/50">
              <th className="px-4 py-3 text-left font-medium text-foreground">Name</th>
              <th className="px-4 py-3 text-left font-medium text-foreground">Health score</th>
              <th className="px-4 py-3 text-left font-medium text-foreground">Active users</th>
              <th className="px-4 py-3 text-left font-medium text-foreground">Last activity</th>
            </tr>
          </thead>
          <tbody>
            {customers.length === 0 ? (
              <tr>
                <td colSpan={4} className="px-4 py-8 text-center text-muted-foreground">
                  No customers yet. Add GET /api/v1/customers on the backend to list customers.
                </td>
              </tr>
            ) : (
              customers.map((c) => (
                <tr key={c.id} className="border-b border-border last:border-0 hover:bg-muted/30">
                  <td className="px-4 py-3">
                    <Link href={`/customers/${c.id}`} className="font-medium text-primary hover:underline">
                      {c.name}
                    </Link>
                  </td>
                  <td className="px-4 py-3">{c.health_score}</td>
                  <td className="px-4 py-3">{c.active_users}</td>
                  <td className="px-4 py-3 text-muted-foreground">{formatDate(c.last_activity)}</td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
