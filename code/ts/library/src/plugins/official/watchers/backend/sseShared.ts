import { DATASTAR, InitContext } from "../../../../engine";
import { kebabize } from "../../../../utils/text";

export const DATASTAR_SSE_EVENT = "datastar-sse";
export const DEFAULT_SETTLE_DURATION_RAW = "300";
export const SETTLING_CLASS = `${DATASTAR}-settling`;
export const SWAPPING_CLASS = `${DATASTAR}-swapping`;

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

export function datastarSSEEventWatcher(
    ctx: InitContext,
    name: string,
    fn: (argsRaw: Record<string, string>) => void,
) {
    const kebabName = kebabize(name);
    document.addEventListener(
        DATASTAR_SSE_EVENT,
        (event: CustomEvent<DatastarSSEEvent>) => {
            if (event.detail.type != kebabName) return;
            const { argsRaw } = event.detail;
            fn(argsRaw);

            ctx.sendDatastarEvent(
                "plugin",
                "backend",
                "sse",
                name,
                JSON.stringify(argsRaw),
            );
        },
    );
}
