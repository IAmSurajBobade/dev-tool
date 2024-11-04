import React from "react";

function SettingsModal({ visibleTabs, updateVisibleTabs, closeModal }) {
  const handleCheckboxChange = (tabName) => {
    updateVisibleTabs({
      ...visibleTabs,
      [tabName]: !visibleTabs[tabName],
    });
  };

  return (
    <div className="modal-overlay">
      <div className="modal">
        <h2>Settings</h2>
        <div className="modal-content">
          <h3>Visible Tabs</h3>
          <label>
            <input
              type="checkbox"
              checked={visibleTabs.base64}
              onChange={() => handleCheckboxChange("base64")}
            />
            Show Base64 Tab
          </label>
          <label>
            <input
              type="checkbox"
              checked={visibleTabs.jwt}
              onChange={() => handleCheckboxChange("jwt")}
            />
            Show JWT Tab
          </label>
        </div>
        <button onClick={closeModal}>Close</button>
      </div>
    </div>
  );
}

export default SettingsModal;
