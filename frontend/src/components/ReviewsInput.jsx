import React, { useState, forwardRef, useImperativeHandle } from "react";
import Spreadsheet from "react-spreadsheet";
import "./ReviewsInput.css";

const SpreadsheetInput = forwardRef((props, ref) => {
  const { columnLabels } = props;
  const [data, setData] = useState(Array(10).fill([{ value: "" }])); // 預設 10 列 1 欄

  // 取得表格數據並忽略空行
  const getData = () => {
    return data
      .map((row) => (row[0] && row[0].value ? row[0].value.trim() : ""))
      .filter((value) => value !== "");
  };

  // 讓父元件可以呼叫 getData()
  useImperativeHandle(ref, () => ({
    getData,
  }));

  // 監聽數據變更
  const handleChange = (newData) => {
    setData(newData);

    // 自動新增行（最後一行輸入時）
    if (newData.length > 0 && newData[newData.length - 1][0].value !== "") {
      setData([...newData, [{ value: "" }]]);
    }
  };

  return (
    <div className="reviews-input">
      <Spreadsheet
        data={data}
        onChange={handleChange}
        hideRowIndicators={true}
        columnLabels={columnLabels}
      />
    </div>
  );
});

export default SpreadsheetInput;
