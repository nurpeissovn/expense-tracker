create extension if not exists "pgcrypto";

create table if not exists transactions (
  id uuid primary key default gen_random_uuid(),
  type text not null check (type in ('income', 'expense')),
  amount numeric(12,2) not null check (amount >= 0),
  category text not null,
  note text,
  date date not null,
  created_at timestamptz not null default now()
);

create index if not exists idx_transactions_type on transactions(type);
create index if not exists idx_transactions_date on transactions(date desc);
