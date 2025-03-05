import { Bar } from "react-chartjs-2";
import ChartStyle from "./ChartStyle";

const RmAttributesBarChart = ({}) => {
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
        position: "bottom",
      },
    },
  };

  const labels = ["January", "February", "March", "April", "May"];

  const data = {
    labels,
    datasets: [
      {
        label: "正面",
        data: labels.map(() => 5),
        borderColor: "rgb(75, 192, 112)",
        backgroundColor: "rgba(75, 192, 112,0.8)",
      },
      {
        label: "負面",
        data: labels.map(() => 10),
        borderColor: "rgb(255, 99, 132)",
        backgroundColor: "rgba(255, 99, 132, 0.8)",
      },
      {
        label: "中立",
        data: labels.map(() => 5),
        borderColor: "rgb(53, 162, 235)",
        backgroundColor: "rgba(53, 162, 235, 0.8)",
      },
    ],
  };

  return (
    <div style={ChartStyle}>
      <h3>前五大屬性在各面向被提及的次數</h3>
      <Bar options={options} data={data} />
    </div>
  );
};

export default RmAttributesBarChart;
