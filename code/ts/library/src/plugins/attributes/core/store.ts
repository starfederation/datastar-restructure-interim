import {
    AttributeContext,
    AttributePlugin,
    RegexpGroups,
} from "library/src/engine";
import { storeFromPossibleContents } from "library/src/utils/signals";

// Setup the global store
export const StoreAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "store",
    removeNewLines: true,
    preprocessors: {
        pre: [
            {
                pluginType: "preprocessor",
                name: "store",
                regexp: /(?<whole>.+)/g,
                replacer: (groups: RegexpGroups) => {
                    const { whole } = groups;
                    return `Object.assign({...ctx.store()}, ${whole})`;
                },
            },
        ],
    },
    allowedModifiers: new Set(["ifmissing"]),
    onLoad: (ctx: AttributeContext) => {
        const possibleMergeStore = ctx.expressionFn(ctx);
        const actualMergeStore = storeFromPossibleContents(
            ctx.store(),
            possibleMergeStore,
            ctx.modifiers.has("ifmissing"),
        );
        ctx.mergeStore(actualMergeStore);

        delete ctx.el.dataset[ctx.rawKey];
    },
};
