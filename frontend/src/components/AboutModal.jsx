import React from "react";

function AboutModal({ appInfo, closeModal }) {
  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>About</h2>
        <div className="modal-content">
          <p>
            <strong>Version:</strong> {appInfo.version}
          </p>
          <p>
            <strong>Build Date:</strong> {appInfo.buildDate}
          </p>
          <p>
            <strong>Commit Hash:</strong> {appInfo.commitHash}
          </p>
        </div>
        <button onClick={closeModal}>Close</button>
      </div>
    </div>
  );
}

export default AboutModal;
