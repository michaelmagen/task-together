import { type User } from "@/lib/typeValidators";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import {
	HoverCard,
	HoverCardContent,
	HoverCardTrigger,
} from "@/components/ui/hover-card";

interface UserCardProps {
	user: User;
}
export default function UserCard({ user }: UserCardProps) {
	const userInitials = user.given_name.charAt(0) + user.family_name.charAt(0);
	return (
		<HoverCard>
			<HoverCardTrigger asChild>
				<Button variant="link">{user.name}</Button>
			</HoverCardTrigger>
			<HoverCardContent className="w-auto">
				<div className="flex justify-between space-x-4 items-center">
					<Avatar>
						<AvatarImage
							src={user.picture}
							alt={`${user.name} picture`}
							crossOrigin="anonymous"
							referrerPolicy="no-referrer"
						/>
						<AvatarFallback>{userInitials}</AvatarFallback>
					</Avatar>
					<div className="space-y-1">
						<h4 className="text-sm font-semibold">{user.name}</h4>
						<p className="text-sm">{user.email}</p>
					</div>
				</div>
			</HoverCardContent>
		</HoverCard>
	);
}
