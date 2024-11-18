// Authors: Delaney Gillilan
// Icon: material-symbols:home-storage
// Slug: Store signals into a singleton per page
// Description: This action stores signals into a singleton per page. This is useful for storing signals that are used across multiple components.

import {
    AttributeContext,
    AttributePlugin,
    RegexpGroups,
} from "../../../../engine";
import {
    PLUGIN_ATTRIBUTE,
    PLUGIN_PREPROCESSOR,
} from "../../../../engine/client_only_consts";
import { storeFromPossibleContents } from "../../../../utils/signals";

// Setup the global store
export const Store: AttributePlugin = {
    pluginType: PLUGIN_ATTRIBUTE,
    name: "store",
    removeNewLines: true,
    preprocessors: {
        pre: [
            {
                pluginType: PLUGIN_PREPROCESSOR,
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
