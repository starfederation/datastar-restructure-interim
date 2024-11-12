import { ActionPlugin, AttributeContext } from "../../../engine";

export const FitInt: ActionPlugin = {
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
