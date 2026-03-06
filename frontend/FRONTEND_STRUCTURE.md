# Behavix-AI Frontend — Project Structure

Phase 1 minimal internal dashboard. Next.js (App Router), React, Tailwind CSS, shadcn/ui, TypeScript.

## Directory layout

```
frontend/
├── public/                    # Static assets
├── src/
│   ├── app/                   # App Router
│   │   ├── (dashboard)/       # Dashboard route group (shared layout)
│   │   │   ├── layout.tsx     # Sidebar + header layout
│   │   │   ├── page.tsx       # Insights Feed (main)
│   │   │   └── customers/
│   │   │       ├── page.tsx   # Customers list
│   │   │       └── [id]/
│   │   │           └── page.tsx  # Customer detail
│   │   ├── layout.tsx        # Root layout (fonts, TooltipProvider)
│   │   └── globals.css       # Tailwind + shadcn theme
│   │
│   ├── components/
│   │   ├── ui/                # shadcn components (card, badge, button, sidebar, etc.)
│   │   └── layout/            # App-specific layout
│   │       ├── app-sidebar.tsx
│   │       └── dashboard-header.tsx
│   │
│   ├── lib/
│   │   ├── utils.ts          # cn() etc.
│   │   └── api/              # Backend API client (server-side fetch)
│   │       ├── config.ts     # API_BASE_URL, BEHAVIX_API_KEY
│   │       ├── client.ts     # apiFetch(), ApiError
│   │       ├── insights.ts   # fetchInsights()
│   │       └── customers.ts  # fetchCustomers(), fetchCustomer(), fetchCustomerUsage()
│   │
│   ├── types/                 # Shared types (align with backend API)
│   │   ├── insight.ts
│   │   └── customer.ts
│   │
│   ├── data/                  # Mock data (Phase 1)
│   │   ├── mockInsights.ts
│   │   └── mockCustomers.ts   # For customers page
│   │
│   └── hooks/                 # shadcn hooks (e.g. use-mobile)
│
├── components.json            # shadcn config
├── package.json
├── tsconfig.json
├── tailwind.config.ts
└── next.config.ts
```

## Layout system

- **Root layout** (`src/app/layout.tsx`): HTML shell, fonts, `TooltipProvider` for shadcn.
- **Dashboard layout** (`src/app/(dashboard)/layout.tsx`): Wraps all dashboard pages with:
  - **Sidebar** (left): Nav links — Insights, Customers. Uses shadcn `Sidebar*` components.
  - **Header** (top): Logo/title, optional actions.
  - **Main**: `SidebarInset` for page content.

## Pages

| Route              | Page            | Purpose                          |
|--------------------|-----------------|----------------------------------|
| `/`                | Insights Feed   | List of AI-generated insights   |
| `/customers`       | Customers       | Table of customers               |
| `/customers/[id]`  | Customer Detail | Overview, usage, insights        |

## Data flow

- **Real API**: Server-side fetch in `src/lib/api/` calls the Go backend.
- **Config**: Set `API_BASE_URL` and `BEHAVIX_API_KEY` in `.env.local` (see `frontend/.env.example`).
- **Insights**: `GET /api/v1/insights` with `Authorization: Bearer <BEHAVIX_API_KEY>`.
- **Customers**: `GET /api/v1/customers`, `GET /api/v1/customers/:id`, `GET /api/v1/customers/:id/usage`. If the backend does not implement these yet, the UI shows empty state or “not available” messages.

## Backend API alignment

- **GET /insights** → `Insight[]` (id, type, title, description, severity, created_at, …).
- **GET /customers** → `Customer[]`.
- **GET /customers/{id}** → single customer.
- **GET /customers/{id}/usage** → usage metrics.

Types in `src/types/` mirror these responses for type-safe integration when wiring to the API.
