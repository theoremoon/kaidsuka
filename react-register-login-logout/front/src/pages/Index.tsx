import * as React from "react";
import { useEffect } from "react";
import { useState } from "react";
import API from "../api";
import { useHistory } from "react-router-dom";

const Index: React.FC = () => {
    const [username, setUsername] = useState("");
    const history = useHistory();

    useEffect(() => {
        API.get("/user")
        .then(res => {
            setUsername(res.data.username)
        })
        .catch(_ => {
            history.push("/login");
        })
    })
    const logout = () => {
        API.post("/logout")
        history.push("/login");
    }

      return (
        <div>
            Hello {username}
            <br/>
            <a href="" onClick={logout}>Logout</a>
        </div>
    );
};

export default Index;
