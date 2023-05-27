import React from "react";
import { Item } from "../Item";
import { Columns } from "react-bulma-components";
interface Item {
  id: number;
  name: string;
  price: number;
  category_name: string;
}

interface Prop {
  items: Item[];
  items_sold: Item[];
}

export const ItemList: React.FC<Prop> = (props) => {
  return (
    <Columns className="m-1">
      {props.items &&
        props.items.map((item) => {
          return <Item item={item} key={item.id} />;
        })}
      {props.items_sold &&
        props.items_sold.map((item) => {
          return <Item item={item} key={item.id} />;
        })}
    </Columns>
  );
};
