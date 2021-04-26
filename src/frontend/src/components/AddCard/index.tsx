import React, { useState } from 'react';
import Button from '@material-ui/core/Button';
import './style.css';
import { ClickAwayListener, TextField } from '@material-ui/core';

interface AddCardProps {
  onAdd?: (name: string) => void;
}

const AddCard: React.FC<AddCardProps> = (props) => {
  const { onAdd } = props;
  const [adding, setAdding] = useState<boolean>(false);
  const [newCardName, setNewCardName] = useState<string>('');

  const onClickAway = () => {
    setAdding(false);
  }

  const onClickAdd = () => {
    setNewCardName('');
    setAdding(true);
  }

  const onSave = () => {
    if (onAdd) onAdd(newCardName);
    setAdding(false);
  }

  return (
    adding
    ? <ClickAwayListener onClickAway={onClickAway}>
        <div>
          <TextField autoFocus onChange={(e) => setNewCardName(e.target.value)} value={newCardName} className="addCardName" />
          <Button color="primary" onClick={onSave} className="saveButton">Save</Button>
        </div>
      </ClickAwayListener>
    : <Button color="primary" onClick={onClickAdd} className="addCard">
        Add Card
      </Button>
  );
}

export default AddCard;
