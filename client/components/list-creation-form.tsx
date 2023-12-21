"use client";
import { Button } from "@/components/ui/button"
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	DialogTrigger,
} from "@/components/ui/dialog"
import { Input } from "@/components/ui/input"
import { PlusIcon } from "lucide-react";
import useSWRMutation from 'swr/mutation'
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/components/ui/form"
import * as z from "zod"
import { zodResolver } from "@hookform/resolvers/zod"
import { useForm } from "react-hook-form"
import fetcher, { Endpoint, FetcherOptions, Method } from "@/lib/fetcher";
import useCookieStore from "@/lib/stores";
import type { List } from "@/lib/typeValidators";
import { Loader2 } from "lucide-react";
import { useState } from "react";

const formSchema = z.object({
	name: z.string().min(1, {
		message: "List name can not be empty."
	}).max(100, {
		message: "Name must be less than 100 characters."
	}),
})

export function ListCreationForm() {
	const [open, setOpen] = useState(false)
	const { cookies } = useCookieStore()
	const listMutator = (key: Endpoint, { arg }: { arg: FetcherOptions }) => fetcher(key, arg)
	const { trigger, isMutating } = useSWRMutation<List, Error, Endpoint, FetcherOptions>(Endpoint.USER_LISTS, listMutator, {
		onSuccess: () => setOpen(false)
	})

	const form = useForm<z.infer<typeof formSchema>>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: "",
		},
	})

	const onSubmitForm = (values: z.infer<typeof formSchema>) => {
		const fetcherOptions: FetcherOptions = {
			method: Method.POST,
			cookieString: cookies,
			body: values,
		}
		trigger(fetcherOptions)
	}

	return (
		<Dialog open={open}>
			<DialogTrigger asChild>
				<Button className="mt-2 w-full" onClick={() => setOpen(true)}>
					<PlusIcon className="mr-2 h-4 w-4" />
					<span>Create New List</span>
				</Button>
			</DialogTrigger>
			<DialogContent className="sm:max-w-[425px]">
				<DialogHeader>
					<DialogTitle>Create new list</DialogTitle>
					<DialogDescription>
						Enter the name for the new list you want to create.
					</DialogDescription>
				</DialogHeader>
				<Form {...form}>
					<form onSubmit={form.handleSubmit(onSubmitForm)}>
						<FormField
							control={form.control}
							name="name"
							render={({ field }) => (
								<FormItem className="grid gap-4 py-4">
									<div className="grid grid-cols-4 items-center gap-x-4 gap-y-2">

										<FormLabel className="text-right">List Name</FormLabel>
										<FormControl>
											<Input placeholder="Groceries" {...field} className="col-span-3" />
										</FormControl>
										<FormMessage className="col-span-3 col-start-2" />
									</div>
								</FormItem>
							)}
						/>
						<DialogFooter>
							<Button type="submit">{
								isMutating ? <Loader2 className="animate-spin" /> : "Create List"}</Button>
						</DialogFooter>
					</form>
				</Form>
			</DialogContent>
		</Dialog>
	)
}

