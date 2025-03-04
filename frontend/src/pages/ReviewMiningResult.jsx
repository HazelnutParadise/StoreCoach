import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import FullScreenLoader from "../components/FullScreenLoader";
import "./ReviewMiningResult.css";
import RmAttributesPieChart from "../charts/RmAttributesPieChart";

const ReviewMiningResult = () => {
  const { dataUUID } = useParams();
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState({
    storeName: "好旺來餐廳",
    attributes: [
      "餐點美味度",
      "價格",
      "環境",
      "服務",
      "衛生",
      "精緻度",
      "cp值",
    ],
    results: [
      {
        reviewContent: "這家店的餐點很好吃，價格也很實惠。",
        miningResults: [
          {
            attribute: "餐點美味度",
            sentiment: "positive",
          },
          {
            attribute: "價格",
            sentiment: "positive",
          },
        ],
      },
      {
        reviewContent: "環境不錯，服務也很好。",
        miningResults: [
          {
            attribute: "環境",
            sentiment: "positive",
          },
          {
            attribute: "服務",
            sentiment: "positive",
          },
        ],
      },
      {
        reviewContent: "這家店的態度很差，也不好吃。",
        miningResults: [
          {
            attribute: "餐點美味度",
            sentiment: "negative",
          },
          {
            attribute: "服務",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "so disgusting",
        miningResults: [
          {
            attribute: "餐點美味度",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "drinks are too expensive",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "價格偏高，但是餐點很美味。",
        miningResults: [
          {
            attribute: "餐點美味度",
            sentiment: "positive",
          },
          {
            attribute: "價格",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },

      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
          {
            attribute: "精緻度",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
          {
            attribute: "精緻度",
            sentiment: "negative",
          },
          {
            attribute: "cp值",
            sentiment: "negative",
          },
        ],
      },
      {
        reviewContent: "不衛生，價格也不便宜。",
        miningResults: [
          {
            attribute: "價格",
            sentiment: "negative",
          },
          {
            attribute: "衛生",
            sentiment: "negative",
          },
        ],
      },
    ],
    summary:
      "好旺來餐廳的餐點美味度普遍受到肯定，環境與服務也有正面評價。然而，價格、衛生及部分服務態度是主要問題。\n\n主要優點：\n餐點美味度高，環境不錯，服務整體而言良好。\n\n主要缺點：\n價格偏高，衛生狀況不佳，部分服務態度差。\n\n中性評價：\n無明顯中性評價。\n\n教練的建議：\n1. 價格策略調整：重新評估菜單定價，考慮推出優惠套餐或會員制度，以提高價格競爭力。同時，針對高價位餐點，需確保其品質與份量能與價格相符。\n2. 加強衛生管理：徹底檢查並改善餐廳的衛生狀況，包括廚房清潔、餐具消毒、環境整潔等。定期進行衛生檢查，並公開透明地展示衛生措施，以贏取顧客信任。\n3. 提升服務品質：加強員工培訓，提升服務態度與專業技能。建立完善的顧客投訴處理機制，積極回應顧客意見，並及時解決問題。",
    timestamp: 1741099358,
  });

  const rmStoreName = result?.storeName;
  const rmAttributes = result?.attributes;
  const rmResults = result?.results;
  let rmSummary = result?.summary;
  const rmTimestamp = result?.timestamp;

  rmSummary = rmSummary.replace(/\n/g, "<br />");

  useEffect(() => {
    // todo: fetch data from backend
    // axios
    //   .get(`/api/review-mining/${dataUUID}`)
    //   .then((res) => {
    //     console.table(res.data);
    //     setResult(res.data);
    //     setIsLoading(false);
    //   })
    //   .catch((err) => {
    //     alert("Failed to fetch data");
    //     setIsLoading(false);
    //     // todo: redirect to home page
    //   });
  }, []);
  if (isLoading) {
    return <FullScreenLoader />;
  } else if (result) {
    return (
      <div className="container">
        <h1>評論探勘報告</h1>
        <div className="info-box">
          <h6>UUID: {dataUUID}</h6>
          <h2>
            商店名稱：<span>{result.storeName}</span>
          </h2>
          <h4>
            分析時間：
            {new Date(result.timestamp * 1000).toLocaleString("zh-TW")}
          </h4>
        </div>
        <div className="report-box">
          <div className="chart-container">
            <div className="div1">
              <div className="chart-box card">
                <RmAttributesPieChart
                  rmAttributes={rmAttributes}
                  rmResults={rmResults}
                />
              </div>
            </div>
            <div className="div2">
              <div className="card">
                <RmAttributesPieChart
                  rmAttributes={rmAttributes}
                  rmResults={rmResults}
                />
              </div>
            </div>
          </div>
          <div className="summary-box">
            <h3 className="title">教練的小叮嚀</h3>
            <div className="summary">
              <div className="chart-box card">
                <p dangerouslySetInnerHTML={{ __html: rmSummary }}></p>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }
};

export default ReviewMiningResult;
