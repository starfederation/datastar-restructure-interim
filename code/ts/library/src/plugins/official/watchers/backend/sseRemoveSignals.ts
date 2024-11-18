// Authors: Delaney Gillilan
// Icon: material-symbols:settings-input-antenna
// Slug: Merge fine grain signals store data from a server using the Datastar SDK interface
// Description: Merge store data from a server using the Datastar SDK interface

import { EventTypes, WatcherPlugin } from "../../../../engine";
import { PLUGIN_WATCHER } from "../../../../engine/client_only_consts";
import { datastarSSEEventWatcher } from "./sseShared";

export const RemoveSignals: WatcherPlugin = {
    pluginType: PLUGIN_WATCHER,
    name: EventTypes.RemoveFromStore,
    onGlobalInit: async (ctx) => {
        datastarSSEEventWatcher(
            EventTypes.RemoveFromStore,
            ({ paths: pathsRaw = "" }) => {
                // replace all whitespace with a single space
                pathsRaw = pathsRaw.replaceAll(/\s+/g, " ");
                if (!!!pathsRaw?.length) {
                    throw new Error(
                        "No paths provided for remove-signals",
                    );
                }
                const paths = pathsRaw.split(" ");
                ctx.removeFromStore(...paths);
            },
        );
    },
};
