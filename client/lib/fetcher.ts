import { env } from "@/env.mjs";
import { ListID, TaskID } from "./typeValidators";

export const endpoints = {
	current_user: "/users",
	user_lists: "/lists",
	tasks: (listID: ListID) => `/lists/${listID}/tasks`,
	singleTask: (listID: ListID, taskID: TaskID) =>
		`/lists/${listID}/tasks/${taskID}`,
};

export enum Method {
	GET = "GET",
	POST = "POST",
	PUT = "PUT",
	DELETE = "DELETE",
	PATCH = "PATCH",
}

export type FetcherOptions = {
	method: Method;
	cookieString: string;
	body?: object;
};

export default async function fetcher(
	endpoint: string,
	options: FetcherOptions,
) {
	const { method, body, cookieString } = options;
	console.log(method, endpoint, body);
	const response = await fetch(env.NEXT_PUBLIC_API_URL + endpoint, {
		method: method,
		headers: {
			Cookie: cookieString,
		},
		credentials: "include",
		body: body ? JSON.stringify(body) : null,
	});
	console.log(response);

	if (!response.ok) {
		throw new Error(`(Fetcher) Failed to fetch ${endpoint}`);
	}
	const data = await response.json();
	return data;
}
