import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";

export const Login = () => {
  const [username, setUsername] = useState<string>();
  const [password, setPassword] = useState<string>();
  const [_, setCookie] = useCookies(["userID", "token"]);

  const navigate = useNavigate();

  const onSubmit = (_: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    fetcher<{ id: number; name: string; token: string }>(`/login_name`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        user_name: username,
        password: password,
      }),
    })
      .then((user) => {
        toast.success("Signed in!");
        console.log("POST success:", user.id);
        setCookie("userID", user.id, { "path": "/" });
        setCookie("token", user.token, { "path": "/" });
        navigate("/");
      })
      .catch((err) => {
        console.log(`POST error:`, err);
        toast.error(err.message);
      });
  };

  return (
    <div>
      <div className="Login">
        <label id="MerInputLabel">User Name</label>
        <input
          type="text"
          name="username"
          id="MerTextInput"
          placeholder="name"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setUsername(String(e.target.value));
          }}
          required
        />
        <label id="MerInputLabel">Password</label>
        <input
          type="password"
          name="password"
          id="MerTextInput"
          placeholder="password"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            setPassword(e.target.value);
          }}
        />
        <button onClick={onSubmit} id="MerButton">
          Login
        </button>
      </div>
    </div>
  );
};
