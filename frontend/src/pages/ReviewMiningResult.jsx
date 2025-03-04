import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import FullScreenLoader from "../components/FullScreenLoader";
import PageNotFound from "./404";
import "./ReviewMiningResult.css";
import RmAttributesPieChart from "../charts/RmAttributesPieChart";

const ReviewMiningResult = () => {
  const { dataUUID } = useParams();
  const [isLoading, setIsLoading] = useState(true);
  const [resultStatus, setResultStatus] = useState();
  const [result, setResult] = useState();
  const [error, setError] = useState();

  const rmStoreName = result?.storeName;
  const rmAttributes = result?.attributes;
  const rmResults = result?.results;
  const rmSummary = result?.summary;
  const rmTimestamp = result?.timestamp;

  useEffect(() => {
    axios
      .get(`/api/review-mining/${dataUUID}`)
      .then((res) => {
        console.table(res.data);
        setResult(res.data);
        setResultStatus(res.status);
        setIsLoading(false);
      })
      .catch((err) => {
        setError(err);
        alert("Failed to fetch data");
        setIsLoading(false);
      });
  }, []);
  if (isLoading) {
    return <FullScreenLoader />;
  } else if (resultStatus === 404) {
    return <PageNotFound />;
  } else if (result) {
    return (
      <div className="container">
        <h1>Review Mining Result</h1>
        <p>UUID: {dataUUID}</p>
        <h2>商店名稱</h2>
        <span>{result.storeName}</span>
        <div className="chart-box">
          <RmAttributesPieChart
            rmAttributes={rmAttributes}
            rmResults={rmResults}
          />
        </div>
      </div>
    );
  }
};

export default ReviewMiningResult;
