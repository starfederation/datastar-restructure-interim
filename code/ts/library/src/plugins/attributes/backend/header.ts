import { AttributePlugin } from "library/src/engine";

export const HeaderPlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "header",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,

    onLoad: (ctx) => {
        ctx.upsertIfMissingFromStore("_dsPlugins.fetch.headers", {});

        const headers: Object = ctx.expressionFn(ctx);
        for (const [key, value] of Object.entries(headers)) {
            ctx.store()._dsPlugins.fetch.headers[key] = value;
        }

        return () => {
            for (const key of Object.keys(headers)) {
                delete ctx.store()._dsPlugins.fetch.headers[key];
            }
        };
    },
};
