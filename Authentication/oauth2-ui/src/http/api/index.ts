import { EFAST } from "./efast";
export * from "./efast";

export interface Error {
    code: number;
    message: string;
    cause: string;
    detail?: {
        [k: string]: any;
    };
    [k: string]: any;
}

export default interface API extends EFAST {}
