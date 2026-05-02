import { auth } from "..";
import type { Context } from "hono";
type SignUpByPasswordDto = {
  name: string;
  email: string;
  password: string;
  image?: string;
  callbackURL?: string;
};

export async function handleSignUpByPassword(c: Context) {
  const { name, email, password }: SignUpByPasswordDto = await c.req.json();
  const data = await auth.api.signUpEmail({
    body: {
      name,
      email,
      password,
    },
    headers: c.req.raw.headers,
  });
  return c.json(data);
}
