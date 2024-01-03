import useCookieStore from "@/lib/stores/cookieStore";
import { cookies } from "next/headers";
import React from "react";
import Sidebar from "@/components/sidebar"
import { MobileSidebar } from "@/components/mobile-sidebar";
import TasksDisplay from "@/components/tasks-display";


export default async function Home() {
	const cookieStore = cookies();
	const cookieString = cookieStore
		.getAll()
		.map((cookie) => `${cookie.name}=${cookie.value}`)
		.join("; ");
	// Store the cookie strings in global store so can be accessed by client components
	useCookieStore.setState({
		cookies: cookieString
	})

	return (
		<main className="flex overflow-hidden">
			<div className="hidden overflow-hidden md:flex">
				<Sidebar />
			</div>
			<MobileSidebar />
			<TasksDisplay />
		</main>
	);
}
