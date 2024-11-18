// Authors: Delaney Gillilan
// Icon: material-symbols:settings-input-antenna
// Slug: Merge fine grain signals store data from a server using the Datastar SDK interface
// Description: Merge store data from a server using the Datastar SDK interface

import {
    DefaultMergeSignalsOnlyIfMissing,
    EventTypes,
    InitExpressionFunction,
    WatcherPlugin,
} from "../../../../engine";
import { PLUGIN_WATCHER } from "../../../../engine/client_only_consts";
import { storeFromPossibleContents } from "../../../../utils/signals";
import { isBoolString } from "../../../../utils/text";
import { datastarSSEEventWatcher } from "./sseShared";

export const MergeSignals: WatcherPlugin = {
    pluginType: PLUGIN_WATCHER,
    name: EventTypes.MergeSignals,
    onGlobalInit: async (ctx) => {
        datastarSSEEventWatcher(EventTypes.MergeSignals, ({
            store = "{}",
            onlyIfMissing: onlyIfMissingRaw =
                `${DefaultMergeSignalsOnlyIfMissing}`,
        }) => {
            const onlyIfMissing = isBoolString(onlyIfMissingRaw);
            const fnContents =
                ` return Object.assign({...ctx.store()}, ${store})`;
            try {
                const fn = new Function(
                    "ctx",
                    fnContents,
                ) as InitExpressionFunction;
                const possibleMergeSignals = fn(ctx);
                const actualMergeSignals = storeFromPossibleContents(
                    ctx.store(),
                    possibleMergeSignals,
                    onlyIfMissing,
                );
                ctx.mergeSignals(actualMergeSignals);
                ctx.applyPlugins(document.body);
            } catch (e) {
                console.log(fnContents);
                console.error(e);
                debugger;
            }
        });
    },
};
