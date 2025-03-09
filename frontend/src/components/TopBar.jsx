import "./TopBar.css";

const TopBar = () => (
  <div style={{ position: "sticky", top: 0, zIndex: 100 }}>
    <div className="top-bar">
      <img
        src="https://src.hazelnut-paradise.com/StoreCoach-logo.png"
        alt="logo"
        height={"100%"}
        width={"auto"}
      />
      <h2 className="title">StoreCoach 商店教練</h2>
    </div>
    {/* <div id="Pistachio-Announcement" style={{zIndex: 99999}}></div> */}
  </div>
);

export default TopBar;
