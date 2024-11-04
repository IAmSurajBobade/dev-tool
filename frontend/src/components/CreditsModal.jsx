import React from "react";

function CreditsModal({ closeModal, openExternalLink }) {
  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>Credits</h2>
        <div className="modal-content">
          <p>This application was built using:</p>
          <ul>
            <li>
              <a
                href="https://wails.io"
                onClick={openExternalLink("https://wails.io")}
              >
                Wails
              </a>
            </li>
            <li>
              <a
                href="https://reactjs.org"
                onClick={openExternalLink("https://reactjs.org")}
              >
                React
              </a>
            </li>
          </ul>
        </div>
        <button onClick={closeModal}>Close</button>
      </div>
    </div>
  );
}

export default CreditsModal;
