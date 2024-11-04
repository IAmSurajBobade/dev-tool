import { useState, useEffect } from "react";
import Base64Tool from "./components/Base64Tool";
import JWTTool from "./components/JWTTool";
import Notification from "./components/Notification";
import SettingsModal from "./components/SettingsModal";
import AboutModal from "./components/AboutModal";
import CreditsModal from "./components/CreditsModal";
import "./css/App.css";
import "./css/Notification.css";
import "./css/Model.css";

function App() {
  const [activeTab, setActiveTab] = useState("base64");
  const [notification, setNotification] = useState(null);
  const [modalOpen, setModalOpen] = useState(null);
  const [visibleTabs, setVisibleTabs] = useState({
    base64: true,
    jwt: true,
  });
  const [appInfo, setAppInfo] = useState({
    version: "1.0.0",
    buildDate: "2023-06-15",
    commitHash: "abc123",
  });
  const [menuOpen, setMenuOpen] = useState(false);

  const showNotification = (message) => {
    setNotification(null);
    setTimeout(() => setNotification(message), 10);
  };

  useEffect(() => {
    if (notification) {
      const timer = setTimeout(() => {
        setNotification(null);
      }, 3000);

      return () => clearTimeout(timer);
    }
  }, [notification]);

  const openExternalLink = (url) => (e) => {
    e.preventDefault();
    if (window.wails && window.wails.BrowserOpenURL) {
      window.wails.BrowserOpenURL(url);
    } else {
      window.open(url, "_blank", "noopener,noreferrer");
    }
  };

  const toggleMenu = () => {
    setMenuOpen(!menuOpen);
  };

  const openModal = (modalType) => {
    setModalOpen(modalType);
    setMenuOpen(false);
  };

  const closeModal = () => {
    setModalOpen(null);
  };

  const updateVisibleTabs = (newVisibleTabs) => {
    setVisibleTabs(newVisibleTabs);
    if (!newVisibleTabs[activeTab]) {
      const firstVisibleTab = Object.keys(newVisibleTabs).find(
        (tab) => newVisibleTabs[tab]
      );
      if (firstVisibleTab) {
        setActiveTab(firstVisibleTab);
      }
    }
  };

  return (
    <div id="App">
      <div className="content-wrapper">
        <div className="header">
          <div className="tab-container">
            {visibleTabs.base64 && (
              <button
                className={`tab-button ${
                  activeTab === "base64" ? "active" : ""
                }`}
                onClick={() => setActiveTab("base64")}
              >
                Base64
              </button>
            )}
            {visibleTabs.jwt && (
              <button
                className={`tab-button ${activeTab === "jwt" ? "active" : ""}`}
                onClick={() => setActiveTab("jwt")}
              >
                JWT
              </button>
            )}
          </div>
        </div>
        <div className="content">
          {activeTab === "base64" && visibleTabs.base64 && (
            <Base64Tool showNotification={showNotification} />
          )}
          {activeTab === "jwt" && visibleTabs.jwt && (
            <JWTTool showNotification={showNotification} />
          )}
        </div>
      </div>
      {notification && <Notification message={notification} />}
      <div className="menu-button-container">
        <button className="menu-button" onClick={toggleMenu}>
          ☰
        </button>
        {menuOpen && (
          <div className="menu-dropdown">
            <button onClick={() => openModal("about")}>About</button>
            <button onClick={() => openModal("settings")}>Settings</button>
            <button onClick={() => openModal("credits")}>Credits</button>
          </div>
        )}
      </div>
      {modalOpen === "about" && (
        <AboutModal appInfo={appInfo} closeModal={closeModal} />
      )}
      {modalOpen === "settings" && (
        <SettingsModal
          visibleTabs={visibleTabs}
          updateVisibleTabs={updateVisibleTabs}
          closeModal={closeModal}
        />
      )}
      {modalOpen === "credits" && (
        <CreditsModal
          closeModal={closeModal}
          openExternalLink={openExternalLink}
        />
      )}
    </div>
  );
}

export default App;
