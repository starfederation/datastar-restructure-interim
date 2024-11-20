// Authors: Delaney Gillilan
// Icon: game-icons:teleport
// Slug: Teleport an element to another element
// Description: This attribute teleports an element to another element in the DOM.

import { AttributePlugin } from "../../../../engine";
import { PLUGIN_ATTRIBUTE } from "../../../../engine/client_only_consts";
import { FragmentMergeModes } from "../../../../engine/consts";
import { ERR_BAD_ARGS } from "../../../../engine/errors";
import { nodeHTMLorSVGElement } from "../../../../utils/dom";

export const Teleport: AttributePlugin = {
    pluginType: PLUGIN_ATTRIBUTE,
    name: "teleport",
    allowedModifiers: new Set([
        FragmentMergeModes.Prepend,
        FragmentMergeModes.Append,
    ]),
    allowedTagRegexps: new Set(["template"]),
    bypassExpressionFunctionCreation: () => true,
    onLoad: (ctx) => {
        const { el, modifiers, expression } = ctx;
        if (!(el instanceof HTMLTemplateElement)) {
            // Element is not a template element
            throw ERR_BAD_ARGS;
        }

        const target = document.querySelector(expression);
        if (!target) {
            // Target element not found
            throw ERR_BAD_ARGS;
        }

        if (!el.content) {
            // Template element must have content
            throw ERR_BAD_ARGS;
        }

        const n = el.content.cloneNode(true);
        const nEl = nodeHTMLorSVGElement(n as Element);
        if (nEl?.firstElementChild) {
            // Empty template
            throw ERR_BAD_ARGS;
        }

        if (modifiers.has(FragmentMergeModes.Prepend)) {
            // Target element must have a parent if using prepend or append
            if (!target.parentNode) throw ERR_BAD_ARGS;
            target.parentNode.insertBefore(n, target);
        } else if (modifiers.has(FragmentMergeModes.Append)) {
            // Target element must have a parent if using prepend or append
            if (!target.parentNode) throw ERR_BAD_ARGS;
            target.parentNode.insertBefore(n, target.nextSibling);
        } else {
            target.appendChild(n);
        }
    },
};
