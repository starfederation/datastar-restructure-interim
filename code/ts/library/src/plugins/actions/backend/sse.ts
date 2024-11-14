import {
    ActionMethod,
    AttributeContext,
    DATASTAR,
    ExpressionFunction,
} from "../../../engine";
import {
    remoteSignals,
    storeFromPossibleContents,
} from "../../../utils/signals";
import {
    docWithViewTransitionAPI,
    supportsViewTransitions,
} from "../../../utils/view-transitions";
import {
    fetchEventSource,
    FetchEventSourceInit,
} from "../../../vendored/fetch-event-source";
import { idiomorph } from "../../../vendored/idiomorph";
import { Signal } from "../../../vendored/preact-core";
import {
    INDICATOR_CLASS,
    INDICATOR_LOADING_CLASS,
} from "../../attributes/backend/fetchIndicator";

const DEFAULT_MERGE_MODE: FragmentMergeMode = "morph";
const DEFAULT_SETTLE_DURATION = 300;
const DEFAULT_USE_VIEW_TRANSITION = false;

const FragmentMergeModes = {
    MorphElement: "morph",
    InnerElement: "inner",
    OuterElement: "outer",
    PrependElement: "prepend",
    AppendElement: "append",
    BeforeElement: "before",
    AfterElement: "after",
    UpsertAttributes: "upsertAttributes",
} as const;
export type FragmentMergeMode =
    (typeof FragmentMergeModes)[keyof typeof FragmentMergeModes];

export function fetcherActionMethod(method: string): ActionMethod {
    return (ctx, urlExpression, onlyRemoteRaw) => {
        const onlyRemotes = ["true", true, undefined].includes(onlyRemoteRaw);
        fetcher(method, urlExpression, ctx, onlyRemotes);
    };
}

type IndicatorReference = { el: HTMLElement; count: number };

async function fetcher(
    method: string,
    urlExpression: string,
    ctx: AttributeContext,
    onlyRemote = true,
) {
    const store = ctx.store();

    if (!urlExpression) {
        throw new Error(`No signal for ${method} on ${urlExpression}`);
    }

    let storeValue = { ...store.value };
    if (onlyRemote) {
        storeValue = remoteSignals(storeValue);
    }
    const storeJSON = JSON.stringify(storeValue);

    const loadingTarget = ctx.el as HTMLElement;

    ctx.sendDatastarEvent(
        "plugin",
        "backend",
        "fetch_start",
        loadingTarget,
        JSON.stringify({ method, urlExpression, onlyRemote, storeJSON }),
    );
    const indicatorElements: HTMLElement[] =
        store?._dsPlugins?.fetch?.indicatorElements
            ? store._dsPlugins.fetch.indicatorElements[loadingTarget.id]
                ?.value || []
            : [];
    const indicatorsVisible: Signal<IndicatorReference[]> | undefined = store
        ?._dsPlugins.fetch?.indicatorsVisible;
    if (!!indicatorElements?.forEach) {
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
    }

    const isWrongContent = (err: any) =>
        `${err}`.includes(
            `Expected content-type to be text/event-stream`,
        );

    const url = new URL(urlExpression, window.location.origin);
    method = method.toUpperCase();
    const req: FetchEventSourceInit = {
        method,
        headers: {
            ["Content-Type"]: "application/json",
            [`${DATASTAR}-request`]: "true",
        },
        onmessage: (evt) => {
            if (!evt.event) {
                return;
            } else if (!evt.event.startsWith(DATASTAR)) {
                console.log(`Unknown event: ${evt.event}`);
                debugger;
            }

            switch (evt.event) {
                case `${DATASTAR}-fragment`:
                    const fragmentLines = evt.data.trim().split("\n");
                    const knownEventTypes = [
                        "selector",
                        "mergeMode",
                        "settleDuration",
                        "useViewTransition",
                        "fragment",
                    ];

                    let fragment = "",
                        selector = "",
                        mergeMode = DEFAULT_MERGE_MODE,
                        settleDuration = DEFAULT_SETTLE_DURATION,
                        useViewTransition = DEFAULT_USE_VIEW_TRANSITION;

                    for (let i = 0; i < fragmentLines.length; i++) {
                        const line = fragmentLines[i];
                        const [dataType, ...rest] = line.split(" ");
                        const data = rest.join(" ");
                        const isDatatype = knownEventTypes.includes(dataType);
                        if (!isDatatype) {
                            throw new Error(`Unknown data type: '${dataType}'`);
                        }

                        switch (dataType) {
                            case "selector":
                                selector = data;
                                break;
                            case "mergeMode":
                                mergeMode = data as FragmentMergeMode;
                                if (
                                    !Object.values(FragmentMergeModes)
                                        .includes(mergeMode)
                                ) {
                                    throw new Error(
                                        `Unknown merge option: ${mergeMode}`,
                                    );
                                }
                                break;
                            case "settleDuration":
                                settleDuration = parseInt(data);
                                break;
                            case "useViewTransition":
                                useViewTransition = data.trim() === "true";
                                break;
                            case "fragment":
                                fragment += data + "\n";
                                break;
                            default:
                                throw new Error(`Unknown data type`);
                        }
                    }

                    if (!fragment?.length) {
                        fragment = "<div></div>";
                    }
                    mergeHTMLFragment(
                        ctx,
                        selector,
                        mergeMode,
                        settleDuration,
                        useViewTransition,
                        fragment,
                    );
                    ctx.sendDatastarEvent(
                        "plugin",
                        "backend",
                        "merge",
                        selector,
                        JSON.stringify({
                            fragment,
                            settleDuration,
                            useViewTransition,
                        }),
                    );

                    break;

                case `${DATASTAR}-signal`:
                    let onlyIfMissing = false;
                    let storeToMerge = "";

                    const signalLines = evt.data.trim().split("\n");
                    for (let i = 0; i < signalLines.length; i++) {
                        const line = signalLines[i];
                        const [signalType, ...signalRest] = line.split(" ");
                        const signalLine = signalRest.join(" ");
                        switch (signalType) {
                            case "onlyIfMissing":
                                onlyIfMissing = signalLine.trim() === "true";
                                break;
                            case "store":
                                storeToMerge += `${signalLine}\n`;
                                break;
                            default:
                                throw new Error(
                                    `Unknown signal type: ${signalType}`,
                                );
                        }
                    }

                    const fnContents =
                        ` return Object.assign({...ctx.store()}, ${storeToMerge})`;
                    try {
                        const fn = new Function(
                            "ctx",
                            fnContents,
                        ) as ExpressionFunction;
                        const possibleMergeStore = fn(ctx);
                        const actualMergeStore = storeFromPossibleContents(
                            ctx.store(),
                            possibleMergeStore,
                            onlyIfMissing,
                        );
                        ctx.mergeStore(actualMergeStore);
                        ctx.applyPlugins(document.body);
                    } catch (e) {
                        console.log(fnContents);
                        console.error(e);
                        debugger;
                    }
                    break;

                case `${DATASTAR}-remove`:
                    const [removePrefix, ...removeRest] = evt.data.trim().split(
                        " ",
                    );

                    switch (removePrefix) {
                        case "selector":
                            const removeSelector = removeRest.join(" ");
                            const removeTargets = document.querySelectorAll(
                                removeSelector,
                            );
                            removeTargets.forEach((target) => target.remove());
                            break;
                        case "paths":
                            const paths = removeRest.join(" ").split(" ");
                            ctx.removeFromStore(...paths);
                            break;
                        case "settleDuration":
                            // TODO: Implement
                            break;
                        case "useViewTransition":
                            // TODO: Implement
                            break;
                        default:
                            throw new Error(
                                `Unknown delete prefix: ${removePrefix}`,
                            );
                    }
                    break;

                case `${DATASTAR}-redirect`:
                    const [redirectSelector, ...redirectRest] = evt.data.trim()
                        .split(" ");
                    if (redirectSelector !== "url") {
                        throw new Error(
                            `Unknown redirect selector: ${redirectSelector}`,
                        );
                    }
                    const redirectTarget = redirectRest.join(" ");
                    ctx.sendDatastarEvent(
                        "plugin",
                        "backend",
                        "redirect",
                        "WINDOW",
                        redirectTarget,
                    );
                    window.location.href = redirectTarget;
                    break;

                case `${DATASTAR}-console`:
                    const [consoleMode, ...consoleRest] = evt.data.trim().split(
                        " ",
                    );
                    const consoleMessage = consoleRest.join(" ");
                    switch (consoleMode) {
                        case "assert":
                        case "clear":
                        case "count":
                        case "countReset":
                        case "debug":
                        case "dir":
                        case "dirxml":
                        case "error":
                        case "group":
                        case "groupCollapsed":
                        case "groupEnd":
                        case "info":
                        case "log":
                        case "table":
                        case "time":
                        case "timeEnd":
                        case "timeLog":
                        case "trace":
                        case "warn":
                            const fn = console[consoleMode] as Function;
                            fn(consoleMessage);
                            break;
                        default:
                            throw new Error(
                                `Unknown console mode: '${consoleMode}', message: '${consoleMessage}'`,
                            );
                    }
            }
        },
        onerror: (err) => {
            if (isWrongContent(err)) {
                // don't retry if the content-type is wrong
                throw err;
            }

            // do nothing and it will retry
        },
        onclose: () => {
            try {
                const store = ctx.store();
                const indicatorsVisible: Signal<IndicatorReference[]> =
                    store?._dsPlugins?.fetch?.indicatorsVisible || [];
                const indicatorElements: HTMLElement[] =
                    store?._dsPlugins?.fetch?.indicatorElements
                        ? store._dsPlugins.fetch
                            .indicatorElements[loadingTarget.id]?.value || []
                        : [];
                const indicatorCleanupPromises: Promise<() => void>[] = [];
                if (indicatorElements?.forEach) {
                    indicatorElements.forEach((indicator) => {
                        if (!indicator || !indicatorsVisible) {
                            return;
                        }
                        const indicatorsVisibleNew = indicatorsVisible.value;
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
                            delete indicatorsVisibleNew[indicatorVisibleIndex];
                        } else if (indicatorVisibleIndex > -1) {
                            indicatorsVisibleNew[indicatorVisibleIndex].count =
                                indicatorsVisibleNew[indicatorVisibleIndex]
                                    .count - 1;
                        }
                        indicatorsVisible.value = indicatorsVisibleNew.filter(
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
                    loadingTarget,
                    JSON.stringify({ method, urlExpression }),
                );
            }
        },
    };

    if (method === "GET") {
        const queryParams = new URLSearchParams(url.search);
        queryParams.append("datastar", storeJSON);
        url.search = queryParams.toString();
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
        await fetchEventSource(url.toString(), req);
    } catch (err) {
        if (!isWrongContent(err)) {
            throw err;
        }

        // exit gracefully and do nothing if the content-type is wrong
        // this can happen if the client is sending a request
        // where no response is expected, and they haven't
        // set the content-type to text/event-stream
    }
}

const SETTLING_CLASS = `${DATASTAR}-settling`;
const SWAPPING_CLASS = `${DATASTAR}-swapping`;

const fragmentContainer = document.createElement("template");
export function mergeHTMLFragment(
    ctx: AttributeContext,
    selector: string,
    mergeMode: FragmentMergeMode,
    settleDuration: number,
    useViewTransition: boolean,
    fragmentsRaw: string,
) {
    const { el } = ctx;

    fragmentContainer.innerHTML = fragmentsRaw.trim();
    const fragments = [...fragmentContainer.content.children];
    fragments.forEach((fragment) => {
        if (!(fragment instanceof Element)) {
            throw new Error(`No fragment found`);
        }
        const applyToTargets = (capturedTargets: Element[]) => {
            for (const initialTarget of capturedTargets) {
                initialTarget.classList.add(SWAPPING_CLASS);
                const originalHTML = initialTarget.outerHTML;
                let modifiedTarget = initialTarget;
                switch (mergeMode) {
                    case FragmentMergeModes.MorphElement:
                        const result = idiomorph(modifiedTarget, fragment, {
                            callbacks: {
                                beforeNodeRemoved: (
                                    oldNode: Element,
                                    _: Element,
                                ) => {
                                    ctx.cleanupElementRemovals(oldNode);
                                    return true;
                                },
                            },
                        });
                        if (!result?.length) {
                            throw new Error(`No morph result `);
                        }
                        const first = result[0] as Element;
                        modifiedTarget = first;
                        break;
                    case FragmentMergeModes.InnerElement:
                        // Replace the contents of the target element with the response
                        modifiedTarget.innerHTML = fragment.innerHTML;
                        break;
                    case FragmentMergeModes.OuterElement:
                        // Replace the entire target element with the response
                        modifiedTarget.replaceWith(fragment);
                        break;
                    case FragmentMergeModes.PrependElement:
                        // Insert the response before the first child of the target element
                        modifiedTarget.prepend(fragment);
                        break;
                    case FragmentMergeModes.AppendElement:
                        // Insert the response after the last child of the target element
                        modifiedTarget.append(fragment);
                        break;
                    case FragmentMergeModes.BeforeElement:
                        // Insert the response before the target element
                        modifiedTarget.before(fragment);
                        break;
                    case FragmentMergeModes.AfterElement:
                        // Insert the response after the target element
                        modifiedTarget.after(fragment);
                        break;
                    case FragmentMergeModes.UpsertAttributes:
                        // Upsert the attributes of the target element
                        fragment.getAttributeNames().forEach((attrName) => {
                            const value = fragment.getAttribute(attrName)!;
                            modifiedTarget.setAttribute(attrName, value);
                        });
                        break;
                    default:
                        throw new Error(`Unknown merge type: ${mergeMode}`);
                }
                ctx.cleanupElementRemovals(modifiedTarget);
                modifiedTarget.classList.add(SWAPPING_CLASS);

                ctx.applyPlugins(document.body);

                setTimeout(() => {
                    initialTarget.classList.remove(SWAPPING_CLASS);
                    modifiedTarget.classList.remove(SWAPPING_CLASS);
                }, settleDuration);

                const revisedHTML = modifiedTarget.outerHTML;

                if (originalHTML !== revisedHTML) {
                    modifiedTarget.classList.add(SETTLING_CLASS);
                    setTimeout(() => {
                        modifiedTarget.classList.remove(SETTLING_CLASS);
                    }, settleDuration);
                }
            }
        };

        const useElAsTarget = selector === "self";

        let targets: Iterable<Element>;
        if (useElAsTarget) {
            targets = [el];
        } else {
            const selectorOrID = selector || `#${fragment.getAttribute("id")}`;
            targets = document.querySelectorAll(selectorOrID) || [];
            if (!!!targets) {
                throw new Error(`No targets found for ${selectorOrID}`);
            }
        }
        const allTargets = [...targets];
        if (!allTargets.length) {
            throw new Error(`No targets found for ${selector}`);
        }

        if (supportsViewTransitions && useViewTransition) {
            docWithViewTransitionAPI.startViewTransition(() =>
                applyToTargets(allTargets)
            );
        } else {
            applyToTargets(allTargets);
        }
    });
}
