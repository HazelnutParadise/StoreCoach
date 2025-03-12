import { useEffect, useState } from "react";
import { Link, useParams, useNavigate } from "react-router-dom";
import FullScreenLoader from "../components/FullScreenLoader";
import "./ReviewMiningResult.css";
import RmAttributesPieChart from "../charts/RmAttributesPieChart";
import RmAttributesBarChart from "../charts/RmAttributesBarChart";
import RmTtest from "../charts/RmTtest";

const testData = {
  storeName: "好旺來餐廳",
  attributes: ["餐點美味度", "價格", "環境", "服務", "衛生", "精緻度", "cp值"],
  tTest: {
    餐點美味度: {
      pValue: 0.1,
      tValue: 5.2,
      df: 10,
    },
    價格: {
      pValue: 0.25,
      tValue: 5.2,
      df: 10,
    },
    環境: {
      pValue: 0.0033,
      tValue: 5.2,
      df: 10,
    },
    服務: {
      pValue: 0.0001,
      tValue: 5.2,
      df: 10,
    },
    衛生: {
      pValue: 0.0001,
      tValue: 5.2,
      df: 10,
    },
    精緻度: {
      pValue: 0.0001,
      tValue: 5.2,
      df: 10,
    },
    cp值: {
      pValue: 0.0001,
      tValue: 5.2,
      df: 10,
    },
  },
  results: [
    {
      reviewContent: "這家店的餐點很好吃，價格也很實惠。",
      reviewRating: 5,
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
      reviewRating: 4,
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
      reviewRating: 1,
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
      reviewRating: 1,
      miningResults: [
        {
          attribute: "餐點美味度",
          sentiment: "negative",
        },
      ],
    },
    {
      reviewContent: "drinks are too expensive",
      reviewRating: 2,
      miningResults: [
        {
          attribute: "價格",
          sentiment: "negative",
        },
      ],
    },
    {
      reviewContent: "價格偏高，但是餐點很美味。",
      reviewRating: 4,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
      reviewRating: 2,
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
};

const ReviewMiningResult = ({ setPageTitle }) => {
  const navigate = useNavigate();
  const { dataUUID } = useParams();
  const [isLoading, setIsLoading] = useState(true);
  const [result, setResult] = useState();

  const rmStoreName = result?.storeName;
  const rmProductName = result?.productName;
  const rmAttributes = result?.attributes;
  const rmResults = result?.results;
  let rmSummary = result?.summary;
  const rmTimestamp = result?.timestamp;
  let rmRatings = result?.results.map((r) => r.reviewRating);
  const rmTtest = result?.tTest;

  let attributeCount = {};
  if (result) {
    rmAttributes.forEach((element) => {
      attributeCount[element] = 0;
    });
    // 計算每個屬性出現的次數
    rmResults.forEach((result) => {
      if(!result.miningResults) {
        alert("找不到探勘結果，可能是評論資料沒有內容");
        navigate("/");
        return;
      };
      const miningResults = result.miningResults;
      miningResults.forEach((miningResult) => {
        if (miningResult.attribute in attributeCount) {
          attributeCount[miningResult.attribute]++;
        }
      });
    });
    rmAttributes.sort((a, b) => attributeCount[b] - attributeCount[a]);
  }

  let rating0count = 0;
  if (rmRatings) {
    rmRatings.forEach((r) => {
      if (r === 0.0) rating0count++;
    });
    if (rating0count === rmRatings.length) rmRatings = null;
  }

  if (rmSummary) rmSummary = rmSummary.replace(/\n/g, "<br />");

  useEffect(() => {
    if (isLoading) {
      setPageTitle("正在探勘評論...");
    } else if (rmProductName) {
      setPageTitle(`${rmStoreName}的${rmProductName}的評論探勘報告`);
    } else {
      setPageTitle(`${rmStoreName}的評論探勘報告`);
    }
  }, [setPageTitle, isLoading]);

  useEffect(() => {
    let isMounted = true;

    const fetchData = async () => {
      const failedTimeout = 10 * 60 * 1000;
      const startTime = new Date().getTime();

      try {
        while (isMounted) {
          try {
            const res = await fetch(`/api/review-mining/${dataUUID}`, {
              headers: {
                "Cache-Control": "no-cache",
                "Content-Type": "application/json",
                Accept: "application/json",
              },
            });
            const resJson = await res.json();

            if (!resJson.storeName) {
              throw new Error("No data returned from the server");
            }

            if (isMounted) {
              setResult(resJson);
              break;
            }
          } catch (err) {
            const nowTime = new Date().getTime();
            if (nowTime - startTime > failedTimeout) {
              if (isMounted) {
                alert("Failed to fetch data");
                console.error(err);
                navigate("/");
              }
              break;
            }
            await new Promise((resolve) => setTimeout(resolve, 10 * 1000));
          }
        }
      } finally {
        if (isMounted) {
          setIsLoading(false);
        }
      }
    };

    fetchData();

    return () => {
      isMounted = false;
    };
  }, [dataUUID, navigate]);
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
            {rmRatings ? (
              <div className="div3" style={{ width: "100%" }}>
                <div className="card rm-card">
                  <h3>顧客評分分布</h3>
                  <div className="rating-container">
                    <div className="rating-box">
                      <h4
                        style={{
                          padding: 0,
                          margin: 0,
                          fontSize: "calc(0.9rem + 0.7vw)",
                        }}
                      >
                        最高評分
                      </h4>
                      <span style={{ fontSize: "calc(1.1rem + 0.7vw)" }}>
                        {Math.max(...rmRatings)}
                      </span>
                    </div>
                    <div className="rating-box">
                      <h4
                        style={{
                          padding: 0,
                          margin: 0,
                          fontSize: "calc(0.9rem + 0.7vw)",
                        }}
                      >
                        評分中位數
                      </h4>
                      <span style={{ fontSize: "calc(1.1rem + 0.7vw)" }}>
                        {((arr) => {
                          arr = arr.sort((a, b) => a - b);
                          let median;
                          if (arr.length % 2 === 0) {
                            // 數目為偶數
                            median =
                              (arr[arr.length / 2] + arr[arr.length / 2 - 1]) /
                              2;
                          } else {
                            // 數目為奇數
                            median = arr[(arr.length - 1) / 2];
                          }
                          return median;
                        })(rmRatings)}
                      </span>
                    </div>
                    <div className="rating-box">
                      <h4
                        style={{
                          padding: 0,
                          margin: 0,
                          fontSize: "calc(0.9rem + 0.7vw)",
                        }}
                      >
                        最低評分
                      </h4>
                      <span style={{ fontSize: "calc(1.1rem + 0.7vw)" }}>
                        {Math.min(...rmRatings)}
                      </span>
                    </div>
                    <div className="rating-box">
                      <h4
                        style={{
                          padding: 0,
                          margin: 0,
                          fontSize: "calc(0.9rem + 0.7vw)",
                        }}
                      >
                        平均評分
                      </h4>
                      <span style={{ fontSize: "calc(1.1rem + 0.7vw)" }}>
                        {(
                          rmRatings.reduce((a, b) => a + b, 0) /
                          rmRatings.length
                        ).toFixed(2)}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            ) : null}
            {rmTtest ? (
              <div className="card rm-card">
                <RmTtest rmAttributes={rmAttributes} rmTtest={rmTtest} />
              </div>
            ) : null}
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
