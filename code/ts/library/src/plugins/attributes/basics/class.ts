import { AttributePlugin } from "library/src/engine";

export const ClassAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "class",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,

    onLoad: (ctx) => {
        return ctx.reactivity.effect(() => {
            const classes: Object = ctx.expressionFn(ctx);
            for (const [k, v] of Object.entries(classes)) {
                if (v) {
                    ctx.el.classList.add(k);
                } else {
                    ctx.el.classList.remove(k);
                }
            }

            return () => {
                ctx.el.classList.remove(...Object.keys(classes));
            };
        });
    },
};
