import React from "react";
import { ItemCard } from "../ItemCard";

interface Item {
  id: number;
  name: string;
  category_id: number;
  category_name: string;
  user_id: number;
  price: number;
  description: string;
  status: number;
}

interface Prop {
  items: Item[];
}

export const ItemList: React.FC<Prop> = (props) => {
  return (
    <div>
      {props.items &&
        props.items.map((item) => {
          return <ItemCard item={item} key={item.id} />;
        })}
    </div>
  );
};
