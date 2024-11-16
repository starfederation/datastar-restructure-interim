// Authors: Delaney Gillilan
// Icon: material-symbols:settings-input-antenna
// Slug: Merge fine grain signals store data from a server using the Datastar SDK interface
// Description: Merge store data from a server using the Datastar SDK interface

import { EffectPlugin, InitExpressionFunction } from "../../../../engine";
import { storeFromPossibleContents } from "../../../../utils/signals";
import { datastarSSEEventWatcher } from "./sseShared";

const name = "mergeStore";
export const MergeStore: EffectPlugin = {
    pluginType: "effect",
    name,
    onGlobalInit: async (ctx) => {
        datastarSSEEventWatcher(ctx, name, ({
            store = "{}",
            onlyIfMissing: onlyIfMissingRaw = "false",
        }) => {
            const onlyIfMissing = onlyIfMissingRaw.trim() === "true";

            const fnContents =
                ` return Object.assign({...ctx.store()}, ${store})`;
            try {
                const fn = new Function(
                    "ctx",
                    fnContents,
                ) as InitExpressionFunction;
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
        });
    },
};
