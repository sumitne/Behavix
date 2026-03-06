"use client";

import { SidebarTrigger } from "@/components/ui/sidebar";
import { Separator } from "@/components/ui/separator";

interface DashboardHeaderProps {
  title?: string;
  description?: string;
}

export function DashboardHeader({ title, description }: DashboardHeaderProps) {
  return (
    <header className="flex h-14 shrink-0 items-center gap-2 border-b border-border bg-background px-4">
      <SidebarTrigger />
      <Separator orientation="vertical" className="h-6" />
      <div className="flex flex-1 flex-col gap-0.5 md:flex-row md:items-center md:gap-2">
        {title && <h1 className="font-semibold text-foreground">{title}</h1>}
        {description && (
          <p className="text-sm text-muted-foreground hidden sm:block">{description}</p>
        )}
      </div>
    </header>
  );
}
