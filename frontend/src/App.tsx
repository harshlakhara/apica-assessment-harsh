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

  useEffect(() => {
    const ws = new WebSocket("http://localhost:3000/ws/cachefeed");
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data);
      setCache(data);
    };

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

  const getSnapshot = async () => {
    try {
      const res = await fetch(`http://localhost:3000/snapshot`);
      const json = await res.json();
      setCache(json);
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
            <div className="dropdown-wrapper">
              <div className="dropdown-toggle">{METHODS[method]}</div>
              <div className="dropdown-body">
                {Object.entries(METHODS).map(([k, v]) => (
                  <div
                    key={k}
                    onClick={() => setMethod(k)}
                    className="dropdown-options"
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
                  onChange={(e) => handleFormChange(e.target)}
                  placeholder="Time to live"
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
          <h2>Insights</h2>
          <div className="cache-visualisation">
            {cache.map((ele: any) => (
              <>
                <div className="cache-block">
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
