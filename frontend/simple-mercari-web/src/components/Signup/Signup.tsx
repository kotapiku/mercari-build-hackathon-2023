import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";

export const Signup = () => {
  const [name, setName] = useState<string>();
  const [password, setPassword] = useState<string>("");
  const [userID, setUserID] = useState<number>();
  const [_, setCookie] = useCookies(["userID"]);

  const [validNameLength, setvalidNameLength] = useState(false);
  const [validPassLength, setvalidPassLength] = useState(false);
  const [hasNumber, setHasNumber] = useState(false);

  const navigate = useNavigate();
  const onSubmit = (_: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    fetcher<{ id: number; name: string }>(`/register`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name: name,
        password: password,
      }),
    })
      .then((user) => {
        toast.success("New account is created!");
        console.log("POST success:", user.id);
        setCookie("userID", user.id, { "path": "/" });
        setUserID(user.id);
        navigate("/");
      })
      .catch((err) => {
        console.log(`POST error:`, err);
        toast.error("This User Name is already used");
      });
  };

  const validateNameLength = (name: string): boolean => {
    return name.length >= 3;
  };
  const validatePassLength = (password: string): boolean => {
    return password.length >= 6;
  };
  const validateNumber = (password: string): boolean => {
    return /\d/.test(password);
  };

  return (
    <div>
      <div className="Signup">
        <label id="MerInputLabel">User Name</label>
        <input
          className="is-fullwidth"
          type="text"
          name="name"
          id="MerTextInput"
          placeholder="name"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            const enteredName = e.target.value;
            setName(enteredName);

            const isvalidNameLength = validateNameLength(enteredName);
            setvalidNameLength(isvalidNameLength);
          }}
          required
        />
        <label id="Caution">
          {validNameLength ? (
            <span></span>
          ) : (
            <span>　＊More than 3 Character</span>
          )}
        </label>
        <br></br>
        <label id="MerInputLabel">Password</label>
        <input
          type="password"
          name="password"
          id="MerTextInput"
          placeholder="password"
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            const enteredPassword = e.target.value;
            setPassword(enteredPassword);

            const isvalidPassLength = validatePassLength(enteredPassword);
            setvalidPassLength(isvalidPassLength);

            const isHasNumber = validateNumber(enteredPassword);
            setHasNumber(isHasNumber);
          }}
        />
        <label id="Caution">
          {validPassLength ? (
            <span></span>
          ) : (
            <span>
              　＊More than 6 Character<br></br>{" "}
            </span>
          )}
          {hasNumber ? <span></span> : <span>　＊Need a Number</span>}
        </label>

        <button
          onClick={onSubmit}
          disabled={!validNameLength || !validPassLength || !hasNumber}
          id="MerButton"
        >
          Signup
        </button>
      </div>
    </div>
  );
};
