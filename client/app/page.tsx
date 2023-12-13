import { cookies } from "next/headers";
import { userSchema } from "@/lib/typeValidators";
import UserCard from "@/components/user-card";
import React from "react";
import { env } from "@/env.mjs";

async function getUser() {
	try {
		const cookieStore = cookies();
		const cookiesString = cookieStore
			.getAll()
			.map((cookie) => `${cookie.name}=${cookie.value}`)
			.join("; ");
		const response = await fetch(env.NEXT_PUBLIC_API_URL + "/users/", {
			method: "GET",
			headers: {
				Cookie: cookiesString,
			},
		});

		if (!response.ok) {
			throw new Error("Failed to get user");
		}

		const data = await response.json();
		const user = userSchema.parse(data);
		return user;
	} catch (error) {
		console.log(error);
		throw error;
	}
}

export default async function Home() {
	const user = await getUser();
	return (
		<main>
			<UserCard user={user} />
		</main>
	);
}
