import React from 'react';
import { DragDropContext, Droppable, DropResult, ResponderProvided } from 'react-beautiful-dnd';
import Board from 'src/models/board';
import { ListComponent } from '../ListComponent';

interface BoardComponentProps {
  board: Board;
}

export const BoardComponent: React.FC<BoardComponentProps> = (props) => {
  const { board } = props;
  const { lists, name, boardID } = board;

  const onDragEnd = (result: DropResult, provided: ResponderProvided): void => {
    
  }

  return (
    <DragDropContext onDragEnd={onDragEnd}>
      <div>Title: {name}</div>
      <Droppable droppableId={boardID} direction='horizontal' type='column'>
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
