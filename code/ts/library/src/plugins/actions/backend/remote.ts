import { remoteSignals } from "library/src/utils/signals";
import { ActionPlugin } from "../../../engine";

export const RemoteActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "remote",
    method: async (ctx) => {
        return remoteSignals(ctx.store().value);
    },
};
