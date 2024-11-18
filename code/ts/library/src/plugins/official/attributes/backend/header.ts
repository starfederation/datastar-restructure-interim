// Authors: Delaney Gillilan
// Icon: fluent:missing-metadata-24-filled
// Slug: Add headers to fetch requests
// Description: This plugin allows you to add headers to fetch requests.  Once you add this attribute the headers will be added to the fetch request.

import { AttributePlugin } from "../../../../engine";
import { PLUGIN_ATTRIBUTE } from "../../../../engine/client_only_consts";

export const Header: AttributePlugin = {
    pluginType: PLUGIN_ATTRIBUTE,
    name: "header",
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
