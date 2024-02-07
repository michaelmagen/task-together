"use client";

import useSWRMutation from "swr/mutation";
import useSWR, { useSWRConfig } from "swr";
import useListStore from "@/lib/stores/listStore";
import useCookieStore from "@/lib/stores/cookieStore";
import { endpoints, FetcherOptions, Method } from "@/lib/fetcher";
import type { Fetcher } from "swr";
import fetcher from "@/lib/fetcher";
import { ListID, Task } from "@/lib/typeValidators";
import Loading from "./loading";
import { Button } from "./ui/button";
import { Loader2, Trash2 } from "lucide-react";
import UserCard from "./user-card";
import { PropsWithChildren, useState } from "react";
import useUserStore from "@/lib/stores/userStore";
import useTaskStore from "@/lib/stores/taskStore";

export default function TasksDisplay() {
  const { selectedList } = useListStore();

  // TODO: Make this look better, esp in mobile view
  if (!selectedList) {
    return <div>Select a list to view it!</div>;
  }

  return (
    <div className="flex flex-grow flex-col items-center p-4">
      <h1 className="text-center text-4xl font-extrabold tracking-tight lg:text-5xl">
        {selectedList.name}
      </h1>
      <TaskTable listID={selectedList.list_id} />
    </div>
  );
}

interface TaskTableProps {
  listID: ListID;
}

function TaskTable({ listID }: TaskTableProps) {
  const { cookies } = useCookieStore();
  const { getTasksForList, addList, removeList } = useTaskStore();

  const fetcherOptions: FetcherOptions = {
    method: Method.GET,
    cookieString: cookies,
  };
  const taskFetcher: Fetcher<Task[], string> = (endpoint) =>
    fetcher(endpoint, fetcherOptions);
  const { data, error, isLoading } = useSWR(
    endpoints.tasks(listID),
    taskFetcher,
    {
      onSuccess: (data) => {
        removeList(listID);
        addList(listID, data);
      },
    },
  );

  if (isLoading) {
    return <Loading />;
  }

  // No tasks yet in current list
  if (data === null) {
    return <div>No tasks in list</div>;
  }

  // Failed to fetch tasks
  if (error || !data) {
    return <div>error loading data</div>;
  }

  // TODO: Make this a scroll view
  return (
    <div className="flex w-full flex-col items-center gap-4 p-4">
      {getTasksForList(listID).map((task) => (
        <Task task={task} key={task.task_id} />
      ))}
    </div>
  );
}

interface TaskProps {
  task: Task;
}

function Task({ task }: TaskProps) {
  const { cookies } = useCookieStore();
  const { user } = useUserStore();
  const { updateTaskInList } = useTaskStore();

  const taskCompleteMutation = (
    key: string,
    { arg }: { arg: FetcherOptions },
  ) => fetcher(key, arg);
  const { trigger, isMutating: mutatingTask } = useSWRMutation<
    Task,
    Error,
    string,
    FetcherOptions
  >(endpoints.singleTask(task.list_id, task.task_id), taskCompleteMutation);

  const toggleTaskCompletion = () => {
    const fetcherOptions: FetcherOptions = {
      method: Method.PATCH,
      cookieString: cookies,
      body: {
        completed: !task.completed,
      },
    };

    let newTask: Task = { ...task };
    if (!task.completed && user !== null) {
      newTask.completed = true;
      newTask.completer = user;
      newTask.completer_id = user.id;
    } else {
      newTask.completed = false;
      const { completer, completer_id, ...taskWithoutCompleter } = newTask;
      newTask = taskWithoutCompleter;
    }

    trigger(fetcherOptions, {
      onSuccess: () => updateTaskInList(newTask.list_id, newTask),
    });
  };

  return (
    <div
      className={`flex w-full flex-col gap-1 rounded-lg border p-3 ${task.completed ? "bg-green-500 bg-opacity-25" : ""
        }`}
    >
      <div className="flex w-full items-center gap-4">
        <span className="flex-grow text-lg">{task.content}</span>
        <Button
          className=""
          variant="outline"
          onClick={toggleTaskCompletion}
          disabled={mutatingTask}
        >
          {mutatingTask && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
          Complete
        </Button>
        <Button className="" variant="destructive" size="icon">
          <Trash2 />
        </Button>
      </div>
      <TaskDescription>
        <span>Created by </span>
        <UserCard user={task.creator} />
        <span>.</span>
        {task.completed && task.completer && (
          <>
            <span> Completed by </span>
            <UserCard user={task.completer} />
            <span>.</span>
          </>
        )}
      </TaskDescription>
    </div>
  );
}

function TaskDescription(props: PropsWithChildren) {
  return (
    <span className="text-muted-foreground text-sm">{props.children} </span>
  );
}
