import { Bar } from "react-chartjs-2";
import ChartStyle from "./ChartStyle";

const RmAttributesBarChart = ({ rmAttributes, rmResults }) => {
  const attributesCount = {};
  rmResults.forEach((result) => {
    result.miningResults.forEach((miningResult) => {
      if (miningResult.attribute in attributesCount) {
        attributesCount[miningResult.attribute]++;
      } else {
        attributesCount[miningResult.attribute] = 1;
      }
    });
  });
  const sortedAttributes = [...rmAttributes].sort(
    (a, b) => attributesCount[b] - attributesCount[a]
  );
  const top5Attributes = sortedAttributes.slice(0, 5);
  const positiveCount = {};
  const negativeCount = {};
  const neutralCount = {};
  rmResults.forEach((result) => {
    result.miningResults.forEach((miningResult) => {
      if (top5Attributes.includes(miningResult.attribute)) {
        switch (miningResult.sentiment) {
          case "positive":
            if (miningResult.attribute in positiveCount) {
              positiveCount[miningResult.attribute]++;
            } else {
              positiveCount[miningResult.attribute] = 1;
            }
            break;
          case "negative":
            if (miningResult.attribute in negativeCount) {
              negativeCount[miningResult.attribute]++;
            } else {
              negativeCount[miningResult.attribute] = 1;
            }
            break;
          case "neutral":
            if (miningResult.attribute in neutralCount) {
              neutralCount[miningResult.attribute]++;
            } else {
              neutralCount[miningResult.attribute] = 1;
            }
            break;
          default:
            break;
        }
      }
    });
  });

  const options = {
    indexAxis: "y",
    elements: {
      bar: {
        borderWidth: 2,
      },
    },
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: "top",
      },
    },
    scales: {
      x: { title: { display: true, text: "次數" } },
      y: { title: { display: true, text: "屬性" } },
    },
  };

  const labels = top5Attributes;

  const data = {
    labels,
    datasets: [
      {
        label: "正面",
        data: labels.map((attribute) => positiveCount[attribute] || 0),
        borderColor: "rgb(75, 192, 112)",
        backgroundColor: "rgba(75, 192, 112,0.8)",
      },
      {
        label: "負面",
        data: labels.map((attribute) => negativeCount[attribute] || 0),
        borderColor: "rgb(255, 99, 132)",
        backgroundColor: "rgba(255, 99, 132, 0.8)",
      },
      {
        label: "中立",
        data: labels.map((attribute) => neutralCount[attribute] || 0),
        borderColor: "rgb(53, 162, 235)",
        backgroundColor: "rgba(53, 162, 235, 0.8)",
      },
    ],
  };

  return (
    <div style={ChartStyle}>
      <h3>前五大屬性被提及時的情緒分布</h3>
      <Bar options={options} data={data} />
    </div>
  );
};

export default RmAttributesBarChart;
