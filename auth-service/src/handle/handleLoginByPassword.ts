import { auth } from "..";
import type { Context } from "hono";
type LoginByPasswordDto = {
  email: string;
  password: string;
};

export async function handleLoginByPassword(c: Context) {
  const { email, password }: LoginByPasswordDto = await c.req.json();
  const data = await auth.api.signInEmail({
    body: {
      email,
      password,
    },
    headers: c.req.raw.headers,
  });
  return c.json(data);
}
