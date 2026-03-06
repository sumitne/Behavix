import { fetchInsights } from "@/lib/api/insights";
import { InsightCard } from "@/components/insight/insight-card";
import { ApiError } from "@/lib/api/client";

export default async function InsightsFeedPage() {
  let insights: Awaited<ReturnType<typeof fetchInsights>> = [];
  let error: string | null = null;

  try {
    insights = await fetchInsights(50);
  } catch (e) {
    if (e instanceof ApiError) {
      if (e.status === 401) {
        error = "Invalid or missing API key. Set BEHAVIX_API_KEY in .env.local.";
      } else {
        error = `Failed to load insights: ${e.status} ${e.message}`;
      }
    } else {
      error = e instanceof Error ? e.message : "Failed to load insights.";
    }
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
