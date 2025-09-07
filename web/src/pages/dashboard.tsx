import { useEffect, useState } from "react";

type Car = { id: string; state: "FREE"|"RESERVED"|"BUSY"; assignedTo?: string; updatedAt: string; ttl?: number };
type Ev = { type: "snapshot"|"update"; data: any };

export default function Dashboard() {
  const [cars, setCars] = useState<Car[]>([]);
  const [status, setStatus] = useState("connecting…");

  useEffect(() => {
    const es = new EventSource("http://localhost:8081/events");
    es.onmessage = (e) => {
      const ev: Ev = JSON.parse(e.data);
      if (ev.type === "snapshot") setCars(ev.data);
      if (ev.type === "update") {
        fetch("http://localhost:8081/rooms")
          .then(r => r.json())
          .then(setCars)
          .catch(()=>{});
      }
      setStatus("live");
    };
    es.onerror = () => setStatus("disconnected");
    return () => es.close();
  }, []);

  return (
    <main style={{ padding: 20, fontFamily: "system-ui" }}>
      <h1>Live Cars <small style={{fontSize:12, color:"#888"}}>({status})</small></h1>
      <table cellPadding={8} style={{ borderCollapse:"collapse", border: "1px solid #ddd", minWidth:600 }}>
        <thead><tr><th>ID</th><th>State</th><th>Assigned</th><th>TTL</th><th>Updated</th></tr></thead>
        <tbody>
        {cars.map(c=>(
          <tr key={c.id}>
            <td>{c.id}</td>
            <td>{c.state}</td>
            <td>{c.assignedTo||"-"}</td>
            <td>{c.state==="RESERVED"? (c.ttl??0) : "-"}</td>
            <td>{new Date(c.updatedAt).toLocaleString()}</td>
          </tr>
        ))}
        </tbody>
      </table>
    </main>
  );
}
