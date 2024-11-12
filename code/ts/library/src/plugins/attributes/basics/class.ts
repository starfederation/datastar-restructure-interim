import { AttributePlugin } from "../../../engine";

export const Class: AttributePlugin = {
    pluginType: "attribute",
    prefix: "class",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,

    onLoad: (ctx) => {
        const classes: Object = ctx.expressionFn(ctx);

        return ctx.reactivity.effect(() => {
            for (const [k, v] of Object.entries(classes)) {
                const clss = k.split(" ");
                if (v) {
                    ctx.el.classList.add(...clss);
                } else {
                    ctx.el.classList.remove(...clss);
                }
            }
        });
    },
};
