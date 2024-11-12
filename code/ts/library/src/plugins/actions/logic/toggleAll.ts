import { ActionPlugin } from "../../../engine";

export const ToggleAll: ActionPlugin = {
    pluginType: "action",
    name: "toggleAll",
    method: (ctx, regexp) => {
        const re = new RegExp(regexp);
        ctx.walkSignals((name, signal) =>
            re.test(name) && (signal.value = !signal.value)
        );
    },
};
