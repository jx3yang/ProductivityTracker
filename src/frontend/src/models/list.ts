import Card from "./card";

export default interface List {
  listID: string;
  name: string
  cards: Card[];
}
