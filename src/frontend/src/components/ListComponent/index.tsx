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
  const { _id: listId, cards, name } = list;

  return (
    <>
      {listId && <Draggable draggableId={listId} index={index} key={listId}>
        {provided => (
          <div ref={provided.innerRef} {...provided.draggableProps} className='list'>
            <div {...provided.dragHandleProps}>
              title: {listId}
            </div>
            <Droppable droppableId={listId} type={CARD} key={listId}>
              {provided => (
                  <div ref={provided.innerRef} {...provided.droppableProps} className='wrapper'>
                    {cards?.map((card, index) => (<CardComponent card={card} index={index} key={card._id} />))}
                    {provided.placeholder}
                  </div>
                )
              }
            </Droppable>
          </div>
        )}
      </Draggable>}
    </>
  );
}
