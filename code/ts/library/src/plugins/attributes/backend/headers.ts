// Authors: Delaney Gillilan
// Icon: material-symbols:scrollable-header-outline
// Slug: Add headers to a SSE request
// Description: Add headers to a Server-Sent Events request.  In general this shouldn't be necessary, but it can be useful for debugging or working with a server that requires headers or when interacting with a backend that requires headers.

import { AttributePlugin, RegexpGroups } from "../../../engine";

export const Headers: AttributePlugin = {
    pluginType: "attribute",
    prefix: "header",
    mustNotEmptyKey: true,
    mustNotEmptyExpression: true,
    preprocessors: {
        post: [
            {
                pluginType: "preprocessor",
                name: "header",
                regexp: /(?<whole>.+)/g,
                replacer: (groups: RegexpGroups) => {
                    const { whole } = groups;
                    return `'${whole}'`;
                },
            },
        ],
    },
    onLoad: (ctx) => {
        ctx.upsertIfMissingFromStore("_dsPlugins.fetch.headers", {});
        const key = ctx.key.replace(/([a-z](?=[A-Z]))/g, "$1-").toUpperCase();
        const value = ctx.expressionFn(ctx);
        ctx.store()._dsPlugins.fetch.headers[key] = value;
        return () => {
            delete ctx.store()._dsPlugins.fetch.headers[key];
        };
    },
};
