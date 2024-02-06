import { create } from "zustand";
import type { User } from "../typeValidators";

// Define the state shape for the store
interface UserStoreState {
	user: User | null;
	updateUser: (user: User) => void;
}

// Create the Zustand store
const useUserStore = create<UserStoreState>((set) => ({
	user: null,
	updateUser: (user) => set(() => ({ user: user })),
}));

export default useUserStore;
