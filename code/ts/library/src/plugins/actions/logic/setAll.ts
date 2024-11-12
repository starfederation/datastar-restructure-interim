import { ActionPlugin } from "../../../engine";

export const SetAll: ActionPlugin = {
    pluginType: "action",
    name: "setAll",
    method: (ctx, regexp, newValue) => {
        const re = new RegExp(regexp);
        ctx.walkSignals((name, signal) =>
            re.test(name) && (signal.value = newValue)
        );
    },
};
