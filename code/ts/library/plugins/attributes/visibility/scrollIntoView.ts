import { AttributeContext, AttributePlugin } from "library/engine";
import { scrollIntoView } from "library/utils/dom";

// Scrolls the element into view
export const ScrollIntoViewAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "scrollIntoView",
    mustHaveEmptyKey: true,
    mustHaveEmptyExpression: true,
    allowedModifiers: new Set([
        "smooth",
        "instant",
        "auto",
        "hstart",
        "hcenter",
        "hend",
        "hnearest",
        "vstart",
        "vcenter",
        "vend",
        "vnearest",
        "focus",
    ]),

    onLoad: ({ el, modifiers, rawKey }: AttributeContext) => {
        if (!el.tabIndex) el.setAttribute("tabindex", "0");
        const opts: ScrollIntoViewOptions = {
            behavior: "smooth",
            block: "center",
            inline: "center",
        };
        if (modifiers.has("smooth")) opts.behavior = "smooth";
        if (modifiers.has("instant")) opts.behavior = "instant";
        if (modifiers.has("auto")) opts.behavior = "auto";
        if (modifiers.has("hstart")) opts.inline = "start";
        if (modifiers.has("hcenter")) opts.inline = "center";
        if (modifiers.has("hend")) opts.inline = "end";
        if (modifiers.has("hnearest")) opts.inline = "nearest";
        if (modifiers.has("vstart")) opts.block = "start";
        if (modifiers.has("vcenter")) opts.block = "center";
        if (modifiers.has("vend")) opts.block = "end";
        if (modifiers.has("vnearest")) opts.block = "nearest";

        scrollIntoView(el, opts, modifiers.has("focus"));
        delete el.dataset[rawKey];
        return () => {};
    },
};
