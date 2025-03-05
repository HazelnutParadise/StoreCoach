import { Doughnut } from "react-chartjs-2";
import ChartDataLabels from "chartjs-plugin-datalabels";
import { Chart } from "chart.js";
import ChartStyle from "./ChartStyle";

Chart.register(ChartDataLabels);

const RmAttributesPieChart = ({ rmAttributes, rmResults }) => {
  const results = rmResults;
  const attributeTimes = {};
  rmAttributes.forEach((element) => {
    attributeTimes[element] = 0;
  });

  results.forEach((result) => {
    const miningResults = result.miningResults;
    miningResults.forEach((miningResult) => {
      if (miningResult.attribute in attributeTimes) {
        attributeTimes[miningResult.attribute]++;
      }
    });
  });

  let sortedAttributes = Object.keys(attributeTimes).sort(
    (a, b) => attributeTimes[b] - attributeTimes[a]
  );

  const sortedAttributeTimesList = [];
  let count = 0;
  sortedAttributes.forEach((attribute) => {
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

  sortedAttributes = sortedAttributes.slice(0, 5);
  sortedAttributes.push("其他");

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
