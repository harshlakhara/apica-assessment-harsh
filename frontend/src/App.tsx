import { useEffect, useState } from "react";
import "./App.css";

const METHODS: { [k: string]: string } = {
  put: "PUT",
  get: "GET",
  delete: "DELETE",
};

function App() {
  const [method, setMethod] = useState("put");
  const [reqBody, setReqBody] = useState({ key: "", value: "", ttl: "" });
  const [cache, setCache] = useState<{ [k: string]: any }[]>([]);
  const [dropdownOpen, setOpendrown] = useState(false);
  const [getValue, setGetValue] = useState("");

  useEffect(() => {
    const ws = new WebSocket("http://localhost:3000/ws/cachefeed");
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data);
      setCache(data);
    };

    getSnapshot();

    return () => ws.close();
  }, []);

  const handleFormChange = (target: HTMLInputElement) => {
    setReqBody({ ...reqBody, [target.name]: target.value });
  };

  const handlePut = async () => {
    if (reqBody.key === "" || reqBody.ttl === "" || reqBody.value === "")
      return;
    try {
      const res = await fetch("http://localhost:3000", {
        method: "POST",
        body: JSON.stringify({ ...reqBody, ttl: parseInt(reqBody.ttl) }),
        headers: {
          "Content-Type": "application/json",
        },
      });
      const json = await res.json();
      if (json.ok) {
        setReqBody({
          key: "",
          value: "",
          ttl: "",
        });
      }
    } catch (error) {
      console.log(error);
    }
  };

  const handleGet = async () => {
    try {
      const key = reqBody.key;
      const res = await fetch(`http://localhost:3000/${key}`);
      const json = await res.json();
      if (json.ok) {
        setGetValue(key);
        setTimeout(() => {
          setGetValue("");
        }, 1000);
        setReqBody({
          key: "",
          value: "",
          ttl: "",
        });
      }
    } catch (error) {
      console.error(error);
    }
  };

  const handleDelete = async () => {
    try {
      const key = reqBody.key;
      const res = await fetch(`http://localhost:3000/${key}`, {
        method: "DELETE",
      });
      const json = await res.json();
      if (json.ok) {
        setReqBody({
          key: "",
          value: "",
          ttl: "",
        });
      }
    } catch (error) {
      console.error(error);
    }
  };

  document.addEventListener("mousedown", (e) => {
    if (
      dropdownOpen &&
      !(e.target as HTMLDivElement).classList.contains("dropdown-options")
    ) {
      setOpendrown(false);
    }
  });

  const getSnapshot = async () => {
    try {
      const res = await fetch(`http://localhost:3000/snapshot`);
      const json = await res.json();
      setCache(json);
    } catch (error) {
      console.error(error);
    }
  };

  const clearCache = async () => {
    try {
      const res = await fetch(`http://localhost:3000/clear`, {
        method: "DELETE",
      });
      const json = await res.json();
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <>
      <div className="layout-wrapper">
        <div className="toolbar">
          <h2>Toolbar</h2>
          <div className="input-section">
            <div className={`dropdown-wrapper ${dropdownOpen ? "show" : ""}`}>
              <div
                className="dropdown-toggle"
                onClick={() => {
                  setOpendrown(true);
                }}
              >
                <span>{METHODS[method]} </span>
                <img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAADIAAAAyCAYAAAAeP4ixAAAACXBIWXMAAAsTAAALEwEAmpwYAAAAvklEQVR4nO3RwQqDMBCE4aGHYt92j75XD7ZPZxEUSsCmarKZ1flhbznMRwCllFJKKeXZDSfYeAfwBNCDNwPwAvDIIcb5elLEON97DZNC2DCWbJsg3dpjVoxtQbBibA+CDWNHECwYK4FojbGSiFYYq4HwxlhNhBfGBVEb44qohWmCKI1piiiFoUAcxVAh9mIoEVsx1Ih/MSEQOUwoxNI0cEiGf98QAfHrZ8L8RA4TEpFiQiOWujMglFJKKaWu3AfcdLEudXvTswAAAABJRU5ErkJggg=="></img>
              </div>
              <div className="dropdown-body">
                {Object.entries(METHODS).map(([k, v]) => (
                  <div
                    key={k}
                    onClick={() => {
                      setMethod(k);
                      setOpendrown(false);
                    }}
                    className={`dropdown-options ${
                      method == k ? "active" : ""
                    }`}
                  >
                    {v}
                  </div>
                ))}
              </div>
            </div>
            {method === "put" ? (
              <div className="form put-form">
                <input
                  type="text"
                  name="key"
                  value={reqBody.key}
                  onChange={(e) => handleFormChange(e.target)}
                  placeholder="Key"
                />
                <input
                  type="text"
                  value={reqBody.value}
                  name="value"
                  onChange={(e) => handleFormChange(e.target)}
                  placeholder="Value"
                />
                <input
                  type="number"
                  value={reqBody.ttl}
                  name="ttl"
                  min={0}
                  onChange={(e) => handleFormChange(e.target)}
                  placeholder="Time to live in seconds"
                />
                <button onClick={() => handlePut()}>Put</button>
              </div>
            ) : method === "get" ? (
              <div className="form get-form">
                <input
                  type="text"
                  value={reqBody.key}
                  name="key"
                  onChange={(e) => handleFormChange(e.target)}
                  placeholder="Key"
                />
                <button onClick={() => handleGet()}>Get</button>
              </div>
            ) : (
              <div className="form delete-form">
                <input
                  type="text"
                  value={reqBody.key}
                  name="key"
                  onChange={(e) => handleFormChange(e.target)}
                  placeholder="Key"
                />
                <button onClick={() => handleDelete()}>Delete</button>
              </div>
            )}
          </div>
        </div>
        <div className="insights">
          <h2>Current State | With Capacity 5</h2>
          <div>
            <button onClick={clearCache}>Clear Cache</button>
          </div>
          <div className="cache-visualisation">
            {cache.map((ele: any) => (
              <>
                <div
                  className={`cache-block ${
                    getValue === ele.key ? "highlight" : ""
                  }`}
                >
                  <div className="key">Key: {ele.key}</div>
                  <div className="value">Value: {ele.value}</div>
                  <div className="ttl">TTL: {ele.ttl}</div>
                </div>
              </>
            ))}
          </div>
        </div>
      </div>
    </>
  );
}

export default App;
