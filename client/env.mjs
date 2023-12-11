import { createEnv } from "@t3-oss/env-nextjs";
import { z } from "zod";

export const env = createEnv({
  client: {
    NEXT_PUBLIC_CLIENT_ID : z.string().min(1),
    NEXT_PUBLIC_AUTH_SCOPES: z.string().min(1),
    NEXT_PUBLIC_REDIRECT_URI: z.string().url(),
    NEXT_PUBLIC_AUTH_STATE: z.string().min(1),
  },

   // For Next.js >= 13.4.4, you only need to destructure client variables:
  experimental__runtimeEnv: {
    NEXT_PUBLIC_CLIENT_ID : process.env.NEXT_PUBLIC_CLIENT_ID,  
    NEXT_PUBLIC_AUTH_SCOPES: process.env.NEXT_PUBLIC_AUTH_SCOPES,
    NEXT_PUBLIC_REDIRECT_URI: process.env.NEXT_PUBLIC_REDIRECT_URI, 
    NEXT_PUBLIC_AUTH_STATE: process.env.NEXT_PUBLIC_AUTH_STATE, 
  },
  
})