import { Button } from "@/components/ui/button";
import { env } from "../../env.mjs";
import Link from "next/link";
import GoogleIcon from "@/components/google-icon";

export default function Home() {
  const authURL = () => {
    const url = new URL("https://accounts.google.com/o/oauth2/auth");
    const parameters = new URLSearchParams();
    parameters.append("client_id", env.NEXT_PUBLIC_CLIENT_ID);
    parameters.append("scope", env.NEXT_PUBLIC_AUTH_SCOPES);
    parameters.append("redirect_uri", env.NEXT_PUBLIC_REDIRECT_URI);
    parameters.append("response_type", "code");
    parameters.append("state", env.NEXT_PUBLIC_AUTH_STATE);
    // The following shows the google prompt every time so that we can get the refresh token every time
    parameters.append("access_type", "offline");
    parameters.append("prompt", "consent");

    url.search = parameters.toString();
    const redirectUrl = url.toString();

    return redirectUrl;
  };

  return (
    <main className="mt-20 flex items-center justify-center">
      <Button asChild>
        <Link href={authURL()}>
          <div className="flex items-center space-x-2">
            <GoogleIcon />
            <span>Sign in with Google</span>
          </div>
        </Link>
      </Button>
    </main>
  );
}
