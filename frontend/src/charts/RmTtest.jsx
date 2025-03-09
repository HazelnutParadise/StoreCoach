const RmTtest = ({ rmResults, rmTtest }) => {
  let attributes = Object.keys(rmTtest);
  let attributeCount = {};
  attributes.forEach((element) => {
    attributeCount[element] = 0;
  });
  // 計算每個屬性出現的次數
  rmResults.forEach((result) => {
    const miningResults = result.miningResults;
    miningResults.forEach((miningResult) => {
      if (miningResult.attribute in attributeCount) {
        attributeCount[miningResult.attribute]++;
      }
    });
  });
  attributes.sort((a, b) => attributeCount[b] - attributeCount[a]).slice(0, 5);
  return (
    <>
      <h3>前五大屬性對顧客評分影響的顯著性（T檢定p值，顯著水準=0.05）</h3>
      <div className="ttest-container">
        {attributes.map((attribute) => (
          <div className="ttest-box">
            <h4
              style={{
                padding: 0,
                margin: 0,
                fontSize: "calc(0.9rem + 0.7vw)",
              }}
            >
              {attribute}
            </h4>
            {(rmTtest[attribute].pValue < 0.05 &&
              rmTtest[attribute].tValue < 0.05 &&
              rmTtest[attribute].df === 0) ||
            !rmTtest[attribute] ? (
              <span style={{ fontSize: "calc(1.1rem + 0.7vw)" }}>無資料</span>
            ) : (
              <>
                <span style={{ fontSize: "0.85rem" }}>
                  {rmTtest[attribute].pValue.toFixed(4)}
                </span>
                <span style={{ fontSize: "calc(1.1rem + 0.7vw)" }}>
                  {rmTtest[attribute].pValue <= 0.05 ? (
                    <span style={{ color: " var(--accent-color)" }}>顯著</span>
                  ) : (
                    <span>不顯著</span>
                  )}
                </span>
              </>
            )}
          </div>
        ))}
      </div>
    </>
  );
};

export default RmTtest;
