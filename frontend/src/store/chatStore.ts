import { create } from "zustand";
import type { ChatMessage } from "../types";
import { SendMessage } from "../../wailsjs/go/main/App";

interface ChatStore {
  messages: ChatMessage[];
  loading: boolean;
  sendMessage: (text: string) => Promise<void>;
}

export const useChatStore = create<ChatStore>((set) => ({
  messages: [],
  loading: false,
  sendMessage: async (text: string) => {
    // Add user message
    const userMsg: ChatMessage = {
      id: Date.now().toString(),
      role: "user",
      text: text,
    };
    set((state) => ({ messages: [...state.messages, userMsg], loading: true }));

    console.log("input text: ", text);
    const result = await SendMessage(text);
    // Add assistant message
    const assistantMsg: ChatMessage = {
      id: (Date.now() + 1).toString(),
      role: "assistant",
      text: result.response.text,
      suggestions: result.response.suggestions || undefined,
      clarificationOptions: result.extension?.clarify?.options || undefined,
    };

    set((state) => ({
      messages: [...state.messages, assistantMsg],
      loading: false,
    }));

    set({ loading: false });
  },
}));
