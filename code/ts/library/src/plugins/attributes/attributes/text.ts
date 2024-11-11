import { AttributePlugin } from "../../../engine";

// Sets the textContent of the element
export const TextAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "text",
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
