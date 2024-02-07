import { create } from "zustand";
import type { Task, ListID, TaskID } from "../typeValidators";

interface ListWithTasks {
	listID: ListID;
	tasks: Task[];
}

type TaskStoreState = {
	taskLists: ListWithTasks[];
};

type TaskStoreActions = {
	getTasksForList: (listID: ListID) => Task[];
	addList: (listID: ListID, tasks: Task[]) => void;
	removeList: (listID: ListID) => void;
	addTaskToList: (listID: ListID, task: Task) => void;
	removeTaskFromList: (listID: ListID, taskID: TaskID) => void;
	updateTaskInList: (listID: ListID, updatedTask: Task) => void;
};

const useTaskStore = create<TaskStoreState & TaskStoreActions>((set, get) => ({
	taskLists: [],
	getTasksForList: (listID) =>
		get()
			.taskLists.filter((list) => list.listID == listID)
			.flatMap((list) => list.tasks),
	addList: (listID, tasks) =>
		set((state) => ({
			taskLists: [...state.taskLists, { listID: listID, tasks: tasks }],
		})),
	removeList: (listID) =>
		set((state) => ({
			taskLists: state.taskLists.filter((tl) => tl.listID != listID),
		})),
	addTaskToList: (listID, task) =>
		set((state) => ({
			taskLists: state.taskLists.map((list) =>
				list.listID === listID
					? { ...list, todos: [...list.tasks, task] }
					: list,
			),
		})),
	removeTaskFromList: (listID, taskID) =>
		set((state) => ({
			taskLists: state.taskLists.map((list) =>
				list.listID === listID
					? {
						...list,
						tasks: list.tasks.filter((task) => task.task_id !== taskID),
					}
					: list,
			),
		})),
	updateTaskInList: (listID, updatedTask) =>
		set((state) => ({
			taskLists: state.taskLists.map((list) =>
				list.listID === listID
					? {
						...list,
						tasks: list.tasks.map((task) =>
							task.task_id === updatedTask.task_id ? updatedTask : task,
						),
					}
					: list,
			),
		})),
}));

export default useTaskStore;
