import React, { ReactNode } from "react";
import { useCookies } from "react-cookie";
import { NotFound } from "../NotFound";
import { Container } from "react-bulma-components";

interface Prop {
  condition?: () => boolean;
  children: ReactNode;
}

export const MerComponent: React.FC<Prop> = (props) => {
  const [cookies] = useCookies(["token", "userID"]);

  if (
    !cookies.token ||
    !cookies.userID ||
    (props.condition && props.condition() === false)
  ) {
    return <NotFound />;
  }
  return <Container className="p-2">{props.children}</Container>;
};
