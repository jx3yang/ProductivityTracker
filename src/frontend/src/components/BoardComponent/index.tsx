import React, { useEffect, useState } from 'react';
import { DragDropContext, Droppable, DropResult, ResponderProvided } from 'react-beautiful-dnd';
import Board from 'src/models/board';
import List from 'src/models/list';
import { ListComponent } from 'src/components/ListComponent';
import { CARD, LIST } from 'src/components/utils/constants';
import Card from 'src/models/card';
import { moveElement, addOneToList, removeOneFromList } from 'src/components/utils/utils';

import { ChangeListOrder, CREATE_LIST, NewList, UPDATE_LIST_ORDER } from 'src/graphql/list';
import { ApolloQueryResult, useMutation } from '@apollo/client';
import { ChangeCardOrder, NewCard, UPDATE_CARD_ORDER, CREATE_CARD } from 'src/graphql/card';
import { AddList } from '../AddList';
import { Grid } from '@material-ui/core';

interface BoardComponentProps {
  board: Board;
  refetch: (variables?: Partial<Record<string, any>> | undefined) => Promise<ApolloQueryResult<Record<string, Board>>>
}

export const BoardComponent: React.FC<BoardComponentProps> = (props) => {
  const { board, refetch } = props;
  const [lists, setLists] = useState<List[]>([]);
  const [boardId, setBoardId] = useState<string>('');
  const [name, setBoardName] = useState<string>('');

  const [updateListOrder] = useMutation(UPDATE_LIST_ORDER);
  const [updateCardOrder] = useMutation(UPDATE_CARD_ORDER);

  const [createList] = useMutation(CREATE_LIST);
  const [createCard] = useMutation(CREATE_CARD);

  useEffect(() => {
    setLists(board.lists || []);
    setBoardId(board._id);
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
      setLists(moveElement(currLists, source.index, destination.index));

      const changeListOrder: ChangeListOrder = {
        listId: draggableId,
        boardId,
        srcIdx: source.index,
        destIdx: destination.index,
      };
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

      const changeCardOrder: ChangeCardOrder = {
        boardId,
        srcListId: source.droppableId,
        destListId: destination.droppableId,
        srcIdx: source.index,
        destIdx: destination.index,
        cardId: draggableId,
      };

      // if update fails, revert to original order
      // TODO: have a graphql subscription to keep board updated
      // TODO: do not revert to original order if user is offline,
      //       instead tell the user they are offline and that changes will not be saved
      updateCardOrder({ variables: { changeCardOrder } })
        .then(res => {
          if (!res.data?.updateCardOrder) {
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
  }

  const addList = (name: string) => {
    if (name !== '') {
      const list: NewList = {
        name,
        parentBoardId: boardId,
      };

      createList({ variables: { list }})
        .then(res => {
          if (res.data?.createList?._id) {
            const { _id } = res.data?.createList;
            setLists([...lists, { name, _id, cards: [] }]);
          }
          else refetch();
        })
        .catch(_ => refetch());
    }
  }

  const addCard = (name: string, parentListId: string) => {
    if (name != '') {
      const card: NewCard = {
        name,
        parentListId,
        parentBoardId: boardId,
      };

      createCard({ variables: { card }})
        .then(res => {
          if (res.data?.createCard?._id) {
            const { _id } = res.data?.createCard;
            const newCard: Card = { _id, name };
            const currLists = [...lists];
            const sourceListIdx = currLists.findIndex(list => list._id === parentListId);
            const sourceList = currLists[sourceListIdx];
            const newOrder = addOneToList(sourceList.cards || [], (sourceList.cards || []).length, newCard);
            currLists[sourceListIdx] = { ...sourceList, cards: newOrder };
            setLists(currLists);
          }
          else refetch();
        })
        .catch(_ => refetch());
    }
  }

  return (
    <DragDropContext onDragEnd={onDragEnd}>
      <div>Title: {name}</div>
      <div style={{ overflowX: 'scroll', height: '95%' }}>
        <Grid container style={{ flexWrap: 'nowrap' }}>
            {boardId && <Droppable droppableId={boardId} direction='horizontal' type={LIST}>
              {provided => (
                  <div ref={provided.innerRef} {...provided.droppableProps} style={{ display: 'flex' }}>
                    {lists.map((list, index) => 
                      <Grid item xs key={list._id}>
                        <ListComponent list={list} index={index} key={list._id} onAddCard={addCard} />
                      </Grid>
                    )}
                    {provided.placeholder}
                  </div>
              )}
            </Droppable>}
          <Grid item xs>
            <AddList onAdd={addList} />
          </Grid>
        </Grid>
      </div>
    </DragDropContext>
  );
}
