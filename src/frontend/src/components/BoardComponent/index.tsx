import React, { useEffect, useState } from 'react';
import { DragDropContext, Droppable, DropResult, ResponderProvided } from 'react-beautiful-dnd';
import Board from 'src/models/board';
import List from 'src/models/list';
import { ListComponent } from 'src/components/ListComponent';
import { CARD, LIST } from 'src/components/utils/constants';
import Card from 'src/models/card';
import { moveElement, addOneToList, removeOneFromList } from 'src/components/utils/utils';

import { ChangeListOrder, UPDATE_LIST_ORDER } from 'src/graphql/list';
import { ApolloQueryResult, useMutation } from '@apollo/client';

interface BoardComponentProps {
  board: Board;
  refetch: (variables?: Partial<Record<string, any>> | undefined) => Promise<ApolloQueryResult<Record<string, Board>>>
}

export const BoardComponent: React.FC<BoardComponentProps> = (props) => {
  const { board, refetch } = props;
  const [lists, setLists] = useState<List[]>([]);
  const [boardID, setBoardID] = useState<string>('');
  const [name, setBoardName] = useState<string>('');

  const [updateListOrder] = useMutation(UPDATE_LIST_ORDER);

  useEffect(() => {
    setLists(board.lists || []);
    setBoardID(board._id);
    setBoardName(board.name);
  }, [board]);

  const onDragEnd = (result: DropResult, provided: ResponderProvided): void => {
    const { destination, source, draggableId, type } = result;
    if (!destination) return;
    if (destination.droppableId === source.droppableId && destination.index === source.index) return;
    const oldLists = [...lists];
    const currLists = [...lists];
    
    // moving lists around
    if (type === LIST) {
      const changeListOrder: ChangeListOrder = {
        listID: draggableId,
        boardID,
        srcIdx: source.index,
        destIdx: destination.index,
      }
      setLists(moveElement(currLists, source.index, destination.index));

      // if update fails, revert to original order
      // TODO: have a graphql subscription to keep board updated
      // TODO: do not revert to original order if user is offline,
      //       instead tell the user they are offline and that changes will not be saved
      updateListOrder({ variables: { changeListOrder } })
        .then(res => {
          if (!res.data?.updateListOrder) {
            setLists(oldLists);
            refetch();
          }
        })
        .catch(_ => {
          setLists(oldLists);
          refetch();
        });

      return;
    }

    // moving cards around
    if (type === CARD) {
      const sourceListIdx = currLists.findIndex(list => list._id === source.droppableId);
      const sourceList = currLists[sourceListIdx];

      // changing order within same list
      if (source.droppableId === destination.droppableId) {
        const newOrder = moveElement(sourceList.cards!, source.index, destination.index);
        currLists[sourceListIdx] = {...sourceList, cards: newOrder}
        setLists(currLists);
      }

      // moving card from one list to another
      else {
        const destinationListIdx = currLists.findIndex(list => list._id === destination.droppableId);
        const destinationList = currLists[destinationListIdx];

        const card: Card = {...sourceList.cards![source.index]};
        const newSourceOrder = removeOneFromList(sourceList.cards!, source.index);
        const newDestinationOrder = addOneToList(destinationList.cards || [], destination.index, card);

        currLists[sourceListIdx] = {...sourceList, cards: newSourceOrder};
        currLists[destinationListIdx] = {...destinationList, cards: newDestinationOrder};

        setLists(currLists);
      }
    }
  }

  return (
    <DragDropContext onDragEnd={onDragEnd}>
      <div>Title: {name}</div>
      {boardID && <Droppable droppableId={boardID} direction='horizontal' type={LIST}>
        {provided => (
          <div ref={provided.innerRef} {...provided.droppableProps} style={{ display: 'flex' }}>
            {lists.map((list, index) => <ListComponent list={list} index={index} key={list._id} />)}
            {provided.placeholder}
          </div>
        )}
      </Droppable>}
    </DragDropContext>
  );
}
