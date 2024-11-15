export const DATASTAR_SSE_EVENT = "datastar-sse";

export interface DatastarSSEEvent {
    type: string;
    argsRaw: Record<string, string>;
}

export interface CustomEventMap {
    "datastar-sse": CustomEvent<DatastarSSEEvent>;
}

declare global {
    interface Document { //adds definition to Document, but you can do the same with HTMLElement
        addEventListener<K extends keyof CustomEventMap>(
            type: K,
            listener: (this: Document, ev: CustomEventMap[K]) => void,
        ): void;
        dispatchEvent<K extends keyof CustomEventMap>(
            ev: CustomEventMap[K],
        ): void;
    }
}
