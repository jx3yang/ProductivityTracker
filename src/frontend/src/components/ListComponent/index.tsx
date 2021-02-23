import React from 'react';
import { Draggable, Droppable } from 'react-beautiful-dnd';
import List from 'src/models/list';

interface ListComponentProps {
  list: List;
  index: number;
}

export const ListComponent: React.FC<ListComponentProps> = (props) => {
  const { list, index } = props;
  const { listID, cards, name } = list;

  const listStyle: React.CSSProperties = {
    margin: '8px',
    borderRadius: '2px',
    border: '1px solid lightgrey',
    backgroundColor: 'white',
    width: '220px',
    display: 'flex',
    flexDirection: 'column',
  }

  return (
    <Draggable draggableId={listID} index={index} key={listID}>
      {provided => (
        <div ref={provided.innerRef} {...provided.draggableProps}>
          <div style={listStyle}>
            <div {...provided.dragHandleProps}>
              title: {name}
            </div>
            <Droppable droppableId={listID}>
              {provided => (
                  <div ref={provided.innerRef} {...provided.droppableProps} style={{ padding: '8px' }}>
                    {cards.map(card => (<div style={{ padding: '8px' }}> {card.name} </div>))}
                  </div>
                )
              }
            </Droppable>
          </div>
        </div> 
      )}
    </Draggable>
  );
}
