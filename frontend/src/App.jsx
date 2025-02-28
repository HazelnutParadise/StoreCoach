import { useEffect, useRef } from "react";
import TopBar from "./components/TopBar";
import ReviewsInput from "./components/ReviewsInput";
import "./App.css";

const App = () => {
  useEffect(() => {
    getNavbar();
  }, []);

  const currentYear = new Date().getFullYear();

  const spreadsheetRef = useRef(null);

  // 取得表格數據（會自動忽略空白列）
  const handleStartReviewMining = () => {
    if (spreadsheetRef.current) {
      console.log("表格資料:", spreadsheetRef.current.getData());
    }
  };

  return (
    <>
      <TopBar />
      <div className="container">
        <div className="card">
          <h2>AI 評論探勘</h2>
          <ReviewsInput
            ref={spreadsheetRef}
            columnLabels={["輸入或貼上評論內容 (列數會自動擴張)："]}
          />
          <button onClick={handleStartReviewMining}>開始探勘</button>
        </div>
      </div>
      <footer>
        <div id="navbar-placeholder"></div>
        <div>{currentYear} © HazelnutParadise</div>
      </footer>
    </>
  );
};

export default App;
