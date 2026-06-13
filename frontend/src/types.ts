import {runtime} from "../wailsjs/go/models";

export type { runtime };
export type RuntimeResult = runtime.RuntimeResult;
export type ResponseDTO = { text: string; suggestions?: string[] };
export type RuntimeExtension = runtime.RuntimeExtension;
export type ConfidenceData = runtime.ConfidenceData;
export type ClarificationData = runtime.ClarificationData;
export type ClarificationOption = runtime.ClarificationOption;

export type ChatMessage = {
  id: string;
  role: "user" | "assistant";
  text: string;
  suggestions?: string[];
  clarificationOptions?: ClarificationOption[];
};
