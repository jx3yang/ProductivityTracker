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
  const { _id: cardId, name } = card;

  return (
    <Draggable draggableId={cardId} index={index} key={cardId}>
      {provided => (
        <div
          ref={provided.innerRef}
          {...provided.draggableProps}
          {...provided.dragHandleProps}
          className="card"
        >
          card: {cardId}
        </div>
      )}
    </Draggable>
  );
}
