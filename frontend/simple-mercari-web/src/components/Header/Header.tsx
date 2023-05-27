import { useCookies } from "react-cookie";
import "./Header.css";
import { useNavigate } from "react-router-dom";
import { FaCamera, FaHome, FaUser, FaSearch } from "react-icons/fa";
import { Button, Heading, Tabs, Form } from "react-bulma-components";
import "bulma/css/bulma.min.css";
import { useState } from "react";
import { fetcher } from "../../helper";
import { Item } from "../../common/context";
import { useItems } from "../../common/context";
import { toast } from "react-toastify";

export const Header: React.FC = () => {
  const [cookies, _, removeCookie] = useCookies(["userID", "token"]);
  const [keyWord, setKeyWord] = useState("");

  const onLogout = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    removeCookie("userID");
    removeCookie("token");
  };

  const navigate = useNavigate();

  const { setItems } = useItems();

  const tabs = [
    { name: "Home", icon: <FaHome />, to: "/" },
    { name: "Listing", icon: <FaCamera />, to: "/sell" },
    { name: "MyPage", icon: <FaUser />, to: `/user/${cookies.userID}` },
  ];

  const handleTabClick = (tabName: string, to: string) => {
    setActiveTab(tabName);
    localStorage.setItem("activeTab", tabName);
    navigate(to);
  };

  const search = async () => {
    if (keyWord) {
      fetcher<Item[]>(`/search?name=${keyWord}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Accept: "application/json",
        },
      })
        .then((data) => {
          console.log("GET success:", data);
          setItems(data);
          localStorage.setItem("searchResult", JSON.stringify(data));
        })
        .catch((err) => {
          console.log(`GET error:`, err);
          toast.error(err.message);
        });
    }
  };

  const [activeTab, setActiveTab] = useState(() => {
    return localStorage.getItem("activeTab") || "Home";
  });

  return (
    <>
      <header>
        <div className="titleWrapper">
          <Heading className="is-size-3-desktop is-size-4-mobile">
            Simple Mercari
          </Heading>
          <Form.Control>
            <Form.Input
              placeholder="Search items"
              type="text"
              onChange={(e) => setKeyWord(e.target.value)}
            />
          </Form.Control>
          <Button color="" onClick={search} className="is-responsive">
            <FaSearch />
          </Button>
        </div>
        <div className="LogoutButtonContainer">
          <Button color="" onClick={onLogout} className="is-responsive">
            Logout
          </Button>
        </div>
        <Tabs type="boxed" className="is-centered is-medium">
          {tabs.map((tab) => {
            return (
              <Tabs.Tab
                active={activeTab === tab.name}
                onClick={() => handleTabClick(tab.name, tab.to)}
              >
                <div className="tabItem">
                  {tab.icon}
                  <span>{tab.name}</span>
                </div>
              </Tabs.Tab>
            );
          })}
        </Tabs>
      </header>
    </>
  );
};
