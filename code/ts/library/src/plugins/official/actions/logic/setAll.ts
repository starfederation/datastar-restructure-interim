// Authors: Delaney Gillilan
// Icon: ion:checkmark-round
// Slug: Set all signals that match a regular expression

import { ActionPlugin } from "../../../../engine";
import { PLUGIN_ACTION } from "../../../../engine/client_only_consts";

export const SetAll: ActionPlugin = {
    pluginType: PLUGIN_ACTION,
    name: "setAll",
    method: (ctx, regexp, newValue) => {
        const re = new RegExp(regexp);
        ctx.walkSignals((name, signal) =>
            re.test(name) && (signal.value = newValue)
        );
    },
};
