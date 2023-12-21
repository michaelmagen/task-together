import { env } from "@/env.mjs";

export enum Endpoint {
	CURRENT_USER = '/users',
	USER_LISTS = '/lists',
}

export enum Method {
	GET = 'GET',
	POST = 'POST',
	PUT = 'PUT',
	DELETE = 'DELETE',
	PATCH = 'PATCH',
}

export type FetcherOptions = {
	method: Method
	cookieString: string
	body?: object
}

export default async function fetcher(endpoint: Endpoint, options: FetcherOptions) {
	const { method, body, cookieString } = options
	console.log(method, cookieString)
	const response = await fetch(env.NEXT_PUBLIC_API_URL + endpoint, {
		method: method,
		headers: {
			Cookie: cookieString,
		},
		credentials: "include",
		body: body ? JSON.stringify(body) : null
	});


	if (!response.ok) {
		throw new Error(`(Fetcher) Failed to get ${endpoint}`);
	}
	const data = await response.json();
	return data;
}

