// Authors: Delaney Gillilan
// Icon: mdi:cursor-pointer
// Slug: Reference an element
// Description: This action references an element that can be used in other expressions.

import { ActionPlugin } from "../../../../engine";
import { PLUGIN_ACTION } from "../../../../engine/client_only_consts";

export const RefAction: ActionPlugin = {
    pluginType: PLUGIN_ACTION,
    name: "ref",
    method: (ctx, text) => {
        const selectorSignal = ctx.store()?._dsPlugins?.refs?.[text];
        if (!selectorSignal) {
            throw new Error(`Reference '${text}' not found`);
        }
        const selector = selectorSignal?.value;
        if (!selector) {
            throw new Error(`Selector '${selector}' not found`);
        }

        const el = document.querySelector(selector);
        if (!el) {
            throw new Error(`Elements for selector '${selector}' not found`);
        }

        return el;
    },
};
