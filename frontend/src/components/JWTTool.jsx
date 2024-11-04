import { useState } from "react";
import { EncodeJWT, DecodeJWT } from "../../wailsjs/go/main/App";

function JWTTool({ showNotification }) {
  const [input, setInput] = useState("");
  const [secret, setSecret] = useState("");
  const [output, setOutput] = useState("");

  const handleEncode = async () => {
    try {
      const result = await EncodeJWT(input, secret);
      setOutput(result);
      showNotification("JWT encoded successfully");
    } catch (error) {
      console.error("Error encoding JWT:", error);
      showNotification("Error encoding JWT");
    }
  };

  const handleDecode = async () => {
    try {
      const result = await DecodeJWT(input, secret);
      setOutput(JSON.stringify(result, null, 2));
      showNotification("JWT decoded successfully");
    } catch (error) {
      console.error("Error decoding JWT:", error);
      showNotification("Error decoding JWT");
    }
  };

  return (
    <div className="tool-container">
      <h2>JWT Encoder/Decoder</h2>
      <textarea
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Enter payload to encode or JWT to decode"
      />
      <input
        type="text"
        value={secret}
        onChange={(e) => setSecret(e.target.value)}
        placeholder="Secret key"
      />
      <div className="button-group">
        <button onClick={handleEncode}>Encode</button>
        <button onClick={handleDecode}>Decode</button>
      </div>
      <textarea value={output} readOnly placeholder="Result" />
    </div>
  );
}

export default JWTTool;
