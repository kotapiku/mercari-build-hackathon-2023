import { useState, useEffect } from "react";
import { useCookies } from "react-cookie";
import { useNavigate } from "react-router-dom";
import { fetcherBlob } from "../../helper";
import "bulma/css/bulma.min.css";
import { Card, Media, Heading, Content, Columns } from "react-bulma-components";
import { FaYenSign, FaHashtag } from "react-icons/fa";

interface Item {
  id: number;
  name: string;
  price: number;
  category_name: string;
  status: number;
}
export const Item: React.FC<{ item: Item }> = ({ item }) => {
  const navigate = useNavigate();
  const [itemImage, setItemImage] = useState<string>("");
  const [cookies] = useCookies(["token"]);

  async function getItemImage(itemId: number): Promise<Blob> {
    return await fetcherBlob(`/items/${itemId}/image`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
    });
  }

  useEffect(() => {
    async function fetchData() {
      const image = await getItemImage(item.id);
      setItemImage(URL.createObjectURL(image));
    }

    fetchData();
  }, [item]);

  return (
    <Columns.Column size={4}>
      <Card className={item.status == 3 ? "SoldOut" : "OnSale"}>
        <Card.Image
          size="1by1"
          src={itemImage}
          onClick={() => navigate(`/item/${item.id}`)}
        />
        <Card.Content className={item.status == 3 ? "SoldOutContent" : ""}>
          <Media>
            <Media.Item>
              <Heading size={4}>{item.name}</Heading>
            </Media.Item>
            <Media.Item>
              <Content>
                <p>
                  <span>
                    <FaHashtag style={{ verticalAlign: -2 }} />
                    {item.category_name}
                  </span>
                  <br />
                  <span>
                    <FaYenSign style={{ verticalAlign: -2 }} />
                    {item.price}
                  </span>
                  <br />
                </p>
              </Content>
            </Media.Item>
          </Media>
        </Card.Content>
      </Card>
    </Columns.Column>
  );
};
