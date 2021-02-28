import { gql } from "@apollo/client";

export interface ChangeListOrder {
  listID: string
  boardID: string
  srcIdx: number
  destIdx: number
}

const UPDATE_LIST_ORDER = gql`
  mutation updateListOrder($changeListOrder: ChangeListOrder!) {
    updateListOrder(changeListOrder: $changeListOrder)
  }
`

export {
  UPDATE_LIST_ORDER
}
