import { createContext, useContext, useState } from "react";

export interface Item {
  id: number;
  name: string;
  category_id: number;
  category_name: string;
  user_id: number;
  price: number;
  description: string;
  status: number;
}

interface ItemsContextType {
  items: Item[];
  setItems: React.Dispatch<React.SetStateAction<Item[]>>;
}

const ItemsContext = createContext<ItemsContextType | undefined>(undefined);

export const ItemsProvider: React.FC = ({ children }) => {
  const [items, setItems] = useState<Item[]>([]);

  return (
    <ItemsContext.Provider value={{ items, setItems }}>
      {children}
    </ItemsContext.Provider>
  );
};

export const useItems = () => {
  const context = useContext(ItemsContext);

  if (context === undefined) {
    throw new Error("useItems must be used within an ItemsProvider");
  }

  return context;
};
