import useCookieStore from "@/lib/stores";
import { cookies } from "next/headers";
import React from "react";
import Sidebar from "@/components/sidebar"


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
			<Sidebar />
		</main>
	);
}
