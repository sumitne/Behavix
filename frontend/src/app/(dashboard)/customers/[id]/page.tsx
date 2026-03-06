import Link from "next/link";
import { notFound } from "next/navigation";
import { fetchCustomer, fetchCustomerUsage } from "@/lib/api/customers";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ArrowLeftIcon } from "lucide-react";

interface CustomerDetailPageProps {
  params: Promise<{ id: string }>;
}

export default async function CustomerDetailPage({ params }: CustomerDetailPageProps) {
  const { id } = await params;
  const [customer, usage] = await Promise.all([
    fetchCustomer(id),
    fetchCustomerUsage(id),
  ]);
  if (!customer) notFound();

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
              Usage data not available. Backend GET /api/v1/customers/{id}/usage can provide this.
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
