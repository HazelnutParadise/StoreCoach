import Logo from "../assets/react.svg";
import "./TopBar.css";

const TopBar = () => (
  <div className="top-bar">
    <img src={Logo} alt="logo" />
    <h1>StoreCoach</h1>
  </div>
);

export default TopBar;
