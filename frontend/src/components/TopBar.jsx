import Logo from "../assets/react.svg";
import "./TopBar.css";

const TopBar = () => (
  <div className="top-bar">
    <img src={Logo} alt="logo" />
    <h2 className="title">StoreCoach 商店教練</h2>
  </div>
);

export default TopBar;
