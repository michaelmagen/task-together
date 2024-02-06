"use client";

import useSWR, { Fetcher } from "swr";
import fetcher, { FetcherOptions, Method, endpoints } from "@/lib/fetcher";
import { List } from "@/lib/typeValidators";
import { Button } from "./ui/button";
import { cn } from "@/lib/utils";
import useListStore from "@/lib/stores/listStore";
import useCookieStore from "@/lib/stores/cookieStore";
import Loading from "./loading";
import { Separator } from "@/components/ui/separator"
import { ListChecks } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area"

export default function ListSelection() {
	const { selectedList, updateSelectedList } = useListStore()
	const { cookies } = useCookieStore()
	const fetcherOptions: FetcherOptions = {
		method: Method.GET,
		cookieString: cookies
	}

	const listFetcher: Fetcher<List[], string> = (endpoint) => fetcher(endpoint, fetcherOptions)
	const { data, error, isLoading } = useSWR(endpoints.user_lists, listFetcher)

	if (isLoading) {
		return <Loading />;
	}

	if (error || !data) {
		return <div>Failed to load lists</div>;
	}

	if (data && data.length === 0) {
		return <div>Create your first list!</div>
	}

	return (
		<ScrollArea className="w-full">
			{data.map((list) => (
				<div key={list.list_id}>
					<Button
						variant="ghost"
						className={cn(
							selectedList?.list_id === list.list_id
								? "bg-muted hover:bg-muted"
								: "hover:bg-transparent hover:underline",
							"w-full justify-start"
						)}
						onClick={() => updateSelectedList(list)}
					>
						<ListChecks className="mr-2" />
						{list.name}
					</Button>
					<Separator className="my-2" />
				</div>
			))}
		</ScrollArea>
	);
}

