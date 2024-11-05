import { ActionPlugin, AttributeContext } from "library/src/engine";
import { scrollIntoView } from "library/src/utils/dom";

export const ScrollActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "scroll",
    method: async (
        _: AttributeContext,
        selector: string,
        opts: {
            behavior: "smooth" | "instant" | "auto"; // smooth is default
            vertical: "start" | "center" | "end" | "nearest"; // center is default
            horizontal: "start" | "center" | "end" | "nearest"; // center is default
            shouldFocus: boolean;
        },
    ) => {
        const allOpts = Object.assign(
            {
                behavior: "smooth",
                vertical: "center",
                horizontal: "center",
                shouldFocus: true,
            },
            opts,
        );
        const el = document.querySelector(selector);
        scrollIntoView(el as HTMLElement, allOpts);
    },
};
