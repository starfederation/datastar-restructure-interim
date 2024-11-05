import { ActionPlugin } from "library/src/engine";
import { remoteSignals } from "library/src/utils/signals";

export const RemoteActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "remote",
    method: async (ctx) => {
        return remoteSignals(ctx.store().value);
    },
};
