import { useLocation } from "react-router-dom";
import { useEffect, useState } from "react";
import { ItemList } from "../ItemList";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import { Item, useItems } from "../../common/context";

export const Search: React.FC = () => {
  const location = useLocation();
  const { items, setItems } = useItems();
  const [loading, setLoading] = useState(false);

  const fetchItems = () => {
    const keyword = new URLSearchParams(location.search).get("keyword");
    if (keyword) {
      fetcher<Item[]>(`/search?keyword=${keyword}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
      })
        .then((data) => {
          console.log("GET success:", data);
          setItems(data);
          setLoading(false);
        })
        .catch((err) => {
          console.log(`GET error:`, err);
          toast.error(err.message);
        });
    }
  };

  useEffect(() => {
    setLoading(true);
    fetchItems();
  }, [location.search]);

  if (loading) {
    return <div>Loading...</div>;
  } else if (!items) {
    return <div>No search result.</div>;
  } else {
    return <ItemList items={items} />;
  }
};
