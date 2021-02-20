import List from "./list";

export default interface Board {
  boardID: string;
  name: string;
  lists: List[];
}
