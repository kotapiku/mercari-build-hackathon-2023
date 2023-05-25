import { useCookies } from "react-cookie";
import "./Header.css";
import { useNavigate } from "react-router-dom";
import { FaCamera, FaHome, FaUser } from "react-icons/fa";
import { Button, Heading, Tabs } from "react-bulma-components";
import "bulma/css/bulma.min.css";
import { useState } from "react";
export const Header: React.FC = () => {
  const [cookies, _, removeCookie] = useCookies(["userID", "token"]);

  const onLogout = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    removeCookie("userID");
    removeCookie("token");
  };

  const navigate = useNavigate();

  const tabs = [
    { name: "Home", icon: <FaHome />, to: "/" },
    { name: "Listing", icon: <FaCamera />, to: "/sell" },
    { name: "MyPage", icon: <FaUser />, to: `/user/${cookies.userID}` },
  ];

  const handleTabClick = (tabName: string, to: string) => {
    setActiveTab(tabName);
    navigate(to);
  };

  const [activeTab, setActiveTab] = useState("Home");

  return (
    <>
      <header>
        <div className="titleWrapper">
          <Heading className="is-size-2-desktop is-size-4-mobile">
            Simple Mercari
          </Heading>
        </div>
        <div className="LogoutButtonContainer">
          <Button color="" onClick={onLogout} className="is-responsive">
            Logout
          </Button>
        </div>
        <div className="tabsWrapper">
          <Tabs>
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
        </div>
      </header>
    </>
  );
};
