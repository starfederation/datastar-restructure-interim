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
            ({ autoRemoveScript: autoRemoveScriptRaw = "true", script }) => {
                const autoRemoveScript = autoRemoveScriptRaw.trim() === "true";
                if (!script?.length) {
                    throw new Error("No script provided");
                }

                const scriptEl = document.createElement("script");
                scriptEl.text = script;
                document.head.appendChild(scriptEl);
                if (autoRemoveScript) {
                    scriptEl.remove();
                }
            },
        );
    },
};
