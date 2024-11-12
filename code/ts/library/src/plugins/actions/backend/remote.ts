import { ActionPlugin } from "../../../engine";
import { remoteSignals } from "../../../utils/signals";

export const RemoteActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "remote",
    method: async (ctx) => {
        return remoteSignals(ctx.store().value);
    },
};
