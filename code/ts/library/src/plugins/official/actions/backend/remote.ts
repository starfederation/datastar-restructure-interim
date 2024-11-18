// Authors: Delaney Gillilan
// Icon: material-symbols:public
// Slug: Filter to only signals that should be sent to the server
// Description: This is a nested action that filters out signals that should not be sent to the server.

import { ActionPlugin } from "../../../../engine";
import { PLUGIN_ACTION } from "../../../../engine/client_only_consts";
import { remoteSignals } from "../../../../utils/signals";

export const RemoteSignals: ActionPlugin = {
    pluginType: PLUGIN_ACTION,
    name: "remote",
    method: async (ctx) => {
        return remoteSignals(ctx.store().value);
    },
};
