'use client';

import { useEffect, useState } from 'react';

export default function DndForestUI() {
  const [player, setPlayer] = useState('');
  const [playerSet, setPlayerSet] = useState(false);
  const [story, setStory] = useState<string[]>([]);
  const [input, setInput] = useState('');
  const [scene, setScene] = useState('forest'); // default scene

  const sceneBackgrounds: Record<string, string> = {
    forest: '/scenes/forest.jpg',
    cave: '/scenes/cave.jpg',
    river: '/scenes/river.jpg',
    village: '/scenes/village.jpg',
  };

  // Start game once player is set
  useEffect(() => {
    if (!playerSet) return;

    console.log('Starting game for player:', player);

    fetch(`http://localhost:3001/start-game?player=${encodeURIComponent(player)}`)
      .then(res => res.json())
      .then(data => {
        console.log('Start game response:', data);
        if (data.story) {
          setStory([data.story]);  // Reset story with initial line
          setScene('forest');     // Reset scene to default
        } else {
          setStory(['No initial story received.']);
        }
      })
      .catch(err => {
        console.error('Failed to fetch start-game:', err);
        setStory(['Failed to start game.']);
      });
  }, [playerSet]);

  // Detect scene from response text
  const detectScene = (text: string) => {
    const lower = text.toLowerCase();
    if (lower.includes('cave')) return 'cave';
    if (lower.includes('river')) return 'river';
    if (lower.includes('village')) return 'village';
    if (lower.includes('forest')) return 'forest';
    return scene; // keep current if no match
  };

  const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();
  if (!input.trim() || !playerSet) return;

  const action = input.trim();
  setInput('');
  setStory(prev => [...prev, `🧍 You: ${action}`]);

  const res = await fetch('http://localhost:3001/player-action', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ player, action }),
  });

  const data = await res.json();
  if (data.response) {
    // Truncate response if it contains "Player's action:"
    let dmResponse = data.response;
    const cutoffIndex = dmResponse.indexOf("Player's action:");
    if (cutoffIndex !== -1) {
      dmResponse = dmResponse.slice(0, cutoffIndex).trim();
    }

    setStory(prev => [...prev, `🧙 DM: ${dmResponse}`]);
    const nextScene = detectScene(dmResponse);
    setScene(nextScene);
  }
};


  if (!playerSet) {
    // Show input to enter player name
    return (
      <div style={{
        height: '100vh',
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        fontFamily: "'Georgia', serif",
        backgroundColor: '#111',
        color: 'white',
        flexDirection: 'column',
        padding: 20,
      }}>
        <h2>Enter your adventurer name to begin:</h2>
        <input
          type="text"
          value={player}
          onChange={(e) => setPlayer(e.target.value)}
          placeholder="Your name"
          style={{
            padding: 12,
            fontSize: 18,
            borderRadius: 8,
            border: 'none',
            outline: 'none',
            width: 250,
            textAlign: 'center',
          }}
          onKeyDown={(e) => {
            if (e.key === 'Enter' && player.trim()) {
              setPlayerSet(true);
            }
          }}
        />
        <button
          onClick={() => {
            if (player.trim()) setPlayerSet(true);
          }}
          style={{
            marginTop: 15,
            padding: '10px 20px',
            fontSize: 16,
            borderRadius: 8,
            cursor: 'pointer',
          }}
        >
          Start Game
        </button>
      </div>
    );
  }

  // Main game UI after player name is set
  return (
    <div style={{
      height: '100vh',
      width: '100vw',
      backgroundImage: `url(${sceneBackgrounds[scene]})`,
      backgroundSize: 'cover',
      backgroundPosition: 'center',
      transition: 'background-image 1s ease-in-out',
      color: 'white',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      padding: 20,
      fontFamily: "'Georgia', serif",
      position: 'relative',
    }}>
     

      <div id="story-log" style={{
        backgroundColor: 'rgba(0,0,0,0.6)',
        borderRadius: 10,
        padding: 15,
        maxWidth: 600,
        maxHeight: '50vh',
        overflowY: 'auto',
        marginBottom: 20,
        width: '100%',
      }}>
        {story.map((line, idx) => (
          <p key={idx} style={{ margin: '8px 0' }}>📜 {line}</p>
        ))}
      </div>

      <form onSubmit={handleSubmit} style={{ width: '100%', maxWidth: 600 }}>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="What do you do?"
          style={{
            width: '100%',
            padding: 12,
            fontSize: 18,
            borderRadius: 8,
            border: 'none',
            outline: 'none',
            boxShadow: '0 0 5px 2px rgba(0,255,0,0.5)',
            backgroundColor: 'rgba(0,0,0,0.7)',
            color: 'white',
          }}
        />
      </form>
    </div>
  );
}
