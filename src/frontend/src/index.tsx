import { ApolloProvider } from '@apollo/client';
import React from 'react';
import ReactDOM from 'react-dom';
import 'src/index.css';
import reportWebVitals from 'src/reportWebVitals';
import { client } from './graphql/client';
import { Router } from './router';

ReactDOM.render(
  <React.StrictMode>
    <ApolloProvider client={client}>
      <Router />
    </ApolloProvider>
  </React.StrictMode>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
