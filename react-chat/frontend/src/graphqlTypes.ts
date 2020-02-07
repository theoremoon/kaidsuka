/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: getMessages
// ====================================================

export interface getMessages_getMessages_user {
  readonly __typename: "User";
  readonly username: string;
}

export interface getMessages_getMessages {
  readonly __typename: "Message";
  readonly id: string;
  readonly text: string;
  readonly user: getMessages_getMessages_user;
  readonly edited: boolean;
  /**
   * posted_at is unix time
   */
  readonly posted_at: number;
}

export interface getMessages {
  readonly getMessages: ReadonlyArray<getMessages_getMessages>;
}

/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL mutation operation: postMessage
// ====================================================

export interface postMessage_postMessage {
  readonly __typename: "Message";
  readonly id: string;
}

export interface postMessage {
  readonly postMessage: postMessage_postMessage;
}

export interface postMessageVariables {
  readonly text: string;
}

/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL query operation: getLoginUser
// ====================================================

export interface getLoginUser_user {
  readonly __typename: "User";
  readonly username: string;
}

export interface getLoginUser {
  readonly user: getLoginUser_user;
}

/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

// ====================================================
// GraphQL subscription operation: newPost
// ====================================================

export interface newPost_newMessage_user {
  readonly __typename: "User";
  readonly username: string;
}

export interface newPost_newMessage {
  readonly __typename: "Message";
  readonly id: string;
  readonly text: string;
  readonly user: newPost_newMessage_user;
  readonly edited: boolean;
  /**
   * posted_at is unix time
   */
  readonly posted_at: number;
}

export interface newPost {
  readonly newMessage: newPost_newMessage;
}

/* tslint:disable */
/* eslint-disable */
// This file was automatically generated and should not be edited.

//==============================================================
// START Enums and Input Objects
//==============================================================

//==============================================================
// END Enums and Input Objects
//==============================================================
