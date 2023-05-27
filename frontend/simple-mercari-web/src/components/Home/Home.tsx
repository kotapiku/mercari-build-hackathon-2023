import { Login } from "../Login";
import { Signup } from "../Signup";
import { ItemList } from "../ItemList";
import { useCookies } from "react-cookie";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import "react-toastify/dist/ReactToastify.css";
import { MerComponent } from "../MerComponent";
import { Container } from "react-bulma-components";

interface Item {
  id: number;
  name: string;
  price: number;
  category_name: string;
  status: number;
}
export const Home = () => {
  const [cookies] = useCookies(["userID", "token"]);
  const [items, setItems] = useState<Item[]>([]);
  const [items_sold, setItemsSold] = useState<Item[]>([]);

  const fetchItems = () => {
    fetcher<Item[]>(`/items`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    })
      .then((data) => {
        console.log("GET success:", data);
        setItems(data);
      })
      .catch((err) => {
        console.log(`GET error:`, err);
        toast.error(err.message);
      });
  };
  const fetchItemsSold = () => {
    fetcher<Item[]>(`/items_sold`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    })
      .then((data) => {
        console.log("GET success:", data);
        setItemsSold(data);
      })
      .catch((err) => {
        console.log(`GET error:`, err);
        toast.error(err.message);
      });
  };

  useEffect(() => {
    fetchItems();
    fetchItemsSold();
  }, []);

  const signUpAndSignInPage = (
    <Container className="p-2">
      <div className="columns is-centered">
        <div className="column is-4">
          <Signup />
        </div>
      </div>
      <div className="columns is-centered">
        <div className="column is-4">
          <Login />
        </div>
      </div>
    </Container>
  );

  const itemListPage = (
    <MerComponent>
      <ItemList items={items} items_sold={items_sold} />
    </MerComponent>
  );
  return <>{cookies.token ? itemListPage : signUpAndSignInPage}</>;
};
