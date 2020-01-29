import * as React from "react";
import { useState } from "react";
import API from "../api";
import { useHistory } from "react-router-dom";

const Login: React.FC = () => {
    const [username, setUsername] = useState("");
    const [password, setPassword] = useState("");
    const [errorMsg, setErrorMsg] = useState("");
    const history = useHistory();

    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        API.post("/login", {
            "username": username,
            "password": password,
        })
            .then(res => {
                history.push("/");
            })
            .catch(_ => {
                setErrorMsg("Failed to Login");
            })
    };

    return (
        <div>
            {errorMsg && <p>{errorMsg}</p>}
            <form action="" onSubmit={onSubmit}>
                <p>
                    username: <input type="text" name="username" onChange={(e) => setUsername(e.target.value)} value={username} />
                </p>
                <p>
                    password: <input type="password" name="password" onChange={(e) => setPassword(e.target.value)} value={password}/>
                </p>
                <p>
                    <input type="submit" value="Login"/>
                </p>
            </form>
        </div>
    );
};

export default Login;
