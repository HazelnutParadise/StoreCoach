import { use, useEffect } from "react";
import { Link } from "react-router-dom";

const PageNotFound = (props) => {
  const { setPageTitle } = props;
  useEffect(() => {
    setPageTitle("404 Page Not Found");
  }, [setPageTitle]);
  return (
    <>
      <div className="container">
        <h2>
          <span>404</span>
          <br />
          Page Not Found
        </h2>
        <p>抱歉，您所尋找的頁面不存在。您可能輸入了錯誤的網址。</p>
        <br />
        <Link to="/" className="btn">
          返回首頁
        </Link>
      </div>
      <style>{`
        span {
          font-size: 6.5rem;
        }
        h2 {
          margin-top: 0;
          font-size: 2rem;
          text-align: center;
          margin-bottom: 10px;
        }
        p {
          font-size: 1.1rem;
        }
      `}</style>
    </>
  );
};

export default PageNotFound;
