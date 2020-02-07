import React, {ReactElement, useState, ChangeEvent} from 'react';
import gql from 'graphql-tag';
import {useQuery, useMutation} from '@apollo/react-hooks';
import {getMessages, getLoginUser, postMessage, newPost} from '../graphqlTypes';
import './Chat.css';

const GET_MESSAGES = gql`
query getMessages {
    getMessages {
        id
        text
            user {
                username
            }
        edited
        posted_at
    }
}
`;

const POST_MESSAGE = gql`
mutation postMessage($text: String!) {
  postMessage(text: $text) {
    id
  }
}
`;

const GET_LOGIN_USER = gql`
query getLoginUser {
  user {
    username
  }
}
`;

const SUBSCRIBE_NEW_POST = gql`
subscription newPost {
  newMessage {
    id
    text
    user {
      username
    }
    edited
    posted_at
  }
}
`;

const Chat = (): ReactElement => {
  const {loading, error, data, subscribeToMore} = useQuery<getMessages>(GET_MESSAGES);
  subscribeToMore<newPost>({
    document: SUBSCRIBE_NEW_POST,
    updateQuery: (prev, {subscriptionData}): getMessages => {
      if (!subscriptionData.data) {
        return prev;
      }
      const feedItem: newPost = subscriptionData.data;
      if (prev.getMessages.filter((x) => x.id == feedItem.newMessage.id).length > 0) {
        return prev;
      }
      return {
        getMessages: [feedItem.newMessage, ...prev.getMessages],
      };
    },
  });
  const {loading: userLoading, error: userError, data: userData} = useQuery<getLoginUser>(GET_LOGIN_USER);
  const [postMessage] = useMutation<postMessage>(POST_MESSAGE);

  const formatPostedDate = (postedAt: number): string => {
    return new Date(postedAt * 1000).toISOString();
  };

  const [state, setState] = useState({
    comment: '',
  });
  const onSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    postMessage({
      variables: {
        text: state.comment,
      },
    });
    setState({...state, comment: ''});
  };

  if (loading || userLoading) {
    return <>Loading...</>;
  }
  if (error) {
    return <>{error.message}</>;
  }

  return (
    <>
      <div className="messages">
        {data && data.getMessages.slice(0).reverse().map((msg) => (
          <div className="message" key={msg.id}>
            <div className="user">{msg.user.username}</div>
            <div className="text">{msg.text}
              <span className="postedat">{formatPostedDate(msg.posted_at)}</span>
            </div>
          </div>
        ))}
      </div>
      {userError ? <>
      Please Login
      </> : <form onSubmit={onSubmit}>
        {userData&& userData.user.username}:
        <input
          type="text"
          value={state.comment}
          onChange={(e: ChangeEvent<HTMLInputElement>) => {
            setState({...state, comment: e.target.value});
          }}/><button>Send</button>
      </form>}
    </>
  );
};

export default Chat;
