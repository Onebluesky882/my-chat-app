import { Hono } from "hono";
import { handleSignUpByPassword } from "../handle/handleSignUpByPassword";
import { handleLoginByPassword } from "../handle/handleLoginByPassword";
import { handleSignOut } from "../handle/handleSignOut";

const authRouter = new Hono();

authRouter.post("/sign-up/email", handleSignUpByPassword);
authRouter.post("/sign-in/email", handleLoginByPassword);
authRouter.post("/sign-out", handleSignOut);

export default authRouter;
