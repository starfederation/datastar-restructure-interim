// Authors: Delaney Gillilan
// Icon: tabler:file-type-js
// Slug: Execute JavaScript from a Server-Sent Event
// Description: Execute JavaScript from a Server-Sent Event

import { EventTypes, WatcherPlugin } from "../../../../engine";
import { PLUGIN_WATCHER } from "../../../../engine/client_only_consts";
import { datastarSSEEventWatcher } from "./sseShared";

export const ExecuteJS: WatcherPlugin = {
    pluginType: PLUGIN_WATCHER,
    name: EventTypes.ExecuteJs,
    onGlobalInit: async () => {
        datastarSSEEventWatcher(
            EventTypes.ExecuteJs,
            (
                { autoRemove: autoRemoveRaw = "true", type = "module", script },
            ) => {
                const autoRemove = autoRemoveRaw.trim() === "true";
                if (!script?.length) {
                    throw new Error("No script provided");
                }

                const scriptEl = document.createElement("script");
                scriptEl.type = type.trim();
                scriptEl.text = script;
                document.head.appendChild(scriptEl);
                if (autoRemove) {
                    scriptEl.remove();
                }
            },
        );
    },
};
