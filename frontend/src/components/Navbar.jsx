import { useEffect, useRef } from "react";

const Navbar = () => {
  const navbarPlaceholder = useRef(null);

  const loadNavbar = () => {
    fetch(
      "https://src.hazelnut-paradise.com/navbar.html?content-type=text/html",
      {
        method: "GET",
      }
    )
      .then((response) => response.text())
      .then((html) => {
        if (navbarPlaceholder.current) {
          navbarPlaceholder.current.innerHTML = html;
        }
      })
      .catch((error) => {
        console.error("Error loading navbar:", error);
      });
  };

  useEffect(() => {
    loadNavbar();
  }, []);

  return <div ref={navbarPlaceholder} id="navbar-placeholder"></div>;
};

export default Navbar;
