import { gql } from "@apollo/client";

export interface ChangeCardOrder {
  boardId: string;
  srcListId: string;
  destListId: string;
  cardId: string;
  srcIdx: number;
  destIdx: number;
}

export interface NewCard {
  name: string;
  parentListId: string;
  parentBoardId: string;
}

const UPDATE_CARD_ORDER = gql`
  mutation updateCardOrder($changeCardOrder: ChangeCardOrder!) {
    updateCardOrder(changeCardOrder: $changeCardOrder)
  }
`

const CREATE_CARD = gql`
  mutation createCard($card: NewCard!) {
    createCard(card: $card) {
      _id
    }
  }
`

export {
  UPDATE_CARD_ORDER,
  CREATE_CARD,
}
