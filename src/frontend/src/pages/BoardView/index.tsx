import React from 'react'
import { BoardComponent } from 'src/components/BoardComponent';
import Board from 'src/models/board';

interface BoardViewProps {
  
}

export const BoardView: React.FC<BoardViewProps> = (props) => {

  const testBoard: Board = {
    boardID: 'id',
    name: 'board',
    lists: [
      {
        cards: [
          {
            cardID: '1',
            name: '1',
          }
        ],
        listID: '1',
        name: 'list1',
      },
      {
        cards: [
          {
            cardID: '21',
            name: '21',
          }
        ],
        listID: '2',
        name: 'list2',
      },
    ]
  }

  return (
    <>
      <div> Stub for board view </div>
      <BoardComponent board={testBoard} />
    </>
  );
}
