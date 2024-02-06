import { create } from "zustand";

type CookieStore = {
	cookies: string;
	updateCookies: (cookies: string) => void;
};

const useCookieStore = create<CookieStore>((set) => ({
	cookies: "",
	updateCookies: (cookies) => set(() => ({ cookies: cookies })),
}));

export default useCookieStore;
