import {
    AttributePlugin,
} from "library/src/engine";

export const QueryStringAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "querystring",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,

    onLoad: (ctx) => {
        const url = new URL(window.location.href);
        const { expression, store } = ctx
        const keys = expression.split(" ");

        for (const key of keys) {
            const parts = key.split(".");
            const params: Record<string, any> = {};
            let subParams = params;
            let subStore = store;
            for (let i = 0; i < parts.length - 1; i++) {
                const part = parts[i];
                if (!subParams[part]) {
                    subParams[part] = {};
                }
                subParams = subParams[part];
            }

            const lastPart = parts[parts.length - 1];
            subParams[lastPart] = subStore[lastPart];

            const queryString = new URLSearchParams(params).toString();
            url.searchParams.set(key, queryString);
        }

        window.history.replaceState({}, '', url);
    },
};
