import { useEffect, useRef, useState } from "react";
import { Route, Routes, useLocation, useNavigate } from "react-router-dom";
import axios from "axios";
import {
  ArcElement,
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Legend,
  Tooltip,
} from "chart.js";
import TopBar from "./components/TopBar";
import ReviewsInput from "./components/ReviewsInput";
import ReviewMiningResult from "./pages/ReviewMiningResult";
import PageNotFound from "./pages/404";
import "./App.css";
import TextAD from "./components/TextAD";

const App = () => {
  ChartJS.register(
    ArcElement,
    Tooltip,
    Legend,
    CategoryScale,
    LinearScale,
    BarElement
  );
  ChartJS.defaults.backgroundColor = [
    "rgba(255, 99, 132, 1)",
    "rgba(54, 162, 235, 1)",
    "rgba(255, 206, 86, 1)",
    "rgb(75, 192, 112)",
    "rgba(153, 102, 255, 1)",
    "rgba(255, 159, 64, 1)",
  ];

  ChartJS.defaults.borderWidth = 1;
  ChartJS.defaults.hoverBorderColor = "#46237a";
  ChartJS.defaults.hoverBorderWidth = 1.2;

  const [title, setTitle] = useState();
  const location = useLocation();
  const navigate = useNavigate();
  const setPageTitle = (str) => {
    const newTitle = str
      ? str + " | StoreCoach 商店教練 - 榛果繽紛樂"
      : "StoreCoach 商店教練 - 榛果繽紛樂";
    setTitle(newTitle);
  };
  useEffect(() => {
    setPageTitle();
  }, [location]);

  useEffect(() => {
    document.title = title;
  }, [title]);

  useEffect(() => {
    getNavbar();
  }, []);

  const currentYear = new Date().getFullYear();

  const spreadsheetRef = useRef(null);
  const [storeName, setStoreName] = useState("");
  const [productName, setProductName] = useState("");

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
          storeName,
          productName,
          reviews,
        })
        .then((res) => {
          const uuid = res.data.dataUUID;
          console.log("data_uuid:", uuid);
          navigate(`/review-mining/${uuid}`);
        });
    }
  };

  return (
    <>
      <TopBar />
      <Routes>
        <Route
          path="/"
          element={
            <div className="container">
              <TextAD />
              <div className="card">
                <h2>AI 評論探勘</h2>
                <br />
                <div className="input-container">
                  <div
                    style={{
                      display: "flex",
                      justifyContent: "center",
                      paddingBottom: "20px",
                    }}
                  >
                    <a
                      className="btn"
                      style={{ padding: "8px 12px", borderRadius: "50px" }}
                      href="https://gmaps-reviews.hazelnut-paradise.com/?mode=storecoach"
                    >
                      使用 <strong>估咩評論小扒手</strong> 爬取評論
                    </a>
                  </div>

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
                  <div className="input-box">
                    <label htmlFor="product-name">商品名稱(選填)</label>
                    <input
                      type="text"
                      id="product-name"
                      onChange={(e) => setProductName(e.target.value)}
                    />
                  </div>
                </div>
                <br />
                <ReviewsInput
                  ref={spreadsheetRef}
                  columnLabels={["輸入或貼上評論內容 (列數會自動擴張)："]}
                />
                <button onClick={handleStartReviewMining}>開始探勘</button>
              </div>
              <div className="tips card">
                此專案採用 Apache 2.0 License 開源，歡迎前往 Github 查看原始碼：
                <a
                  target="_blank"
                  href="https://github.com/HazelnutParadise/StoreCoach"
                  style={{wordBreak: "break-all"}}
                >
                  https://github.com/HazelnutParadise/StoreCoach
                </a>
              </div>
            </div>
          }
        />
        <Route
          path="/review-mining/:dataUUID"
          element={<ReviewMiningResult setPageTitle={setPageTitle} />}
        />
        <Route
          path="*"
          element={<PageNotFound setPageTitle={setPageTitle} />}
        />
      </Routes>
      <footer>
        <div id="navbar-placeholder"></div>
        <div>{currentYear} © HazelnutParadise</div>
      </footer>
    </>
  );
};

export default App;
