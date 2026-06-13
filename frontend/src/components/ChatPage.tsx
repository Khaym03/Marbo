import React, { useState } from "react";
import { useChatStore } from "../store/chatStore";
import type { ChatMessage } from "../types";

export const ChatPage: React.FC = () => {
  const { messages, loading, sendMessage } = useChatStore();
  const [input, setInput] = useState("");

  const handleSend = () => {
    if (!input.trim() || loading) return;
    sendMessage(input);
    setInput("");

    messages.map((o) => console.log(o));
  };

  return (
    <div className="flex flex-col h-screen p-4">
      <div className="flex-1 overflow-y-auto mb-4 space-y-4">
        {messages.map((msg: ChatMessage) => (
          <div
            key={msg.id}
            className={`p-2 rounded ${msg.role === "user" ? "bg-blue-100 self-end" : "bg-gray-100 self-start"}`}
          >
            <p>{msg.text}</p>
            {msg.suggestions && (
              <div className="flex gap-2 mt-2">
                {msg.suggestions.map((s: string) => (
                  <button
                    key={s}
                    onClick={() => sendMessage(s)}
                    className="bg-white border rounded px-2"
                  >
                    {s}
                  </button>
                ))}
              </div>
            )}
            {msg.clarificationOptions && (
              <div className="flex flex-col gap-2 mt-2">
                {msg.clarificationOptions.map((o) => (
                  <button
                    key={o.intent_id}
                    onClick={() => sendMessage(o.label)}
                    className="bg-white border rounded px-2 text-left"
                  >
                    {o.label}
                  </button>
                ))}
              </div>
            )}
          </div>
        ))}
        {loading && <div className="text-gray-500">Thinking...</div>}
      </div>
      <div className="flex gap-2">
        <input
          className="border p-2 flex-1"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          disabled={loading}
          onKeyDown={(e) => e.key === "Enter" && handleSend()}
        />
        <button
          className="bg-blue-500 text-white p-2"
          onClick={handleSend}
          disabled={loading}
        >
          Send
        </button>
      </div>
    </div>
  );
};
