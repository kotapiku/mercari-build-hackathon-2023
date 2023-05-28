import React, { useEffect, useState } from "react";
import { useCookies } from "react-cookie";
import { MerComponent } from "../MerComponent";
import { toast } from "react-toastify";
import { fetcher } from "../../helper";
import { Button, Form } from "react-bulma-components";

interface Category {
  id: number;
  name: string;
}

type formDataType = {
  name: string;
  category_id: number;
  price: number;
  description: string;
  image: string | File;
};

export const Listing: React.FC = () => {
  const initialState = {
    name: "",
    category_id: 1,
    price: 0,
    description: "",
    image: "",
  };
  const [values, setValues] = useState<formDataType>(initialState);
  const [categories, setCategories] = useState<Category[]>([]);
  const [cookies] = useCookies(["token", "userID"]);

  const onValueChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value,
    });
  };

  const onSelectChange = (event: React.ChangeEvent<HTMLSelectElement>) => {
    setValues({
      ...values,
      [event.target.name]: event.target.value,
    });
  };

  const onFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setValues({
      ...values,
      [event.target.name]: event.target.files![0],
    });
  };

  const onSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const data = new FormData();
    data.append("name", values.name);
    data.append("category_id", values.category_id.toString());
    data.append("price", values.price.toString());
    data.append("description", values.description);
    data.append("image", values.image);

    fetcher<{ id: number }>(`/items`, {
      method: "POST",
      body: data,
      headers: {
        Authorization: `Bearer ${cookies.token}`,
      },
    })
      .then((res) => {
        sell(res.id);
      })
      .catch((error: Error) => {
        toast.error(error.message);
        console.error("POST error:", error);
      });
  };

  const sell = (itemID: number) =>
    fetcher(`/sell`, {
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        Authorization: `Bearer ${cookies.token}`,
      },
      body: JSON.stringify({
        item_id: itemID,
      }),
    })
      .then((_) => {
        toast.success("Item added successfully!");
      })
      .catch((error: Error) => {
        toast.error(error.message);
        console.error("POST error:", error);
      });

  const fetchCategories = () => {
    fetcher<Category[]>(`/items/categories`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    })
      .then((items) => setCategories(items))
      .catch((err) => {
        console.log(`GET error:`, err);
        toast.error(err.message);
      });
  };

  const createDescription = (event: React.MouseEvent<HTMLElement>) => {
    event.preventDefault();

    fetcher<string>(`/description`, {
      method: "POST",
      body: JSON.stringify({
        name: values.name,
        description: values.description,
      }),
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
    })
      .then((res) => {
        console.log(res);
        setValues({ ...values, description: res });
      })
      .catch((error: Error) => {
        toast.error(error.message);
        console.error("POST error:", error);
      });
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  return (
    <MerComponent>
      <div className="columns is-centered">
        <div className="column is-4">
          <form onSubmit={onSubmit}>
            <Form.Field>
              <Form.Label>Item name</Form.Label>
              <Form.Control>
                <Form.Input
                  type="text"
                  name="name"
                  id="MerTextInput"
                  placeholder="name"
                  value={values.name}
                  onChange={onValueChange}
                  required
                />
              </Form.Control>
            </Form.Field>
            <Form.Field>
              <Form.Control>
                <Form.Label>Category</Form.Label>
                <Form.Select
                  name="category_id"
                  value={values.category_id}
                  onChange={onSelectChange}
                >
                  {categories &&
                    categories.map((category) => {
                      return (
                        <option value={category.id} key={category.id}>
                          {category.name}
                        </option>
                      );
                    })}
                </Form.Select>
              </Form.Control>
            </Form.Field>
            <Form.Field>
              <Form.Control>
                <Form.Label>Price</Form.Label>
                <Form.Input
                  type="number"
                  name="price"
                  value={values.price}
                  min="0"
                  onChange={onValueChange}
                  required
                />
              </Form.Control>
            </Form.Field>
            <Form.Field>
              <Form.Control>
                <Form.Label>Description</Form.Label>
                <Form.Input
                  type="text"
                  name="description"
                  value={values.description}
                  placeholder="description"
                  onChange={onValueChange}
                  required
                />
                <Button onClick={createDescription}>Create by ChatGPT</Button>
              </Form.Control>
            </Form.Field>
            <Form.Field>
              <Form.Control>
                <Form.Label>Image file</Form.Label>
                <Form.Input
                  type="file"
                  name="image"
                  onChange={onFileChange}
                  required
                />
              </Form.Control>
            </Form.Field>
            <Button type="submit" id="MerButton">
              List this item
            </Button>
          </form>
        </div>
      </div>
    </MerComponent>
  );
};
