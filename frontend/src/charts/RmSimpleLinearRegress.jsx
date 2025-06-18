const RmSimpleLinearRegress = ({ rmAttributes, rmSimpleLinearRegress }) => {
  let attributes = rmAttributes;
  attributes = attributes.slice(0, 5);

  // 計算顯著屬性的影響程度排名
  const significantAttributes = attributes
    .filter(
      (attr) =>
        rmSimpleLinearRegress[attr] &&
        rmSimpleLinearRegress[attr].pValue <= 0.05
    )
    .map((attr) => ({
      attribute: attr,
      slope: Math.abs(rmSimpleLinearRegress[attr].slope), // 使用絕對值來排名影響程度
      originalSlope: rmSimpleLinearRegress[attr].slope,
      r2: rmSimpleLinearRegress[attr].r2,
      pValue: rmSimpleLinearRegress[attr].pValue,
    }))
    .sort((a, b) => b.slope - a.slope); // 按影響程度降序排列

  // 為每個顯著屬性分配排名
  const attributeRanks = {};
  significantAttributes.forEach((item, index) => {
    attributeRanks[item.attribute] = index + 1;
  });

  // 定義排名樣式
  const getRankStyle = (rank) => {
    const styles = {
      1: {
        borderColor: "#ffd700",
        backgroundColor: "#fffbf0",
        badgeColor: "#ffd700",
        textColor: "#b8860b",
      },
      2: {
        borderColor: "#c0c0c0",
        backgroundColor: "#f8f9fa",
        badgeColor: "#c0c0c0",
        textColor: "#6c757d",
      },
      3: {
        borderColor: "#cd7f32",
        backgroundColor: "#fef9f6",
        badgeColor: "#cd7f32",
        textColor: "#8b4513",
      },
      4: {
        borderColor: "#20c997",
        backgroundColor: "#f0fdf9",
        badgeColor: "#20c997",
        textColor: "#059669",
      },
      5: {
        borderColor: "#6f42c1",
        backgroundColor: "#f8f6ff",
        badgeColor: "#6f42c1",
        textColor: "#5a2d91",
      },
    };
    return styles[rank] || styles[5];
  };

  return (
    <>
      {" "}
      <h3>
        提及次數前五多的屬性得分分別對商家評分的影響（簡單線性回歸分析R²值與顯著性，顯著水準=0.05）
      </h3>
      <p style={{ fontSize: "0.9rem", color: "#6c757d", marginBottom: "1rem" }}>
        顯著屬性按影響程度排名，排名徽章顏色：🥇金牌、🥈銀牌、🥉銅牌、第4名、第5名
      </p>
      <div className="ttest-container">
        {attributes.map((attribute) => {
          const isSignificant =
            rmSimpleLinearRegress[attribute] &&
            rmSimpleLinearRegress[attribute].pValue <= 0.05;
          const rank = attributeRanks[attribute];
          const rankStyle = isSignificant && rank ? getRankStyle(rank) : null;

          return (
            <div
              className="ttest-box"
              key={attribute}
              style={
                rankStyle
                  ? {
                      borderColor: rankStyle.borderColor,
                      borderWidth: "2px",
                      backgroundColor: rankStyle.backgroundColor,
                      position: "relative",
                    }
                  : {}
              }
            >
              {/* 排名徽章 */}
              {isSignificant && rank && (
                <div
                  style={{
                    position: "absolute",
                    top: "-8px",
                    left: "8px",
                    backgroundColor: rankStyle.badgeColor,
                    color: "white",
                    borderRadius: "50%",
                    width: "24px",
                    height: "24px",
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center",
                    fontSize: "0.8rem",
                    fontWeight: "bold",
                    boxShadow: "0 2px 4px rgba(0,0,0,0.2)",
                  }}
                >
                  {rank}
                </div>
              )}
              <h4
                style={{
                  padding: 0,
                  margin: 0,
                  fontSize: "calc(0.9rem + 0.7vw)",
                  color: rankStyle ? rankStyle.textColor : "inherit",
                  fontWeight: rankStyle ? "600" : "normal",
                }}
              >
                {attribute}
                {isSignificant && rank && (
                  <span
                    style={{
                      fontSize: "0.7rem",
                      fontWeight: "normal",
                      color: "#6c757d",
                      marginLeft: "0.5rem",
                    }}
                  >
                    (影響程度第{rank}名)
                  </span>
                )}
              </h4>{" "}
              {!rmSimpleLinearRegress[attribute] ||
              (rmSimpleLinearRegress[attribute].r2 === 0 &&
                rmSimpleLinearRegress[attribute].pValue === 0) ? (
                <span style={{ fontSize: "calc(1.1rem + 0.7vw)" }}>無資料</span>
              ) : (
                <>
                  <div
                    style={{
                      display: "flex",
                      flexDirection: "column",
                      alignItems: "center",
                      gap: "0.2rem",
                    }}
                  >
                    <div
                      style={{
                        fontSize: "0.75rem",
                        color: "#666",
                        textAlign: "center",
                      }}
                    >
                      R² = {rmSimpleLinearRegress[attribute].r2.toFixed(3)}
                    </div>
                    <div
                      style={{
                        fontSize: "0.75rem",
                        color: "#666",
                        textAlign: "center",
                      }}
                    >
                      p = {rmSimpleLinearRegress[attribute].pValue.toFixed(4)}
                    </div>
                    <span
                      style={{
                        fontSize: "calc(1.1rem + 0.7vw)",
                        color: rankStyle ? rankStyle.textColor : "inherit",
                        fontWeight: rankStyle ? "600" : "normal",
                      }}
                    >
                      {rmSimpleLinearRegress[attribute].pValue <= 0.05 ? (
                        <span
                          style={{
                            color: rankStyle
                              ? rankStyle.textColor
                              : "var(--accent-color)",
                          }}
                        >
                          顯著
                        </span>
                      ) : (
                        <span>不顯著</span>
                      )}
                    </span>
                    {rmSimpleLinearRegress[attribute].pValue <= 0.05 && (
                      <div
                        style={{
                          fontSize: "0.7rem",
                          color: rankStyle ? rankStyle.textColor : "#888",
                          textAlign: "center",
                          marginTop: "0.2rem",
                        }}
                      >
                        影響程度:{" "}
                        {rmSimpleLinearRegress[attribute].slope.toFixed(3)}
                        <br />
                        <span style={{ fontSize: "0.6rem", color: "#aaa" }}>
                          (該屬性分數每提升1分對商家評分的影響)
                        </span>
                      </div>
                    )}
                  </div>
                </>
              )}
            </div>
          );
        })}
      </div>
    </>
  );
};

export default RmSimpleLinearRegress;
