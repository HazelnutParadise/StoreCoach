import { Doughnut } from "react-chartjs-2";
import ChartDataLabels from "chartjs-plugin-datalabels";
import { Chart, Legend } from "chart.js";
import ChartStyle from "./ChartStyle";

Chart.register(ChartDataLabels);

const RmAttributesPieChart = ({
  rmAttributes,
  rmResults,
  topFiveAttributesByMentions,
  attributeMentionCounts,
}) => {
  // 如果有預計算的數據，使用它們；否則回退到原始計算
  let sortedAttributes, sortedAttributeTimesList;

  if (topFiveAttributesByMentions && attributeMentionCounts) {
    // 使用預計算的數據
    sortedAttributes = topFiveAttributesByMentions.map(
      (item) => item.attribute
    );
    sortedAttributeTimesList = topFiveAttributesByMentions.map(
      (item) => item.count
    );

    // 計算"其他"類別
    const totalMentions = Object.values(attributeMentionCounts).reduce(
      (a, b) => a + b,
      0
    );
    const topFiveMentions = sortedAttributeTimesList.reduce((a, b) => a + b, 0);
    const otherMentions = totalMentions - topFiveMentions;

    if (otherMentions > 0) {
      sortedAttributes.push("其他");
      sortedAttributeTimesList.push(otherMentions);
    }
  } else {
    // 回退到原始計算邏輯
    const attributeTimes = {};
    rmAttributes.forEach((element) => {
      attributeTimes[element] = 0;
    });

    rmResults.forEach((result) => {
      const miningResults = result.miningResults;
      miningResults.forEach((miningResult) => {
        if (miningResult.attribute in attributeTimes) {
          attributeTimes[miningResult.attribute]++;
        }
      });
    });

    let sortedAttributesTemp = Object.keys(attributeTimes).sort(
      (a, b) => attributeTimes[b] - attributeTimes[a]
    );

    sortedAttributeTimesList = [];
    let count = 0;
    sortedAttributesTemp.forEach((attribute) => {
      if (count < 5) {
        sortedAttributeTimesList.push(attributeTimes[attribute]);
        count++;
      } else if (count === 5) {
        sortedAttributeTimesList[5] = attributeTimes[attribute];
        count++;
      } else {
        console.log(attribute);
        sortedAttributeTimesList[5] += attributeTimes[attribute];
      }
    });

    sortedAttributes = sortedAttributesTemp.slice(0, 5);
    if (sortedAttributeTimesList[5] && sortedAttributeTimesList[5] !== 0)
      sortedAttributes.push("其他");
  }

  console.log(sortedAttributes);
  console.log(sortedAttributeTimesList);

  const data = {
    labels: sortedAttributes,
    datasets: [
      {
        label: "提及次數",
        data: sortedAttributeTimesList,
      },
    ],
  };

  const options = {
    maintainAspectRatio: false,
    plugins: {
      datalabels: {
        formatter: (value, context) => {
          const total = context.chart.data.datasets[0].data.reduce(
            (a, b) => a + b,
            0
          );
          const percentageNumber = (value / total) * 100;
          if (percentageNumber < 6.25) {
            return "";
          }
          const percentage = percentageNumber.toFixed(2) + "%";
          return percentage;
        },
        color: "#000",
        // backgroundColor: "#fff",
        // anchor: "end",
        // align: "end",
      },
    },
  };

  return (
    <div style={ChartStyle}>
      <h3>顧客在意的屬性</h3>
      <Doughnut data={data} options={options} width={"100%"} />
    </div>
  );
};

export default RmAttributesPieChart;
