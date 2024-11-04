import { ActionPlugin, AttributeContext } from "library/engine";

export const FitActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "fit",
    method: (
        _: AttributeContext,
        v: number,
        oldMin: number,
        oldMax: number,
        newMin: number,
        newMax: number,
    ) => {
        return ((v - oldMin) / (oldMax - oldMin)) * (newMax - newMin) + newMin;
    },
};
