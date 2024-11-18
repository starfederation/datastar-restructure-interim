// Authors: Delaney Gillilan
// Icon: tabler:typography
// Slug: Set the text content of an element
// Description: This attribute sets the text content of an element to the result of the expression.

import { AttributePlugin } from "../../../../engine";
import { PLUGIN_ATTRIBUTE } from "../../../../engine/client_only_consts";

export const Text: AttributePlugin = {
    pluginType: PLUGIN_ATTRIBUTE,
    name: "text",
    mustHaveEmptyKey: true,
    onLoad: (ctx) => {
        const { el, expressionFn } = ctx;
        if (!(el instanceof HTMLElement)) {
            throw new Error("Element is not HTMLElement");
        }
        return ctx.reactivity.effect(() => {
            const res = expressionFn(ctx);
            el.textContent = `${res}`;
        });
    },
};
