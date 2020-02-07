import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import {split, InMemoryCache} from 'apollo-boost';
import {ApolloClient} from 'apollo-client';
import {HttpLink} from 'apollo-link-http';
import {WebSocketLink} from 'apollo-link-ws';
import {ApolloProvider} from '@apollo/react-hooks';
import {getMainDefinition} from 'apollo-utilities';
import env from './env';


const wsLink = new WebSocketLink({
  uri: 'ws' + env.SERVER_ADDRESS.slice(4) + '/query',
  options: {
    reconnect: true,
  },
});
const httpLink = new HttpLink({
  uri: env.SERVER_ADDRESS + '/query',
  credentials: 'include',
});
interface Definition {
  kind: string;
  operation?: string;
};
const link = split(
    // split based on operation type
    ({query}) => {
      const {kind, operation}: Definition = getMainDefinition(query);
      return kind === 'OperationDefinition' && operation === 'subscription';
    },
    wsLink,
    httpLink,
);

const client = new ApolloClient({
  link: link,
  cache: new InMemoryCache(),
});


ReactDOM.render(
    <ApolloProvider client={client}>
      <App />
    </ApolloProvider>
    , document.getElementById('root'));

