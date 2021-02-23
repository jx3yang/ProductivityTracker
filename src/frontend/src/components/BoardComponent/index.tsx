import React, { useEffect, useState } from 'react';
import { DragDropContext, Droppable, DropResult, ResponderProvided } from 'react-beautiful-dnd';
import Board from 'src/models/board';
import List from 'src/models/list';
import { ListComponent } from 'src/components/ListComponent';
import { CARD, LIST } from 'src/components/utils/constants';

function moveElement<T>(oldList: Array<T>, sourceIdx: number, destIdx: number): Array<T> {
  const newList = [...oldList];
  const movedElement = newList.splice(sourceIdx, 1)[0];
  newList.splice(destIdx, 0, movedElement);
  return newList;
}

interface BoardComponentProps {
  board: Board;
}

export const BoardComponent: React.FC<BoardComponentProps> = (props) => {
  const { board } = props;
  const [lists, setLists] = useState<List[]>([]);
  const [boardID, setBoardID] = useState<string>('');
  const [name, setBoardName] = useState<string>('');

  useEffect(() => {
    setLists(board.lists);
    setBoardID(board.boardID);
    setBoardName(board.name);
  }, [board]);

  const onDragEnd = (result: DropResult, provided: ResponderProvided): void => {
    const { destination, source, draggableId, type } = result;
    if (!destination) return;
    if (destination.droppableId === source.droppableId && destination.index === source.index) return;
    
    // moving lists around
    if (type === LIST) {
      setLists(moveElement(lists, source.index, destination.index));
      return;
    }

    // moving cards around
    if (type === CARD) {
      const sourceList = lists.find(list => list.listID === source.droppableId)!;

      // changing order within same list
      if (source.droppableId === destination.droppableId) {
        sourceList.cards = moveElement(sourceList.cards, source.index, destination.index);
        setLists([...lists]);
      }

      // moving card from one list to another
      else {
        const destinationList = lists.find(list => list.listID === destination.droppableId)!;
        const card = sourceList.cards.splice(source.index, 1)[0];
        destinationList.cards.splice(destination.index, 0, card);
        setLists([...lists]);
      }
    }
  }

  return (
    <DragDropContext onDragEnd={onDragEnd}>
      <div>Title: {name}</div>
      <Droppable droppableId={boardID} direction='horizontal' type={LIST}>
        {provided => (
          <div ref={provided.innerRef} {...provided.droppableProps} style={{ display: 'flex' }}>
            {lists.map((list, index) => <ListComponent list={list} index={index} key={list.listID} />)}
            {provided.placeholder}
          </div>
        )}
      </Droppable>
    </DragDropContext>
  );
}
