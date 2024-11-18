// Authors: Delaney Gillilan
// Icon: material-symbols:toggle-off
// Slug: Toggle all signals that match a regular expression

import { ActionPlugin } from "../../../../engine";
import { PLUGIN_ACTION } from "../../../../engine/client_only_consts";

export const ToggleAll: ActionPlugin = {
    pluginType: PLUGIN_ACTION,
    name: "toggleAll",
    method: (ctx, regexp) => {
        const re = new RegExp(regexp);
        ctx.walkSignals((name, signal) =>
            re.test(name) && (signal.value = !signal.value)
        );
    },
};
