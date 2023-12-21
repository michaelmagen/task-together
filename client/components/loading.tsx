import { Loader2 } from "lucide-react";
export default function Loading() {
	return (
		<div className="text-muted-foreground flex items-center justify-center">
			<div className="m-5 flex items-center gap-1">
				<Loader2 className="h-6 w-6 animate-spin" />
				<div className="text-sm">Loading...</div>
			</div>
		</div>
	)
}
