import { FaHome, FaCamera, FaUser } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { useCookies } from "react-cookie";
import "./Footer.css";

export const Footer: React.FC = () => {
  const [cookies] = useCookies(["userID"]);
  const navigate = useNavigate();

  if (!cookies.userID) {
    return <></>;
  }

  return (
    <footer>
      <div className="MerFooterItem" onClick={() => navigate("/")}>
        <FaHome />
        <p>Home</p>
      </div>
      <div className="MerFooterItem" onClick={() => navigate("/sell")}>
        <FaCamera />
        <p>Listing</p>
      </div>
      <div
        className="MerFooterItem"
        onClick={() => navigate(`/user/${cookies.userID}`)}
      >
        <FaUser />
        <p>MyPage</p>
      </div>
    </footer>
  );
};
