import React, { useRef, useEffect } from "react";
import WordCloud from "wordcloud";

const RealWordCloud = ({ title, attributes, color }) => {
  const canvasRef = useRef(null);

  useEffect(() => {
    if (
      canvasRef.current &&
      attributes &&
      Array.isArray(attributes) &&
      attributes.length > 0
    ) {
      const canvas = canvasRef.current;
      const ctx = canvas.getContext("2d");

      // 高 DPI 支持，提升文字清晰度
      const dpr = window.devicePixelRatio || 1;
      const rect = canvas.getBoundingClientRect();

      // 設置實際尺寸（考慮設備像素比）
      canvas.width = rect.width * dpr;
      canvas.height = rect.height * dpr;

      // 縮放上下文以匹配設備像素比
      ctx.scale(dpr, dpr);

      // 設置畫布顯示尺寸
      canvas.style.width = rect.width + "px";
      canvas.style.height = rect.height + "px";

      // 清除畫布
      ctx.clearRect(0, 0, canvas.width, canvas.height); // 準備文字雲數據
      const wordList = attributes.map((attr) => [attr.attribute, attr.count]); // 計算當前數據集的最大和最小值，用於動態調整大小
      const counts = attributes.map((attr) => attr.count);
      const maxCount = Math.max(...counts);
      const minCount = Math.min(...counts);
      const countRange = maxCount - minCount || 1; // 避免除以0

      // 文字雲選項
      const options = {
        list: wordList,
        gridSize: Math.max(4, Math.round((6 * canvas.width) / 2048)), // 針對高解析度調整網格
        weightFactor: function (size) {
          // 基於當前數據集的範圍動態調整大小，針對高解析度優化
          const normalizedSize = (size - minCount) / countRange;
          const baseSize = Math.max(32, (canvas.width / 2048) * 80); // 針對高解析度提升基礎大小
          const sizeMultiplier = 1.2 + normalizedSize * 2.5; // 1.2倍到3.7倍之間
          return baseSize * sizeMultiplier;
        },
        fontFamily: 'Arial, "Microsoft YaHei", "PingFang SC", sans-serif',
        color: function (word, weight) {
          // 根據權重決定顏色深淺
          const opacity = Math.min(1, Math.max(0.6, weight / 10));
          if (color === "#2ecc71") {
            return `rgba(46, 204, 113, ${opacity})`;
          } else if (color === "#e74c3c") {
            return `rgba(231, 76, 60, ${opacity})`;
          }
          return color;
        },
        rotateRatio: 0.3,
        rotationSteps: 2,
        backgroundColor: "transparent",
        click: function (item) {
          console.log(`點擊了屬性: ${item[0]}, 提及次數: ${item[1]}`);
        },
        hover: function (item) {
          if (item) {
            canvas.style.cursor = "pointer";
            canvas.title = `${item[0]}: ${item[1]}次`;
          } else {
            canvas.style.cursor = "default";
            canvas.title = "";
          }
        },
      };

      // 生成文字雲
      try {
        WordCloud(canvas, options);
      } catch (error) {
        console.error("WordCloud 生成錯誤:", error);
        // 如果出錯，在畫布上顯示錯誤信息
        ctx.fillStyle = "#666";
        ctx.font = "16px Arial";
        ctx.textAlign = "center";
        ctx.fillText("文字雲載入失敗", canvas.width / 2, canvas.height / 2);
      }
    }
  }, [attributes, color]);

  // 安全檢查：確保 attributes 是數組且不為空
  if (!attributes || !Array.isArray(attributes) || attributes.length === 0) {
    return (
      <div className="wordcloud-container">
        <h3 style={{ textAlign: "center", marginBottom: "1rem", color: color }}>
          {title}
        </h3>
        <div
          style={{
            height: "300px",
            width: "100%",
            display: "flex",
            alignItems: "center",
            justifyContent: "center",
            backgroundColor: "#f8f9fa",
            borderRadius: "8px",
            border: "1px solid #e9ecef",
          }}
        >
          <p style={{ color: "#666", fontSize: "1rem" }}>暫無數據</p>
        </div>
      </div>
    );
  }

  return (
    <div className="wordcloud-container">
      <h3 style={{ textAlign: "center", marginBottom: "1rem", color: color }}>
        {title}
      </h3>
      <div
        style={{
          height: "300px",
          width: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          backgroundColor: "#ffffff",
          borderRadius: "8px",
          border: "1px solid #e9ecef",
          boxShadow: "0 2px 4px rgba(0,0,0,0.1)",
        }}
      >
        {" "}
        <canvas
          ref={canvasRef}
          style={{
            width: "100%",
            height: "100%",
            borderRadius: "6px",
          }}
        />
      </div>
    </div>
  );
};

export default RealWordCloud;
