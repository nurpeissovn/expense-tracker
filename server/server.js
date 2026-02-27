import express from "express";
import cors from "cors";
import dotenv from "dotenv";
import { Pool } from "pg";
import path from "path";
import fs from "fs/promises";
import { fileURLToPath } from "url";

dotenv.config();

const app = express();
const port = process.env.PORT || 8080;
const allowOrigin = process.env.ALLOW_ORIGIN || "*";

const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
  ssl: process.env.DATABASE_URL?.includes("railway") ? { rejectUnauthorized: false } : false,
});

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const publicDir = path.join(__dirname, "public");
const schemaPath = path.join(__dirname, "schema.sql");

async function ensureDatabase() {
  const schema = await fs.readFile(schemaPath, "utf8");
  await pool.query(schema);
  console.log("Database schema ensured");
}

app.use(
  cors({
    origin: allowOrigin === "*" ? true : allowOrigin.split(","),
  })
);
app.use(express.json());
app.use(express.static(publicDir));

app.get("/api/health", async (_req, res) => {
  try {
    const result = await pool.query("select now()");
    res.json({ status: "ok", time: result.rows[0].now });
  } catch (err) {
    res.status(500).json({ status: "error", message: err.message });
  }
});

app.get("/api/transactions", async (_req, res) => {
  try {
    const result = await pool.query(
      "select id, type, amount, category, note, to_char(date, 'YYYY-MM-DD') as date, created_at from transactions order by date desc, created_at desc"
    );
    res.json(result.rows);
  } catch (err) {
    res.status(500).json({ error: "Failed to fetch transactions" });
  }
});

app.post("/api/transactions", async (req, res) => {
  const { type, amount, category, date, note } = req.body || {};

  if (!["income", "expense"].includes(type)) {
    return res.status(400).json({ error: "type must be income or expense" });
  }
  const numericAmount = Number(amount);
  if (!Number.isFinite(numericAmount) || numericAmount < 0) {
    return res.status(400).json({ error: "amount must be a positive number" });
  }
  if (!category || !date) {
    return res.status(400).json({ error: "category and date are required" });
  }

  try {
    const result = await pool.query(
      "insert into transactions (type, amount, category, note, date) values ($1, $2, $3, $4, $5) returning id, type, amount, category, note, to_char(date, 'YYYY-MM-DD') as date, created_at",
      [type, numericAmount, category, note || null, date]
    );
    res.status(201).json(result.rows[0]);
  } catch (err) {
    console.error("create transaction failed:", err);
    res.status(500).json({ error: err.message || "Failed to create transaction" });
  }
});

app.delete("/api/transactions/:id", async (req, res) => {
  const { id } = req.params;
  if (!id) return res.status(400).json({ error: "id required" });
  try {
    const result = await pool.query("delete from transactions where id = $1 returning id", [id]);
    if (!result.rowCount) return res.status(404).json({ error: "Not found" });
    res.json({ ok: true });
  } catch (err) {
    console.error("delete transaction failed:", err);
    res.status(500).json({ error: err.message || "Failed to delete transaction" });
  }
});

app.get("*", (_req, res) => {
  res.sendFile(path.join(publicDir, "index.html"));
});

(async () => {
  try {
    await ensureDatabase();
    app.listen(port, () => {
      console.log(`Expense Tracker API running on :${port}`);
    });
  } catch (err) {
    console.error("Failed to initialize database:", err);
    process.exit(1);
  }
})();
