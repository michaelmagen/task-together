import useCookieStore from "@/lib/stores/cookieStore";
import { cookies } from "next/headers";
import React from "react";
import Sidebar from "@/components/sidebar";
import { MobileSidebar } from "@/components/mobile-sidebar";
import { userSchema } from "@/lib/typeValidators";
import TasksDisplay from "@/components/tasks-display";
import fetcher, { endpoints, Method } from "@/lib/fetcher";
import AppInitializer from "@/components/app-initializer";

async function getUser(cookieString: string) {
  try {
    const data = await fetcher(endpoints.current_user, {
      method: Method.GET,
      cookieString: cookieString,
    });
    const user = userSchema.parse(data);
    return user;
  } catch (error) {
    console.log(error);
    throw error;
  }
}

export default async function Home() {
  const cookieStore = cookies();
  const cookieString = cookieStore
    .getAll()
    .map((cookie) => `${cookie.name}=${cookie.value}`)
    .join("; ");

  const user = await getUser(cookieString);

  return (
    <main className="flex flex-col overflow-hidden md:flex-row">
      <AppInitializer user={user} cookies={cookieString}>
        <div className="hidden overflow-hidden md:flex">
          <Sidebar />
        </div>
        <MobileSidebar />
        <TasksDisplay />
      </AppInitializer>
    </main>
  );
}
