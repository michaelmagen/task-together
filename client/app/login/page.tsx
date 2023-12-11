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

    url.search = parameters.toString();
    const redirectUrl = url.toString();

    return redirectUrl;
  };
  return (
    <main className="flex justify-center items-center mt-20">
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
