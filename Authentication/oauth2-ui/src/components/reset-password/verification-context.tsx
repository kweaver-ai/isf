import { createContext, useContext } from "react";
import { IForgetPasswordState } from "./type";
export const VerificationContext = createContext<null | IForgetPasswordState>(null);
export const useVerification = () => useContext(VerificationContext);
