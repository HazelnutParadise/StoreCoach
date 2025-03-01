import Logo from "../assets/react.svg";
import "./TopBar.css";

const TopBar = () => (
  <div className="top-bar">
    <img src={Logo} alt="logo" />
    <h1>
      <i>StoreCoach 商店教練</i>
    </h1>
  </div>
);

export default TopBar;
