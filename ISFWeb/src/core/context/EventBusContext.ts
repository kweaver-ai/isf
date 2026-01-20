import { createContext } from "react";
import { EventEmitter } from "events";

export const EventBusContext = createContext<EventEmitter>(null as any);