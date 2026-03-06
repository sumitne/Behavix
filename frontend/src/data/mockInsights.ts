import type { Insight } from "@/types/insight";

export const mockInsights: Insight[] = [
  {
    id: "ins-1",
    type: "usage_drop",
    title: "Weekly active users dropped 23%",
    description:
      "Acme Corp's weekly active users fell from 412 to 317 over the last 14 days. Largest drop in the Reports and Export modules.",
    severity: "high",
    created_at: "2025-03-06T10:30:00Z",
    metric_value: 317,
    baseline_value: 412,
    customer_name: "Acme Corp",
  },
  {
    id: "ins-2",
    type: "feature_adoption",
    title: "New dashboard rarely used",
    description:
      "Only 12% of Beta Inc users have used the new Analytics Dashboard in the last 30 days. Consider in-app prompts or email campaign.",
    severity: "medium",
    created_at: "2025-03-05T14:00:00Z",
    metric_value: 12,
    baseline_value: 45,
    customer_name: "Beta Inc",
  },
  {
    id: "ins-3",
    type: "engagement_change",
    title: "Session depth increased",
    description:
      "Gamma LLC shows a 18% increase in average events per session compared to last month. Engagement trend is positive.",
    severity: "low",
    created_at: "2025-03-04T09:15:00Z",
    metric_value: 8.2,
    baseline_value: 6.9,
    customer_name: "Gamma LLC",
  },
  {
    id: "ins-4",
    type: "usage_drop",
    title: "Login frequency decline",
    description:
      "Delta Co daily logins are down 15% week-over-week. Possible churn risk; recommend reaching out to power users.",
    severity: "medium",
    created_at: "2025-03-03T16:45:00Z",
    customer_name: "Delta Co",
  },
  {
    id: "ins-5",
    type: "feature_adoption",
    title: "API usage growing",
    description:
      "Epsilon Ltd API calls up 40% in the last 7 days. Indicates successful developer adoption of the new API.",
    severity: "low",
    created_at: "2025-03-02T11:00:00Z",
    customer_name: "Epsilon Ltd",
  },
];
