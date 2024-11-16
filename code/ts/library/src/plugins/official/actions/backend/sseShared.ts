import { ActionMethod, DATASTAR } from "../../../../engine";
import { remoteSignals } from "../../../../utils/signals";
import {
    fetchEventSource,
    FetchEventSourceInit,
} from "../../../../vendored/fetch-event-source";
import { Signal } from "../../../../vendored/preact-core";
import {
    INDICATOR_CLASS,
    INDICATOR_LOADING_CLASS,
} from "../../attributes/backend/fetchIndicator";
import {
    DATASTAR_SSE_EVENT,
    DatastarSSEEvent,
    DEFAULT_SETTLE_DURATION_RAW,
} from "../../effects/backend/sseShared";

export type IndicatorReference = { el: HTMLElement; count: number };

const DEFAULT_SETTLE_DURATION = parseInt(DEFAULT_SETTLE_DURATION_RAW);
const isWrongContent = (err: any) =>
    `${err}`.includes(
        `Expected content-type to be text/event-stream`,
    );

export function sendSSERequest(
    method: string,
): ActionMethod {
    return async (ctx, url, onlyRemotes = true) => {
        if (!!!url?.length) {
            throw new Error("URL is required for SSE request");
        }

        const currentStore = ctx.store().value;
        let store = Object.assign({}, currentStore);
        if (onlyRemotes) {
            store = remoteSignals(store);
        }
        const storeJSON = JSON.stringify(store);

        const sendFromElement = ctx.el as HTMLElement;
        ctx.sendDatastarEvent(
            "plugin",
            "backend",
            "fetch_start",
            sendFromElement,
            JSON.stringify({ method, url, onlyRemotes, storeJSON }),
        );

        const indicatorElements: HTMLElement[] =
            store?._dsPlugins?.fetch?.indicatorElements
                ? store?._dsPlugins?.fetch
                    ?.indicatorElements?.[sendFromElement.id]
                    ?.value || []
                : [];
        const indicatorsVisible: Signal<IndicatorReference[]> | undefined =
            store
                ?._dsPlugins?.fetch?.indicatorsVisible;

        indicatorElements.forEach((indicator) => {
            if (!indicator || !indicatorsVisible) {
                return;
            }
            const indicatorVisibleIndex = indicatorsVisible.value.findIndex(
                (indicatorVisible) => {
                    if (!indicatorVisible) {
                        return false;
                    }
                    return indicator.isSameNode(indicatorVisible.el);
                },
            );
            if (indicatorVisibleIndex > -1) {
                const indicatorVisible =
                    indicatorsVisible.value[indicatorVisibleIndex];
                const indicatorsVisibleNew = [...indicatorsVisible.value];
                delete indicatorsVisibleNew[indicatorVisibleIndex];
                indicatorsVisible.value = [
                    ...indicatorsVisibleNew.filter((indicator) => {
                        return !!indicator;
                    }),
                    { el: indicator, count: indicatorVisible.count + 1 },
                ];
            } else {
                indicator.classList.remove(INDICATOR_CLASS);
                indicator.classList.add(INDICATOR_LOADING_CLASS);
                indicatorsVisible.value = [
                    ...indicatorsVisible.value,
                    {
                        el: indicator,
                        count: 1,
                    },
                ];
            }
        });

        const urlInstance = new URL(url, window.location.origin);
        const req: FetchEventSourceInit = {
            method,
            headers: {
                ["Content-Type"]: "application/json",
                [`${DATASTAR}-request`]: "true",
            },
            onmessage: (evt) => {
                if (!evt.event.startsWith(DATASTAR)) {
                    return;
                }
                const type = evt.event.slice(DATASTAR.length + 1);
                const argsRawLines: Record<string, string[]> = {};

                const lines = evt.data.split("\n");
                for (const line of lines) {
                    const colonIndex = line.indexOf(" ");
                    const key = line.slice(0, colonIndex);
                    let argLines = argsRawLines[key];
                    if (!argLines) {
                        argLines = [];
                        argsRawLines[key] = argLines;
                    }
                    const value = line.slice(colonIndex + 1).trim();
                    argLines.push(value);
                }

                const argsRaw: Record<string, string> = {};
                for (const [key, lines] of Object.entries(argsRawLines)) {
                    argsRaw[key] = lines.join("\n");
                }

                // if you aren't seeing your event you can debug by using this line in the console
                // document.addEventListener("datastar-sse",(e) => console.log(e));

                const datastarSSEEvent = new CustomEvent<DatastarSSEEvent>(
                    DATASTAR_SSE_EVENT,
                    {
                        detail: { type, argsRaw },
                    },
                );
                document.dispatchEvent(datastarSSEEvent);
            },
            onerror: (err) => {
                if (isWrongContent(err)) {
                    // don't retry if the content-type is wrong
                    throw err;
                }
                // do nothing and it will retry
                console.error(`Error with SSE request: ${err.message}`);
            },
            onclose: () => {
                try {
                    const store = ctx.store();
                    const indicatorsVisible: Signal<IndicatorReference[]> =
                        store?._dsPlugins?.fetch?.indicatorsVisible || [];
                    const indicatorElements: HTMLElement[] =
                        store?._dsPlugins?.fetch?.indicatorElements
                            ? store._dsPlugins.fetch
                                .indicatorElements[sendFromElement.id]?.value ||
                                []
                            : [];
                    const indicatorCleanupPromises: Promise<() => void>[] = [];
                    if (indicatorElements?.forEach) {
                        indicatorElements.forEach((indicator) => {
                            if (!indicator || !indicatorsVisible) {
                                return;
                            }
                            const indicatorsVisibleNew =
                                indicatorsVisible.value;
                            const indicatorVisibleIndex = indicatorsVisibleNew
                                .findIndex((indicatorVisible) => {
                                    if (!indicatorVisible) {
                                        return false;
                                    }
                                    return indicator.isSameNode(
                                        indicatorVisible.el,
                                    );
                                });
                            const indicatorVisible =
                                indicatorsVisibleNew[indicatorVisibleIndex];
                            if (!indicatorVisible) {
                                return;
                            }
                            if (indicatorVisible.count < 2) {
                                indicatorCleanupPromises.push(
                                    new Promise(() =>
                                        setTimeout(() => {
                                            indicator.classList.remove(
                                                INDICATOR_LOADING_CLASS,
                                            );
                                            indicator.classList.add(
                                                INDICATOR_CLASS,
                                            );
                                        }, DEFAULT_SETTLE_DURATION)
                                    ),
                                );
                                delete indicatorsVisibleNew[
                                    indicatorVisibleIndex
                                ];
                            } else if (indicatorVisibleIndex > -1) {
                                indicatorsVisibleNew[indicatorVisibleIndex]
                                    .count = indicatorsVisibleNew[
                                        indicatorVisibleIndex
                                    ]
                                        .count - 1;
                            }
                            indicatorsVisible.value = indicatorsVisibleNew
                                .filter(
                                    (indicator) => {
                                        return !!indicator;
                                    },
                                );
                        });
                    }

                    Promise.all(indicatorCleanupPromises);
                } catch (e) {
                    console.error(e);
                    debugger;
                } finally {
                    ctx.sendDatastarEvent(
                        "plugin",
                        "backend",
                        "fetch_end",
                        sendFromElement,
                        JSON.stringify({ method, url }),
                    );
                }
            },
        };

        if (method === "GET") {
            const queryParams = new URLSearchParams(urlInstance.search);
            queryParams.append("datastar", storeJSON);
            urlInstance.search = queryParams.toString();
        } else {
            req.body = storeJSON;
        }

        const headers = store?._dsPlugins?.fetch?.headers || {};
        if (req.headers) {
            for (const [key, value] of Object.entries(headers)) {
                if (key.startsWith("_")) {
                    continue;
                }
                req.headers[key] = `${value}`;
            }
        }

        try {
            const revisedURL = urlInstance.toString();
            await fetchEventSource(revisedURL, req);
        } catch (err) {
            if (!isWrongContent(err)) {
                throw new Error(`Failed to fetch ${url}: ${err}`);
            }

            // exit gracefully and do nothing if the content-type is wrong
            // this can happen if the client is sending a request
            // where no response is expected, and they haven't
            // set the content-type to text/event-stream
        }
    };
}
