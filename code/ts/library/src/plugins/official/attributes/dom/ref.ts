// Authors: Delaney Gillilan
// Icon: mdi:cursor-pointer
// Slug: Create a reference to an element
// Description: This attribute creates a reference to an element that can be used in other expressions.

import { AttributePlugin } from "../../../../engine";
import { PLUGIN_ATTRIBUTE } from "../../../../engine/client_only_consts";

// Sets the value of the element
export const Ref: AttributePlugin = {
    pluginType: PLUGIN_ATTRIBUTE,
    name: "ref",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,
    bypassExpressionFunctionCreation: () => true,
    onLoad: (ctx) => {
        const { expressionFn, el, upsertSignal, removeSignals } = ctx;
        const signalPath = expressionFn(ctx);
        upsertSignal(signalPath, el);

        return () => {
            removeSignals(signalPath);
        };
    },
};
