import { useRef } from "react";
import TopBar from "./components/TopBar";
import ReviewsInput from "./components/ReviewsInput";
import "./App.css";

const App = () => {
  const spreadsheetRef = useRef(null);

  // 取得表格數據（會自動忽略空白列）
  const handleStartReviewMining = () => {
    if (spreadsheetRef.current) {
      console.log("表格資料:", spreadsheetRef.current.getData());
    }
  };

  return (
    <div className="App">
      <TopBar />
      <div className="container">
        <div className="card">
          <h2>AI 評論探勘</h2>
          <ReviewsInput ref={spreadsheetRef} columnLabels={["輸入評論內容"]} />
          <button onClick={handleStartReviewMining}>開始探勘</button>
        </div>
      </div>
    </div>
  );
};

export default App;
