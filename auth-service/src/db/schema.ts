import * as authSchema from "./auth-schema";

export const schema = {
  user: authSchema.user,
  session: authSchema.session,
  account: authSchema.account,
  verification: authSchema.verification,
};

export * from "./auth-schema";
