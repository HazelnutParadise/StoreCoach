import { useEffect, useRef, useState } from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import axios from "axios";
import TopBar from "./components/TopBar";
import ReviewsInput from "./components/ReviewsInput";
import "./App.css";

const App = () => {
  useEffect(() => {
    getNavbar();
  }, []);

  const currentYear = new Date().getFullYear();

  const spreadsheetRef = useRef(null);
  let [storeName, setStoreName] = useState("");

  // 取得表格數據（會自動忽略空白列）
  const handleStartReviewMining = () => {
    if (!storeName) {
      alert("請輸入商店名稱");
      return;
    }
    if (spreadsheetRef.current) {
      const reviews = spreadsheetRef.current.getData();
      console.log("表格資料:", reviews);
      if (reviews.length === 0) {
        alert("請輸入評論內容");
        return;
      }
      axios
        .post("/api/review-mining", {
          store_name: storeName,
          reviews: reviews,
        })
        .then((res) => {
          console.log("uuid:", res.data);
        });
    }
  };

  return (
    <Router>
      <TopBar />
      <Routes>
        <Route
          path="/"
          element={
            <div className="container">
              <div className="card">
                <h2>AI 評論探勘</h2>
                <br />
                <div className="input-box">
                  <label htmlFor="store-name">商店名稱</label>
                  <input
                    type="text"
                    id="store-name"
                    value={storeName}
                    onChange={(e) => setStoreName(e.target.value)}
                  />
                </div>
                <br />
                <ReviewsInput
                  ref={spreadsheetRef}
                  columnLabels={["輸入或貼上評論內容 (列數會自動擴張)："]}
                />
                <button onClick={handleStartReviewMining}>開始探勘</button>
              </div>
            </div>
          }
        />
      </Routes>
      <footer>
        <div id="navbar-placeholder"></div>
        <div>{currentYear} © HazelnutParadise</div>
      </footer>
    </Router>
  );
};

export default App;
