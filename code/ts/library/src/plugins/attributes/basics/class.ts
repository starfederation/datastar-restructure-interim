import { AttributePlugin } from "library/src/engine";

export const ClassAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "class",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,

    onLoad: (ctx) => {
        const classes: Object = ctx.expressionFn(ctx);

        return ctx.reactivity.effect(() => {
            for (const [key, value] of Object.entries(classes)) {
                if (value) {
                    ctx.el.classList.add(key);
                } else {
                    ctx.el.classList.remove(key);
                }
            }

            return () => {
                ctx.el.classList.remove(...Object.keys(classes));
            };
        });
    },
};
