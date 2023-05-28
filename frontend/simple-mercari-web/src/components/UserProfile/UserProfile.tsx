import { useState, useEffect } from "react";
import { useCookies } from "react-cookie";
import { useParams } from "react-router-dom";
import { toast } from "react-toastify";
import { MerComponent } from "../MerComponent";
import { ItemList } from "../ItemList";
import { fetcher } from "../../helper";
import { Button, Icon, Form } from "react-bulma-components";
import { FaYenSign } from "react-icons/fa";

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

export const UserProfile: React.FC = () => {
  const [items, setItems] = useState<Item[]>([]);
  const [balance, setBalance] = useState<number>();
  const [addedbalance, setAddedBalance] = useState<number>();
  const [cookies] = useCookies(["token"]);
  const params = useParams();

  const fetchItems = () => {
    fetcher<Item[]>(`/users/${params.id}/items`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
    })
      .then((items) => setItems(items))
      .catch((err) => {
        console.log(`GET error:`, err);
        toast.error(err.message);
      });
  };

  const fetchUserBalance = () => {
    fetcher<{ balance: number }>(`/balance`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
    })
      .then((res) => {
        setBalance(res.balance);
      })
      .catch((err) => {
        console.log(`GET error:`, err);
        toast.error(err.message);
      });
  };

  useEffect(() => {
    fetchItems();
    fetchUserBalance();
  }, []);

  const onBalanceSubmit = () => {
    fetcher(`/balance`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
      body: JSON.stringify({
        balance: addedbalance,
      }),
    })
      .then((_) => window.location.reload())
      .catch((err) => {
        console.log(`POST error:`, err);
        toast.error("Please enter a number greater than 0");
      });
  };

  return (
    <MerComponent>
      <div className="columns is-centered">
        <div className="column is-4">
          <Form.Field>
            <Form.Label>Balance: {balance}</Form.Label>
            <Form.Control>
              <Form.Input
                type="number"
                name="balance"
                placeholder="0"
                min="0"
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                  setAddedBalance(Number(e.target.value));
                }}
                required
              />
              <Icon align="left" size="small">
                <FaYenSign style={{ verticalAlign: -2 }} />
              </Icon>
            </Form.Control>
            <Button onClick={onBalanceSubmit} id="MerButton">
              Add balance
            </Button>
          </Form.Field>
        </div>
      </div>

      <ItemList items={items} />
    </MerComponent>
  );
};
