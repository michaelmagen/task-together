import { create } from 'zustand'

type CookieStore = {
	cookies: string,
}

const useCookieStore = create<CookieStore>(() => ({
	cookies: "",
}))

export default useCookieStore;
