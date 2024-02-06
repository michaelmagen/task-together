"use client";
import useUserStore from "@/lib/stores/userStore";
import { User } from "@/lib/typeValidators";
import { useEffect, type PropsWithChildren } from "react";

interface AppInitializerProps extends PropsWithChildren {
  cookies: string;
  user: User;
}

// Initializes the app with starting info like user info, cookie strings, etc.
export default function AppInitializer({
  cookies,
  user,
  children,
}: AppInitializerProps) {
  const { updateUser } = useUserStore();
  useEffect(() => {
    updateUser(user);
  }, [updateUser]);
  return <>{children}</>;
}
