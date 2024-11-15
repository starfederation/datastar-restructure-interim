// Authors: Delaney Gillilan
// Icon: material-symbols:cloud-download
// Slug: Use Server-Sent Events to fetch data from a server using the Datastar SDK interface
// Description: Remember, SSE is just a regular SSE request but with the ability to send 0-inf messages to the client.

import { AttribtueExpressionFunction, EffectPlugin } from "../../../engine";
import { storeFromPossibleContents } from "../../../utils/signals";
import { DATASTAR_SSE_EVENT, DatastarSSEEvent } from "./types";

export const MergeHTMLFragments: EffectPlugin = {
    pluginType: "effect",
    name: "signals",
    onGlobalInit: async (ctx) => {
        document.addEventListener(
            DATASTAR_SSE_EVENT,
            (event: CustomEvent<DatastarSSEEvent>) => {
                if (event.detail.type != "signals") return;

                const {
                    store = "{}",
                    onlyIfMissing: onlyIfMissingRaw = "false",
                } = event.detail.argsRaw;

                const onlyIfMissing = onlyIfMissingRaw.trim() === "true";

                const fnContents =
                    ` return Object.assign({...ctx.store()}, ${store})`;
                try {
                    const fn = new Function(
                        "ctx",
                        fnContents,
                    ) as AttribtueExpressionFunction;
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
            },
        );
    },
};
