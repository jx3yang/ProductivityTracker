import List from "./list";

export default interface Board {
  _id: string;
  name: string;
  lists?: List[];
}
