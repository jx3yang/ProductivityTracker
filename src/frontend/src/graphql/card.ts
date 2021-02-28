import { gql } from "@apollo/client";

export interface ChangeCardOrder {
  boardId: string;
  srcListId: string;
  destListId: string;
  cardId: string;
  srcIdx: number;
  destIdx: number;
}

const UPDATE_CARD_ORDER = gql`
  mutation updateCardOrder($changeCardOrder: ChangeCardOrder!) {
    updateCardOrder(changeCardOrder: $changeCardOrder)
  }
`

export {
  UPDATE_CARD_ORDER
}
