import { AttributePlugin } from "../../../engine";
import { nodeHTMLorSVGElement } from "../../../utils/dom";

const PREPEND = "prepend";
const APPEND = "append";
const teleportParentErr = new Error(
    "Target element must have a parent if using prepend or append",
);
// Teleports the element to another element
export const TeleportAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "teleport",
    allowedModifiers: new Set([PREPEND, APPEND]),
    allowedTagRegexps: new Set(["template"]),
    bypassExpressionFunctionCreation: () => true,
    onLoad: (ctx) => {
        const { el, modifiers, expression } = ctx;
        if (!(el instanceof HTMLTemplateElement)) {
            throw new Error(`el must be a template element`);
        }

        const target = document.querySelector(expression);
        if (!target) {
            throw new Error(`Target element not found: ${expression}`);
        }

        if (!el.content) {
            throw new Error("Template element must have content");
        }

        const n = el.content.cloneNode(true);
        const nEl = nodeHTMLorSVGElement(n as Element);
        if (nEl?.firstElementChild) {
            throw new Error("Empty template");
        }

        if (modifiers.has(PREPEND)) {
            if (!target.parentNode) throw teleportParentErr;
            target.parentNode.insertBefore(n, target);
        } else if (modifiers.has(APPEND)) {
            if (!target.parentNode) throw teleportParentErr;
            target.parentNode.insertBefore(n, target.nextSibling);
        } else {
            target.appendChild(n);
        }
    },
};
