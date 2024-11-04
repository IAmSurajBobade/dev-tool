import { useState } from "react";
import { EncodeBase64, DecodeBase64 } from "../../wailsjs/go/main/App";

function Base64Tool({ showNotification }) {
  const [input, setInput] = useState("");
  const [output, setOutput] = useState("");

  const handleEncode = async () => {
    try {
      const result = await EncodeBase64(input);
      setOutput(result);
      showNotification("Text encoded successfully");
    } catch (error) {
      console.error("Error encoding:", error);
      showNotification("Error encoding text");
    }
  };

  const handleDecode = async () => {
    try {
      const result = await DecodeBase64(input);
      setOutput(result);
      showNotification("Text decoded successfully");
    } catch (error) {
      console.error("Error decoding:", error);
      showNotification("Error decoding text");
    }
  };

  return (
    <div className="tool-container">
      <h2>Base64 Encoder/Decoder</h2>
      <textarea
        value={input}
        onChange={(e) => setInput(e.target.value)}
        placeholder="Enter text to encode/decode"
      />
      <div className="button-group">
        <button onClick={handleEncode}>Encode</button>
        <button onClick={handleDecode}>Decode</button>
      </div>
      <textarea value={output} readOnly placeholder="Result" />
    </div>
  );
}

export default Base64Tool;
