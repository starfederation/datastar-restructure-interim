// Authors: Delaney Gillilan
// Icon: tabler:file-type-js
// Slug: Execute JavaScript from a Server-Sent Event
// Description: Execute JavaScript from a Server-Sent Event

import { WatcherPlugin } from "../../../../engine";
import { datastarSSEEventWatcher } from "./sseShared";

const name = "executeJs";
export const ExecuteJS: WatcherPlugin = {
    pluginType: "effect",
    name,
    onGlobalInit: async (ctx) => {
        datastarSSEEventWatcher(
            ctx,
            name,
            ({
                autoRemove: autoRemoveRaw = "true",
                type = "module",
                script,
            }) => {
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
