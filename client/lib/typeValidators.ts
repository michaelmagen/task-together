// TODO: Rename this file, or split it to other files
import { z } from "zod";

export type ListID = number;
export type TaskID = number;
type UserID = string;

export type List = {
  list_id: ListID;
  name: string;
  creator_id: UserID;
  created_at: string; // Assuming time.Time is represented as a string, update if needed
};

export type User = {
  id: UserID;
  email: string;
  verified_email: boolean;
  name: string;
  given_name: string;
  family_name: string;
  picture: string;
  locale: string;
};

export type Task = {
  task_id: TaskID;
  content: string;
  list_id: ListID;
  creator_id: UserID;
  completer_id?: UserID;
  created_at: Date;
  completed: boolean;
  creator: User;
  completer?: User;
};

// Define the Zod schema for the List type
const listSchema = z.object({
  list_id: z.number(),
  name: z.string(),
  creator_id: z.string(),
  created_at: z.string(),
});

export const userSchema = z.object({
  id: z.string(), // Assuming UserID is a string
  email: z.string(),
  verified_email: z.boolean(),
  name: z.string(),
  given_name: z.string(),
  family_name: z.string(),
  picture: z.string(),
  locale: z.string(),
});

// Define the Zod schema for an array of lists
export const listsSchema = z.array(listSchema);
export const usersSchema = z.array(userSchema);
