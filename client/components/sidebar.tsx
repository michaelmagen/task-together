import ListSelection from "@/components/list-selection";
import { ListCreationForm } from "@/components/list-creation-form";

export default async function Sidebar() {
	// TODO: Make w 1/5 when there is also a small view port sidebar
	return (
		<div className="m-4 flex w-full max-w-xs flex-col items-center justify-center rounded-lg border-2 p-2">
			<div className="text-lg font-bold">Select List</div>
			<ListSelection />
			<ListCreationForm />
		</div>
	);
}
