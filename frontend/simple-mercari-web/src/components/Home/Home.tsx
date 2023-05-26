import { Login } from "../Login";
import { Signup } from "../Signup";
import { ItemList } from "../ItemList";
import { useCookies } from "react-cookie";
import { MerComponent } from "../MerComponent";
import { useEffect } from "react";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import "react-toastify/dist/ReactToastify.css";
import { Item, useItems } from "../../common/context";

export const Home = () => {
  const [cookies] = useCookies(["userID", "token"]);
  const { items, setItems } = useItems();

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

  useEffect(() => {
    fetchItems();
  }, []);

  const signUpAndSignInPage = (
    <>
      <div>
        <Signup />
      </div>
      or
      <div>
        <Login />
      </div>
    </>
  );

  const itemListPage = (
    <MerComponent>
      <div>
        <span>
          <p>Logined User ID: {cookies.userID}</p>
        </span>
        <ItemList items={items} />
      </div>
    </MerComponent>
  );

  return <>{cookies.token ? itemListPage : signUpAndSignInPage}</>;
};
