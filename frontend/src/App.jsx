import { useState, useEffect } from "react";
import Base64Tool from "./components/Base64Tool";
import JWTTool from "./components/JWTTool";
import Notification from "./components/Notification";
import "./App.css";
import "./Notification.css";

function App() {
  const [activeTab, setActiveTab] = useState("base64");
  const [notification, setNotification] = useState(null);

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
    BrowserOpenURL(url);
  };

  return (
    <div id="App">
      <div className="content-wrapper">
        <div className="tab-container">
          <button
            className={`tab-button ${activeTab === "base64" ? "active" : ""}`}
            onClick={() => setActiveTab("base64")}
          >
            Base64
          </button>
          <button
            className={`tab-button ${activeTab === "jwt" ? "active" : ""}`}
            onClick={() => setActiveTab("jwt")}
          >
            JWT
          </button>
        </div>
        <div className="content">
          {activeTab === "base64" && (
            <Base64Tool showNotification={showNotification} />
          )}
          {activeTab === "jwt" && (
            <JWTTool showNotification={showNotification} />
          )}
        </div>
      </div>
      {notification && <Notification message={notification} />}
      <footer className="footer">Built with ❤️ using Wails and React</footer>
    </div>
  );
}

export default App;
