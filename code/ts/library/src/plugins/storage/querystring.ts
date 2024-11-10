import {
    AttributePlugin,
} from "library/src/engine";

export const QueryStringAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "querystring",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,

    onLoad: (ctx) => {
        // console.log("QueryStringAttributePlugin onLoad ", ctx);

        const keys = new Set<string>();
        let expression = ctx.expression;
        const parts = expression.split(" ");
        for (const part of parts) {
            keys.add(part);
        }

        if (keys.size > 0) {
            const url = new URL(window.location.href);
            const store = ctx.store();
            const queryStringParams: Record<string, any> = {};
            for (const key of keys) {
                const parts = key.split(".");
                let subQueryStringParams = queryStringParams;
                let subStore = store;
                for (let i = 0; i < parts.length - 1; i++) {
                    const part = parts[i];
                    if (!subQueryStringParams[part]) {
                        subQueryStringParams[part] = {};
                    }
                    subQueryStringParams = subQueryStringParams[part];
                }

                const lastPart = parts[parts.length - 1];
                subQueryStringParams[lastPart] = subStore[lastPart];

                const marshalledStore = JSON.stringify(queryStringParams[key]);
                url.searchParams.set(key, marshalledStore);
            }

            window.history.replaceState({}, '', url);
        }
    },
};
