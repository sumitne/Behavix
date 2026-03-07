"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import type { Customer, CustomerUsage } from "@/types/customer";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ArrowLeftIcon } from "lucide-react";
import { clientFetch } from "@/lib/api/client-browser";

export default function CustomerDetailPage() {
  const params = useParams();
  const id = typeof params?.id === "string" ? params.id : "";
  const [customer, setCustomer] = useState<Customer | null>(null);
  const [usage, setUsage] = useState<CustomerUsage | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    if (!id) {
      setLoading(false);
      setNotFound(true);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setError(null);
    setNotFound(false);
    clientFetch<Customer>(`/api/customers/${encodeURIComponent(id)}`)
      .then((c) => {
        if (cancelled) return;
        setCustomer(c);
        return clientFetch<CustomerUsage>(`/api/customers/${encodeURIComponent(id)}/usage`).then((u) => {
          if (!cancelled) setUsage(u);
        });
      })
      .catch((e) => {
        if (!cancelled) {
          const status = (e as { status?: number }).status;
          if (status === 404) setNotFound(true);
          else setError(e instanceof Error ? e.message : "Failed to load customer.");
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, [id]);

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/customers" aria-label="Back to customers">
              <ArrowLeftIcon className="h-4 w-4" />
            </Link>
          </Button>
          <p className="text-sm text-muted-foreground">Loading…</p>
        </div>
      </div>
    );
  }

  if (notFound || !customer) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/customers" aria-label="Back to customers">
              <ArrowLeftIcon className="h-4 w-4" />
            </Link>
          </Button>
          <p className="text-sm text-muted-foreground">Customer not found.</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link href="/customers" aria-label="Back to customers">
              <ArrowLeftIcon className="h-4 w-4" />
            </Link>
          </Button>
          <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-sm text-destructive">
            {error}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link href="/customers" aria-label="Back to customers">
            <ArrowLeftIcon className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h2 className="text-lg font-semibold text-foreground">{customer.name}</h2>
          <p className="text-sm text-muted-foreground">Customer overview</p>
        </div>
      </div>
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader className="pb-2">
            <p className="text-sm font-medium text-muted-foreground">Health score</p>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-semibold">{customer.health_score}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <p className="text-sm font-medium text-muted-foreground">Active users</p>
          </CardHeader>
          <CardContent>
            <p className="text-2xl font-semibold">{customer.active_users}</p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <p className="text-sm font-medium text-muted-foreground">Last activity</p>
          </CardHeader>
          <CardContent>
            <p className="text-sm">{new Date(customer.last_activity).toLocaleString()}</p>
          </CardContent>
        </Card>
      </div>
      <Card>
        <CardHeader>
          <p className="text-sm font-medium text-muted-foreground">Usage & feature adoption</p>
        </CardHeader>
        <CardContent>
          {usage ? (
            <div className="space-y-2 text-sm">
              <p>Period: {usage.period}</p>
              <p>Active users: {usage.active_users}</p>
              <p>Events: {usage.events_count}</p>
              {usage.feature_adoption && Object.keys(usage.feature_adoption).length > 0 && (
                <div className="mt-2">
                  <p className="font-medium text-foreground">Feature adoption</p>
                  <ul className="list-inside list-disc text-muted-foreground">
                    {Object.entries(usage.feature_adoption).map(([feature, pct]) => (
                      <li key={feature}>{feature}: {pct}%</li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          ) : (
            <p className="text-sm text-muted-foreground">
              Usage data not available.
            </p>
          )}
        </CardContent>
      </Card>
      <Card>
        <CardHeader>
          <p className="text-sm font-medium text-muted-foreground">Recent insights</p>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            Recent insights for this customer can be shown here when the API supports filtering by customer/tenant.
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
