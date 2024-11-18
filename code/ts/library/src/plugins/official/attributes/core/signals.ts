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
const name = "signals";
export const Signals: AttributePlugin = {
    pluginType: PLUGIN_ATTRIBUTE,
    name,
    removeNewLines: true,
    preprocessors: {
        pre: [
            {
                pluginType: PLUGIN_PREPROCESSOR,
                name,
                regexp: /(?<whole>.+)/g,
                replacer: (groups: RegexpGroups) => {
                    const { whole } = groups;
                    return `Object.assign({...ctx.signals()}, ${whole})`;
                },
            },
        ],
    },
    allowedModifiers: new Set(["ifmissing"]),
    onLoad: (ctx: AttributeContext) => {
        const possibleMergeSignals = ctx.expressionFn(ctx);
        const actualMergeSignals = storeFromPossibleContents(
            ctx.signals(),
            possibleMergeSignals,
            ctx.modifiers.has("ifmissing"),
        );
        ctx.mergeSignals(actualMergeSignals);

        delete ctx.el.dataset[ctx.rawKey];
    },
};
