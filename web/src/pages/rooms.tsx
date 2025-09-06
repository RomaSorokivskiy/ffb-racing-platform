import { useEffect, useState } from "react";

type Car = { id: string; state: "FREE" | "RESERVED" | "BUSY"; assignedTo?: string; updatedAt: string };

export default function Rooms() {
  const [cars, setCars] = useState<Car[]>([]);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    (async () => {
      try {
        const res = await fetch("http://localhost:8081/rooms");
        if (!res.ok) throw new Error(await res.text());
        const data = await res.json();
        setCars(data);
      } catch (e: any) {
        setError(e.message || "Failed to load rooms");
      }
    })();
  }, []);

  return (
    <main style={{ padding: 24 }}>
      <h1>Cars (Matchmaker)</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <table cellPadding={8} style={{ borderCollapse: "collapse", border: "1px solid #ddd", minWidth: 500 }}>
        <thead>
        <tr>
          <th align="left">ID</th>
          <th align="left">State</th>
          <th align="left">Assigned To</th>
          <th align="left">Updated</th>
        </tr>
        </thead>
        <tbody>
        {cars.map((c) => (
          <tr key={c.id}>
            <td>{c.id}</td>
            <td>{c.state}</td>
            <td>{c.assignedTo || "-"}</td>
            <td>{new Date(c.updatedAt).toLocaleString()}</td>
          </tr>
        ))}
        </tbody>
      </table>
    </main>
  );
}
