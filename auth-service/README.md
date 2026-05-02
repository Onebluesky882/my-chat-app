add better auth schema

path : auth/auth.ts

```ts
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
```

cli better auth schema : bunx @better-auth/cli generate --config=./src/auth/auth.ts
