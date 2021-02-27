import Card from "./card";

export default interface List {
  _id: string;
  name: string
  cards?: Card[];
}
