import React from 'react';
import { Draggable } from 'react-beautiful-dnd';
import Card from 'src/models/card';
import './style.css';

interface CardComponentProps {
  card: Card;
  index: number;
}

export const CardComponent: React.FC<CardComponentProps> = (props) => {
  const { card, index } = props;
  const { _id: cardID, name } = card;

  return (
    <Draggable draggableId={cardID} index={index} key={cardID}>
      {provided => (
        <div
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...provided.dragHandleProps}
          className="card"
        >
          card: {name}
        </div>
      )}
    </Draggable>
  );
}
