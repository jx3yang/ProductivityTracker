import React from 'react';
import { Draggable, Droppable } from 'react-beautiful-dnd';
import List from 'src/models/list';
import { CardComponent } from 'src/components/CardComponent';
import { CARD } from 'src/components/utils/constants';
import './style.css';

interface ListComponentProps {
  list: List;
  index: number;
}

export const ListComponent: React.FC<ListComponentProps> = (props) => {
  const { list, index } = props;
  const { listID, cards, name } = list;

  return (
    <Draggable draggableId={listID} index={index} key={listID}>
      {provided => (
        <div ref={provided.innerRef} {...provided.draggableProps} className='list'>
          <div {...provided.dragHandleProps}>
            title: {name}
          </div>
          <Droppable droppableId={listID} type={CARD} key={listID}>
            {provided => (
                <div ref={provided.innerRef} {...provided.droppableProps} className='wrapper'>
                  {cards.map((card, index) => (<CardComponent card={card} index={index} key={card.cardID} />))}
                  {provided.placeholder}
                </div>
              )
            }
          </Droppable>
        </div> 
      )}
    </Draggable>
  );
}
