import { Doughnut } from "react-chartjs-2";

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

  const sortedAttributes = Object.keys(attributeTimes).sort(
    (a, b) => attributeTimes[b] - attributeTimes[a]
  );

  const attributeTimesArray = Object.values(attributeTimes).sort(
    (a, b) => b - a
  );

  const data = {
    labels: sortedAttributes,
    datasets: [
      {
        label: "提及次數",
        data: attributeTimesArray,
      },
    ],
  };
  return <Doughnut data={data} />;
};

export default RmAttributesPieChart;
