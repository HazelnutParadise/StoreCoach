import { useEffect } from "react";
import { v4 as uuidv4 } from "uuid";

const TextAD = () => {
  const adUUID = uuidv4();
  
  useEffect(() => {
    const custom_call = function (params) {
      if (params.hasAd) {
        console.log("TD has AD");
      } else {
        console.log("TD AD Empty");
      }
    };

    const ONEAD_TEXT = {
      pub: {
        uid: "2000181",
        slotobj: document.getElementById(`div-onead-draft-${adUUID}`),
        player_mode: "text-drive",
        queryAdCallback: custom_call
      }
    };

    window.ONEAD_text_pubs = window.ONEAD_text_pubs || [];
    window.ONEAD_text_pubs.push(ONEAD_TEXT);
  }, [adUUID]);

  return <div id={`div-onead-draft-${adUUID}`}></div>;
};

export default TextAD;