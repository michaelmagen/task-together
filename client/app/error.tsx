"use client"; // Error components must be Client Components

import { useEffect } from "react";
import { Button } from "@/components/ui/button";

export default function Error({
	error,
	reset,
}: {
	error: Error & { digest?: string };
	reset: () => void;
}) {
	useEffect(() => {
		console.error(error);
	}, [error]);

	return (
		<div className="m-12 flex flex-col items-center space-y-3">
			<h3 className="scroll-m-20 text-xl font-semibold tracking-tight">
				Unable to load!
			</h3>
			<Button onClick={() => reset()}>Try again</Button>
		</div>
	);
}
