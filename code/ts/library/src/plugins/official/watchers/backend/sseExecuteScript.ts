// Authors: Delaney Gillilan
// Icon: tabler:file-type-js
// Slug: Execute JavaScript from a Server-Sent Event
// Description: Execute JavaScript from a Server-Sent Event

import {
    DefaultExecuteScriptAttributes,
    EventTypes,
    WatcherPlugin,
} from "../../../../engine";
import { PLUGIN_WATCHER } from "../../../../engine/client_only_consts";
import { isBoolString } from "../../../../utils/text";
import { datastarSSEEventWatcher } from "./sseShared";

export const ExecuteJS: WatcherPlugin = {
    pluginType: PLUGIN_WATCHER,
    name: EventTypes.ExecuteScript,
    onGlobalInit: async () => {
        datastarSSEEventWatcher(
            EventTypes.ExecuteScript,
            (
                {
                    autoRemove: autoRemoveRaw = "true",
                    attribute: attributesRaw = DefaultExecuteScriptAttributes,
                    script,
                },
            ) => {
                const autoRemove = isBoolString(autoRemoveRaw);
                if (!script?.length) {
                    throw new Error("No script provided");
                }
                const scriptEl = document.createElement("script");
                attributesRaw.split("\n").forEach((attr) => {
                    const pivot = attr.indexOf(" ");
                    const key = attr.slice(0, pivot).trim();
                    const value = attr.slice(pivot).trim();
                    scriptEl.setAttribute(key, value);
                });
                scriptEl.text = script;
                document.head.appendChild(scriptEl);
                if (autoRemove) {
                    scriptEl.remove();
                }
            },
        );
    },
};
