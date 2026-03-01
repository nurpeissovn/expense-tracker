# FinSet — Finance Dashboard

A full-stack personal finance dashboard built with **Go + PostgreSQL** backend and a vanilla HTML/CSS/JS frontend. Deployable to **Railway** in minutes.

---

## Project Structure

```
finset/
├── cmd/server/
│   └── main.go              # Entry point, router, static file serving
├── internal/
│   ├── db/
│   │   └── db.go            # PostgreSQL connection pool + queries + migrations
│   ├── handlers/
│   │   └── handlers.go      # HTTP handlers for all API routes
│   └── models/
│       └── transaction.go   # Transaction model + validation
├── static/
│   └── index.html           # Full frontend (HTML + CSS + JS in one file)
├── Dockerfile               # Multi-stage Docker build
├── railway.json             # Railway deployment config
├── go.mod / go.sum
└── .env.example             # Local dev env template
```

---

## Database — Safe Migration Policy

> **Your existing data is NEVER deleted or modified during deployment.**

On every startup the server runs `Migrate()` which uses:
- `CREATE TABLE IF NOT EXISTS` — only creates the table if it doesn't exist
- `ALTER TABLE … ADD COLUMN IF NOT EXISTS` — only adds columns missing from an older schema
- `CREATE INDEX IF NOT EXISTS` — only creates indexes if they don't already exist

This means you can deploy new versions repeatedly without losing any records.

### Schema

```sql
CREATE TABLE transactions (
    id          TEXT           PRIMARY KEY,
    type        TEXT           NOT NULL CHECK (type IN ('income','expense')),
    amount      NUMERIC(12,2)  NOT NULL CHECK (amount > 0),
    category    TEXT           NOT NULL DEFAULT '',
    method      TEXT           NOT NULL DEFAULT 'Cash',
    date        DATE           NOT NULL,
    note        TEXT           NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);
```

---

## REST API

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/health` | Health check (DB ping) |
| GET | `/api/transactions` | List all transactions (date DESC) |
| POST | `/api/transactions` | Create a transaction |
| GET | `/api/transactions/:id` | Get one transaction |
| DELETE | `/api/transactions/:id` | Delete a transaction |
| GET | `/api/stats` | Total income, expense, balance, count |
| GET | `/api/monthly-flow?months=7` | Per-month income/expense |
| GET | `/api/category-breakdown` | Expense totals grouped by category |
| POST | `/api/import` | Bulk import `{transactions:[…]}` (skips duplicates) |

### POST /api/transactions body

```json
{
  "type": "expense",
  "amount": 42.50,
  "category": "Food",
  "method": "Card",
  "date": "2025-01-15",
  "note": "Lunch"
}
```

---

## Deploy to Railway

### Step 1 — Push code to GitHub

```bash
cd finset
git init
git add .
git commit -m "initial commit"
gh repo create finset --public --push
# or:
git remote add origin https://github.com/YOUR_USER/finset.git
git push -u origin main
```

### Step 2 — Create Railway project

1. Go to [railway.app](https://railway.app) → **New Project**
2. Click **Deploy from GitHub repo** → select your repo
3. Railway will detect the `Dockerfile` automatically

### Step 3 — Link your existing PostgreSQL database

You already have a PostgreSQL service on Railway. Link it:

1. In your Railway project, click on your **app service** (not the DB)
2. Go to **Variables** tab
3. Click **+ Add Variable Reference**
4. Select your PostgreSQL service → choose `DATABASE_URL`
5. This injects the correct connection string automatically

> ⚠️ **Important:** Use the `DATABASE_URL` variable reference (not a hard-coded value).  
> Railway rotates credentials — the reference always stays current.

### Step 4 — Deploy

Railway auto-deploys on every `git push`. Watch the build logs in the Railway dashboard.

The server will:
1. Connect to PostgreSQL (retries 10× with 3s delay if DB isn't ready yet)
2. Run safe migration (creates table/indexes only if missing)
3. Start serving on the PORT Railway provides

### Step 5 — Open the app

Click **Generate Domain** in your service settings → your app is live.

---

## Local Development

### Prerequisites
- Go 1.22+
- PostgreSQL running locally (or use Railway's DB with a local tunnel)

```bash
# 1. Copy env file
cp .env.example .env
# Edit .env with your local DATABASE_URL

# 2. Download dependencies
go mod download

# 3. Run
go run ./cmd/server

# App is at http://localhost:8080
```

### Run with Railway DB locally (tunnel)

```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and link project
railway login
railway link

# Start with Railway env vars injected
railway run go run ./cmd/server
```

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| `DATABASE_URL is not set` | Add the variable reference in Railway Variables tab |
| `could not connect after 10 attempts` | Check DB service is running; check DATABASE_URL value |
| `relation "transactions" already exists` | Safe — migration skips it, your data is intact |
| Build fails | Check Go 1.22 syntax; run `go mod tidy` locally |
| 404 on `/api/*` | Check `railway.json` healthcheckPath; ensure Dockerfile CMD is correct |

---

## Environment Variables Reference

| Variable | Required | Description |
|----------|----------|-------------|
| `DATABASE_URL` | ✅ Yes | PostgreSQL connection string |
| `PORT` | Auto (Railway) | HTTP port — Railway sets this automatically |
