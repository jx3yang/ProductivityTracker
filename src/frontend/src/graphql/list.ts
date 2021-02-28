import { gql } from "@apollo/client";

export interface ChangeListOrder {
  listId: string
  boardId: string
  srcIdx: number
  destIdx: number
}

export interface NewList {
  name: string
  parentBoardId: string
}

const UPDATE_LIST_ORDER = gql`
  mutation updateListOrder($changeListOrder: ChangeListOrder!) {
    updateListOrder(changeListOrder: $changeListOrder)
  }
`

const CREATE_LIST = gql`
  mutation createList($list: NewList!) {
    createList(list: $list) {
      _id
    }
  }
`

export {
  UPDATE_LIST_ORDER,
  CREATE_LIST,
}
