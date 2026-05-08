import { Hono } from "hono";
import { cors } from "hono/cors";
import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from "./db";
import { schema } from "./db/schema";
import authRouter from "./router/better-auth";
import { expo } from "@better-auth/expo";
const app = new Hono();
app.use(
  "/api/*",
  cors({
    origin: "http://localhost:3000",
    allowMethods: ["GET", "POST", "OPTIONS"],
    allowHeaders: ["Content-Type", "Authorization"],
    credentials: true,
  }),
);

export const auth = betterAuth({
  database: drizzleAdapter(db, {
    provider: "pg",
    schema,
  }),
  trustedOrigins: ["localhost://3000", "chatapp://", "http://192.168.1.59"],
  plugins: [expo()],
  emailAndPassword: {
    enabled: true,
  },
});
app.get("/", (c) => {
  return c.text("hello");
});
app.all("/api/auth/*", (c) => auth.handler(c.req.raw));

app.route("/api", authRouter);

const port = Number(process.env.PORT ?? 3000);

export default {
  port,
  fetch: app.fetch,
};
