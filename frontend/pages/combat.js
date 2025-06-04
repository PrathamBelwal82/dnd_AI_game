import { useEffect, useRef, useState } from "react";

export default function Combat() {
  const canvasRef = useRef(null);

  // States
  const [playerHP, setPlayerHP] = useState(100);
  const [enemyHP, setEnemyHP] = useState(60);
  const [combatLog, setCombatLog] = useState(["The battle begins!"]);
  const [loading, setLoading] = useState(false);

  // Draw UI on canvas
  useEffect(() => {
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");

    // Clear
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Background
    ctx.fillStyle = "#222";
    ctx.fillRect(0, 0, canvas.width, canvas.height);

    // Draw player HP bar
    ctx.fillStyle = "green";
    ctx.fillRect(50, 50, (playerHP / 100) * 200, 20);
    ctx.strokeStyle = "white";
    ctx.strokeRect(50, 50, 200, 20);
    ctx.fillStyle = "white";
    ctx.font = "16px Arial";
    ctx.fillText(`Player HP: ${playerHP}`, 50, 45);

    // Draw enemy HP bar
    ctx.fillStyle = "red";
    ctx.fillRect(50, 100, (enemyHP / 60) * 200, 20);
    ctx.strokeStyle = "white";
    ctx.strokeRect(50, 100, 200, 20);
    ctx.fillStyle = "white";
    ctx.font = "16px Arial";
    ctx.fillText(`Goblin HP: ${enemyHP}`, 50, 95);
  }, [playerHP, enemyHP]);

  // Send action to backend
  async function sendAction(action) {
    setLoading(true);
    try {
      const res = await fetch("/combat/turn", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ action }),
      });
      const data = await res.json();
      const responseText = data.response;

      // Simple parsing for damage (you can improve this)
      const playerDamageMatch = responseText.match(/dealt (\d+) damage to Goblin/i);
      const enemyDamageMatch = responseText.match(/deals (\d+) damage to you/i);

      if (playerDamageMatch) {
        const dmg = parseInt(playerDamageMatch[1]);
        setEnemyHP((hp) => Math.max(hp - dmg, 0));
      }
      if (enemyDamageMatch) {
        const dmg = parseInt(enemyDamageMatch[1]);
        setPlayerHP((hp) => Math.max(hp - dmg, 0));
      }

      setCombatLog((log) => [...log, responseText]);
    } catch (e) {
      setCombatLog((log) => [...log, "Error: Could not reach server"]);
    }
    setLoading(false);
  }

  return (
    <div style={{ padding: 20, backgroundColor: "#111", color: "white", height: "100vh" }}>
      <canvas ref={canvasRef} width={300} height={150} style={{ border: "1px solid white" }} />

      <div style={{ marginTop: 20 }}>
        <button onClick={() => sendAction("attack")} disabled={loading || playerHP <= 0 || enemyHP <= 0}>
          Attack
        </button>
        <button onClick={() => sendAction("defend")} disabled={loading || playerHP <= 0 || enemyHP <= 0} style={{ marginLeft: 10 }}>
          Defend
        </button>
        <button onClick={() => sendAction("flee")} disabled={loading || playerHP <= 0 || enemyHP <= 0} style={{ marginLeft: 10 }}>
          Flee
        </button>
      </div>

      <div
        style={{
          marginTop: 20,
          backgroundColor: "#222",
          padding: 10,
          height: 150,
          overflowY: "scroll",
          fontFamily: "monospace",
          whiteSpace: "pre-wrap",
        }}
      >
        {combatLog.map((line, i) => (
          <div key={i}>{line}</div>
        ))}
      </div>

      {(playerHP <= 0 || enemyHP <= 0) && (
        <div style={{ marginTop: 20, fontWeight: "bold" }}>
          {playerHP <= 0 && "You Lost! 💀"}
          {enemyHP <= 0 && "You Won! 🎉"}
        </div>
      )}
    </div>
  );
}
