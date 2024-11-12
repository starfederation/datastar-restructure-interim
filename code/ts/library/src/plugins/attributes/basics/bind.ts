import { AttributePlugin } from "../../../engine";
import { kebabize } from "../../../utils/text";

export const Bind: AttributePlugin = {
    pluginType: "attribute",
    prefix: "bind",
    mustNotEmptyKey: true,
    mustNotEmptyExpression: true,

    onLoad: (ctx) => {
        return ctx.reactivity.effect(async () => {
            const key = kebabize(ctx.key);
            const value = ctx.expressionFn(ctx);
            let v: string;
            if (typeof value === "string") {
                v = value;
            } else {
                v = JSON.stringify(value);
            }
            if (!v || v === "false" || v === "null" || v === "undefined") {
                ctx.el.removeAttribute(key);
            } else {
                ctx.el.setAttribute(key, v);
            }
        });
    },
};
