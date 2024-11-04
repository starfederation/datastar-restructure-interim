import { ActionPlugin } from "library/engine";
import { remoteSignals } from "library/utils/signals";

export const RemoteActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "remote",
    method: async (ctx) => {
        return remoteSignals(ctx.store().value);
    },
};
