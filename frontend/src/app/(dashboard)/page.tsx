"use client";

import { useEffect, useState } from "react";
import type { Insight } from "@/types/insight";
import { InsightCard } from "@/components/insight/insight-card";
import { clientFetch } from "@/lib/api/client-browser";

export default function InsightsFeedPage() {
  const [insights, setInsights] = useState<Insight[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setError(null);
    clientFetch<Insight[]>("/api/insights?limit=50")
      .then((data) => {
        if (!cancelled) {
          setInsights(Array.isArray(data) ? data : []);
        }
      })
      .catch((e) => {
        if (!cancelled) {
          const status = (e as { status?: number }).status;
          if (status === 401) {
            setError("Invalid or missing API key. Set BEHAVIX_API_KEY in .env.local.");
          } else {
            setError(e instanceof Error ? e.message : "Failed to load insights.");
          }
        }
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  if (loading) {
    return (
      <div className="space-y-6">
        <div>
          <h2 className="text-lg font-semibold text-foreground">Insights Feed</h2>
          <p className="text-sm text-muted-foreground">
            AI-generated behavioral insights from product usage.
          </p>
        </div>
        <p className="text-sm text-muted-foreground">Loading…</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="space-y-6">
        <div>
          <h2 className="text-lg font-semibold text-foreground">Insights Feed</h2>
          <p className="text-sm text-muted-foreground">
            AI-generated behavioral insights from product usage.
          </p>
        </div>
        <div className="rounded-lg border border-destructive/50 bg-destructive/10 p-4 text-sm text-destructive">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-lg font-semibold text-foreground">Insights Feed</h2>
        <p className="text-sm text-muted-foreground">
          AI-generated behavioral insights from product usage.
        </p>
      </div>
      {insights.length === 0 ? (
        <p className="text-sm text-muted-foreground">No insights yet.</p>
      ) : (
        <ul className="grid gap-4 sm:grid-cols-1 lg:grid-cols-2">
          {insights.map((insight) => (
            <li key={insight.id}>
              <InsightCard insight={insight} />
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
