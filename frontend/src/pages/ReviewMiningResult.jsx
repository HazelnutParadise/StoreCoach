import axios from "axios";
import { useEffect, useState } from "react";
import { Link, useParams, useNavigate } from "react-router-dom";
import FullScreenLoader from "../components/FullScreenLoader";
import "./ReviewMiningResult.css";
import RmAttributesPieChart from "../charts/RmAttributesPieChart";
import RmAttributesBarChart from "../charts/RmAttributesBarChart";

const ReviewMiningResult = ({ setPageTitle }) => {
  const navigate = useNavigate();
  const { dataUUID } = useParams();
  const [isLoading, setIsLoading] = useState(true);
  const [result, setResult] =
    useState(/*{
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
  }*/);

  const rmStoreName = result?.storeName;
  const rmProductName = result?.productName;
  const rmAttributes = result?.attributes;
  const rmResults = result?.results;
  let rmSummary = result?.summary;
  const rmTimestamp = result?.timestamp;

  if (rmSummary) rmSummary = rmSummary.replace(/\n/g, "<br />");

  useEffect(() => {
    if (isLoading) setPageTitle("載入中...");
    else if (rmProductName)
      setPageTitle(`${rmStoreName}的${rmProductName}的評論探勘報告`);
    else setPageTitle(`${rmStoreName}的評論探勘報告`);
  }, [setPageTitle, isLoading]);

  useEffect(() => {
    axios
      .get(`/api/review-mining/${dataUUID}`)
      .then((res) => {
        if (!res.data.storeName)
          throw new Error("No data returned from the server");
        console.table(res.data);
        setResult(res.data);
      })
      .catch((err) => {
        alert("Failed to fetch data");
        console.error(err);
        navigate("/");
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);
  if (isLoading) {
    return <FullScreenLoader />;
  } else if (result) {
    return (
      <div className="container">
        <div
          style={{
            width: "100%",
            display: "flex",
            justifyContent: "space-between",
            flexWrap: "wrap",
          }}
        >
          <h1 style={{ alignSelf: "start" }}>評論探勘報告 📊</h1>
          <div className="info-box">
            <div>
              <h2>
                商店名稱：<span>{rmStoreName}</span>
              </h2>
              {result.productName ? (
                <h3>
                  產品名稱：<span>{rmProductName}</span>
                </h3>
              ) : null}
            </div>

            <h4>
              分析時間：
              {new Date(rmTimestamp * 1000).toLocaleString("zh-TW")}
            </h4>
            <h6 className="uuid">UUID: {dataUUID}</h6>
          </div>
        </div>
        <div className="card rm-card tips">
          <h3>Tips 💡</h3>
          <h4>此頁面將永久存在，複製網址來分享您的分析！</h4>
        </div>
        <Link to="/" className="btn back-btn">
          回首頁
        </Link>

        <div className="report-box">
          <div className="chart-container">
            <div className="div1">
              <div className="chart-box card rm-card">
                <RmAttributesPieChart
                  rmAttributes={rmAttributes}
                  rmResults={rmResults}
                />
              </div>
            </div>
            <div className="div2">
              <div className="chart-box card rm-card">
                <RmAttributesBarChart
                  rmAttributes={rmAttributes}
                  rmResults={rmResults}
                />
              </div>
            </div>
          </div>
          <div className="summary-box">
            <h3 className="title">教練的小叮嚀 💪</h3>
            <div className="summary">
              <div className="card rm-card">
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
