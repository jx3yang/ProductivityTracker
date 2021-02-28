import React, { useState } from 'react';
import Button from '@material-ui/core/Button';
import { ClickAwayListener, TextField } from '@material-ui/core';
import './style.css';

interface AddListProps {
  onAdd?: (name: string) => void;
}

export const AddList: React.FC<AddListProps> = (props) => {
  const { onAdd } = props;
  const [adding, setAdding] = useState<boolean>(false);
  const [newListName, setNewListName] = useState<string>('');

  const onClickAway = () => {
    setAdding(false);
  }

  const onClickAdd = () => {
    setNewListName('');
    setAdding(true);
  }

  const onSave =() => {
    if (onAdd) {
      onAdd(newListName);
    }
    setAdding(false);
  }

  return (
    adding
    ? <ClickAwayListener onClickAway={onClickAway}>
        <div>
          <TextField autoFocus onChange={(e) => setNewListName(e.target.value)} value={newListName} className="addList"/>
          <Button color="primary" onClick={onSave} className="saveButton">Save</Button>
        </div>
      </ClickAwayListener>
    : <Button color="primary" onClick={onClickAdd} className="addList">
        Add list
      </Button>
  )
}
