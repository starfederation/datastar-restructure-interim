import {
    AttributePlugin,
} from "library/src/engine";

export const ReplaceUrlAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "replaceUrl",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,

    onLoad: (ctx) => {
        const value = ctx.expressionFn(ctx);
        const baseUrl = window.location.href;
        const url = new URL(value, baseUrl).toString();

        window.history.replaceState({}, '', url);
    },
};
