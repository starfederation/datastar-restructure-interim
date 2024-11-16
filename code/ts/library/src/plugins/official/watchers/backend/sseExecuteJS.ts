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
