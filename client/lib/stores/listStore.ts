import { create } from 'zustand'
import type { ListID } from '../typeValidators';

type ListStoreState = {
	listID: ListID | null,
}

type ListStoreActions = {
	updateListID: (listID: ListStoreState['listID']) => void
}

const useListStore = create<ListStoreState & ListStoreActions>((set) => ({
	listID: null,
	updateListID: (listID) => set(() => ({ listID: listID }))
}))

export default useListStore;
