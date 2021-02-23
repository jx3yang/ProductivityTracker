import React, { useEffect, useState } from 'react';
import { DragDropContext, Droppable, DropResult, ResponderProvided } from 'react-beautiful-dnd';
import Board from 'src/models/board';
import List from 'src/models/list';
import { ListComponent } from 'src/components/ListComponent';
import { COLUMN } from 'src/components/utils/constants';

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

    if (type === COLUMN) {
      const newLists = [...lists];
      const movedElement = newLists.splice(source.index, 1)[0];
      newLists.splice(destination.index, 0, movedElement);
      setLists(newLists);
    }
  }

  return (
    <DragDropContext onDragEnd={onDragEnd}>
      <div>Title: {name}</div>
      <Droppable droppableId={boardID} direction='horizontal' type={COLUMN}>
        {provided => (
          <div ref={provided.innerRef} {...provided.droppableProps} style={{ display: 'flex' }}>
            {lists.map((list, index) => <ListComponent list={list} index={index}/>)}
            {provided.placeholder}
          </div>
        )}
      </Droppable>
    </DragDropContext>
  );
}
