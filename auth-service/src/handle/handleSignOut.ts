import { auth } from "..";
import type { Context } from "hono";

export async function handleSignOut(c: Context) {
  const data = await auth.api.signOut({
    headers: c.req.raw.headers,
  });
  return c.json(data);
}
