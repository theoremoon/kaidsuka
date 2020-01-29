import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import { model } from "./model";
import { initState } from "./store";


const {store, Provider} = model.createStore({
    initState,
});

ReactDOM.render(
    <Provider>
        <App />
    </Provider>
    ,
    document.getElementById("root")
);
