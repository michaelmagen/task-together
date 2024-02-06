import { create } from 'zustand'
import type { List } from '../typeValidators';

type ListStoreState = {
	selectedList?: List,
}

type ListStoreActions = {
	updateSelectedList: (list: ListStoreState['selectedList']) => void
}

const useListStore = create<ListStoreState & ListStoreActions>((set) => ({
	updateSelectedList: (list) => set(() => ({ selectedList: list }))
}))

export default useListStore;
