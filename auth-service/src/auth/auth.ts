import { drizzleAdapter } from "@better-auth/drizzle-adapter";
import { betterAuth } from "better-auth";

export const auth = betterAuth({
  database: drizzleAdapter(
    {},
    {
      provider: "pg",
    },
  ),
});
