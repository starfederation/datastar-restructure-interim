import { ActionPlugin, AttributeContext } from "library/engine";

export const FitIntActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "fitInt",
    method: (
        _: AttributeContext,
        v: number,
        oldMin: number,
        oldMax: number,
        newMin: number,
        newMax: number,
    ) => {
        return Math.round(
            ((v - oldMin) / (oldMax - oldMin)) * (newMax - newMin) + newMin,
        );
    },
};
