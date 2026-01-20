import { useContext } from "react";
import { EventBusContext } from "../context/EventBusContext";

export const useEventBus: any = () => useContext(EventBusContext);
