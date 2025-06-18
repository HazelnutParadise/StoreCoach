const RmSimpleLinearRegress = ({ rmAttributes, rmSimpleLinearRegress }) => {
  let attributes = rmAttributes;
  attributes = attributes.slice(0, 5);

  return (
    <>
      <h3>
        提及次數前五多的屬性得分分別對商家評分的影響（簡單線性回歸分析R²值與顯著性，顯著水準=0.05）
      </h3>
      <div className="ttest-container">
        {attributes.map((attribute) => (
          <div className="ttest-box" key={attribute}>
            <h4
              style={{
                padding: 0,
                margin: 0,
                fontSize: "calc(0.9rem + 0.7vw)",
              }}
            >
              {attribute}
            </h4>
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
                  <span style={{ fontSize: "calc(1.1rem + 0.7vw)" }}>
                    {rmSimpleLinearRegress[attribute].pValue <= 0.05 ? (
                      <span style={{ color: "var(--accent-color)" }}>顯著</span>
                    ) : (
                      <span>不顯著</span>
                    )}
                  </span>{" "}
                  {rmSimpleLinearRegress[attribute].pValue <= 0.05 && (
                    <div
                      style={{
                        fontSize: "0.7rem",
                        color: "#888",
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
        ))}
      </div>
    </>
  );
};

export default RmSimpleLinearRegress;
