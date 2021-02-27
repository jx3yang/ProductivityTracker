import React from 'react'
import { BoardComponent } from 'src/components/BoardComponent';
import Board from 'src/models/board';
import { GET_BOARD } from 'src/graphql/board';
import { useQuery } from '@apollo/client';

interface BoardViewProps {}

export const BoardView: React.FC<BoardViewProps> = (props) => {
  const { loading, data } = useQuery<Record<string, Board>>(GET_BOARD, { variables: { id: "602ee9c00f16fe8c02efcb94" }});

  if (loading) return null;

  return (
    <>
      <div> Stub for board view </div>
      {data && data.getBoard && <BoardComponent board={{...data.getBoard}} />}
    </>
  );
}
