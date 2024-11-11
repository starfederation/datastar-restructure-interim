import { AttributePlugin } from "../../../engine";

export const ComputedPlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "computed",
    mustNotEmptyKey: true,
    onLoad: (ctx) => {
        const store = ctx.store();
        store[ctx.key] = ctx.reactivity.computed(() => {
            return ctx.expressionFn(ctx);
        });

        return () => {
            const store = ctx.store();
            delete store[ctx.key];
        };
    },
};
