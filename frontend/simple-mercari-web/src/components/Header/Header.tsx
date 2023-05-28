import { useCookies } from "react-cookie";
import "./Header.css";
import { useLocation, useNavigate } from "react-router-dom";
import { FaCamera, FaHome, FaUser, FaSearch } from "react-icons/fa";
import { Button, Heading, Tabs, Form } from "react-bulma-components";
import "bulma/css/bulma.min.css";
import { useEffect, useState } from "react";

export const Header: React.FC = () => {
  const [cookies, _, removeCookie] = useCookies(["userID", "token"]);
  const [keyWord, setKeyWord] = useState("");
  const [activeTab, setActiveTab] = useState("Home");

  const navigate = useNavigate();

  const onLogout = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    removeCookie("userID", { path: "/" });
    removeCookie("token", { path: "/" });
  };

  const tabs = [
    { name: "Home", icon: <FaHome />, to: "/" },
    { name: "Listing", icon: <FaCamera />, to: "/sell" },
    { name: "MyPage", icon: <FaUser />, to: `/user/${cookies.userID}` },
  ];

  const handleTabClick = (tabName: string, to: string) => {
    navigate(to);
  };

  const search = async () => {
    if (keyWord) {
      navigate(`/search?keyword=${keyWord}`);
    }
  };

  const location = useLocation();

  useEffect(() => {
    // identify active tab and set it to state
    const currentTab = tabs.find((tab) => tab.to === location.pathname);
    if (currentTab) {
      setActiveTab(currentTab.name);
    } else {
      setActiveTab("Home");
    }
  }, [location]);

  const tab = (
    <Tabs type="boxed" className="is-centered is-medium">
      {tabs.map((tab) => {
        return (
          <Tabs.Tab
            active={activeTab === tab.name}
            onClick={() => handleTabClick(tab.name, tab.to)}
            key={tab.name}
          >
            <div className="tabItem">
              {tab.icon}
              <span>{tab.name}</span>
            </div>
          </Tabs.Tab>
        );
      })}
    </Tabs>
  );

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
        {cookies.token ? tab : <></>}
      </header>
    </>
  );
};
