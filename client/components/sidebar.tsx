import ListSelection from "@/components/list-selection";
import { ListCreationForm } from "@/components/list-creation-form";
import { Method, Endpoint } from "@/lib/fetcher";
import fetcher from "@/lib/fetcher";
import { userSchema } from "@/lib/typeValidators";
import UserCard from "./user-card";
import { cookies } from "next/headers";

async function getUser(cookieString: string) {
	try {
		const data = await fetcher(Endpoint.CURRENT_USER, {
			method: Method.GET,
			cookieString: cookieString,
		})
		const user = userSchema.parse(data);
		return user;
	} catch (error) {
		console.log(error);
		throw error;
	}
}

export default async function Sidebar() {
	const cookieStore = cookies();
	const cookieString = cookieStore
		.getAll()
		.map((cookie) => `${cookie.name}=${cookie.value}`)
		.join("; ");
	const user = await getUser(cookieString);

	// TODO: Make w 1/5 when there is also a small view port sidebar
	return (
		<div className="m-4 flex w-full max-w-xs flex-col items-center justify-center rounded-lg border-2 p-2">
			<div className="text-lg font-bold">Select List</div>
			<ListSelection />
			<ListCreationForm />
			<UserCard user={user} />
		</div>
	)
}
