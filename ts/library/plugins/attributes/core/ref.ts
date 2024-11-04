import { AttributePlugin } from "library/engine";
import { elemToSelector } from "library/utils/dom";

// Sets the value of the element
export const RefPlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "ref",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,
    bypassExpressionFunctionCreation: () => true,
    onLoad: (ctx) => {
        ctx.upsertIfMissingFromStore("_dsPlugins.refs", {});
        const { el, expression } = ctx;
        const s = ctx.store();

        const revised = {
            _dsPlugins: {
                refs: {
                    ...s._dsPlugins.refs.value,
                    [expression]: elemToSelector(el),
                },
            },
        };
        ctx.mergeStore(revised);

        return () => {
            const s = ctx.store();
            const revised = { ...s._dsPlugins.refs.value };
            delete revised[expression];
            s._dsPlugins.refs = revised;
        };
    },
};
