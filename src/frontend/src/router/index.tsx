import React from 'react'
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import { BoardView } from 'src/pages/BoardView';
import { Welcome } from 'src/pages/Welcome';

interface RouterProps {
  
}

export const Router: React.FC<RouterProps> = (props) => {
  return (
    <BrowserRouter>
      <Switch>
        <Route path="/" exact component={Welcome} />
        <Route path="/board/:board" component={BoardView} />
      </Switch>
    </BrowserRouter>
  );
}
