import { useEffect } from "react";
import "./FullScreenLoader.css";

const FullScreenLoader = () => {
  return (
    <>
      <div className="fullscreen-loader">
        <img
          className="logo"
          src="https://src.hazelnut-paradise.com/StoreCoach-logo.png"
        />
        <div className="loader"></div>
        <h2>LOADING</h2>
      </div>
    </>
  );
};

export default FullScreenLoader;
