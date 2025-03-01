import "../App.css";
import { Link } from "react-router-dom";

const PageNotFound = () => {
  return (
    <>
      <div className="container">
        <h1>404 - Page Not Found</h1>
        <p>抱歉，您所尋找的頁面不存在。您可能輸入了錯誤的網址。</p>
        <br />
        <Link to="/" className="btn">
          返回首頁
        </Link>
      </div>
      <style>{`
        p {
          font-size: 1.1rem;
        }
        .btn {
          display: inline-block;
          padding: 10px 20px;
          background-color: var(--main-color);
          color: white;
          text-decoration: none;
          border-radius: 5px;
        }
        .btn:hover {
          background-color: var(--main-hover-color);
        }
      `}</style>
    </>
  );
};

export default PageNotFound;
