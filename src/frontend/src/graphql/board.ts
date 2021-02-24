import { gql } from "@apollo/client";

const GET_BOARD = gql`
  query getBoard($id: ID!) {
    getBoard(id: $id) {
      name
      lists {
        _id
        name
        cards {
          name
          dueDate
          _id
        }
      }
    }
  }
`
export {
  GET_BOARD
}
