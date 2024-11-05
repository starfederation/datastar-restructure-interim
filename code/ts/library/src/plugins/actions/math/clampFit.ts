import { ActionPlugin, AttributeContext } from "library/src/engine";

export const ClampFitActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "clampFit",
    method: (
        _: AttributeContext,
        v: number,
        oldMin: number,
        oldMax: number,
        newMin: number,
        newMax: number,
    ) => {
        return Math.max(
            newMin,
            Math.min(
                newMax,
                ((v - oldMin) / (oldMax - oldMin)) * (newMax - newMin) + newMin,
            ),
        );
    },
};
