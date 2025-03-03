import axios from "axios";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import PageNotFound from "./404";

const ReviewMiningResult = () => {
  const { dataUUID } = useParams();
  const [resultStatus, setResultStatus] = useState();
  const [result, setResult] = useState();
  const [error, setError] = useState();
  useEffect(() => {
    axios
      .get(`/api/review-mining/${dataUUID}`)
      .then((res) => {
        console.table(res.data);
        setResult(res.data);
        setResultStatus(res.status);
      })
      .catch((err) => {
        setError(err);
      });
  }, []);
  if (error) {
    return (
      <div className="container">
        <h1>Review Mining Result</h1>
        <p>Failed to fetch data</p>
      </div>
    );
  } else if (resultStatus === 404) {
    return <PageNotFound />;
  } else {
    return (
      <div className="container">
        <h1>Review Mining Result</h1>
        <p>UUID: {dataUUID}</p>
      </div>
    );
  }
};

export default ReviewMiningResult;
