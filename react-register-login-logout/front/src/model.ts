import { createModel } from "@prodo/core";

export interface User {
    is_login: boolean;
    username: string;
};

export const model = createModel<User>();
export const {state, watch, dispatch} = model.ctx;