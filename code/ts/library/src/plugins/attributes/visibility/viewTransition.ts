import { AttributePlugin } from "library/src/engine";
import { supportsViewTransitions } from "library/src/utils/view-transitions";

// Setup view transition api
export const ViewTransitionAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "viewTransition",
    onGlobalInit() {
        let hasViewTransitionMeta = false;
        document.head.childNodes.forEach((node) => {
            if (
                node instanceof HTMLMetaElement &&
                node.name === "view-transition"
            ) {
                hasViewTransitionMeta = true;
            }
        });

        if (!hasViewTransitionMeta) {
            const meta = document.createElement("meta");
            meta.name = "view-transition";
            meta.content = "same-origin";
            document.head.appendChild(meta);
        }
    },
    onLoad: (ctx) => {
        if (!supportsViewTransitions) {
            console.error("Browser does not support view transitions");
            return;
        }

        return ctx.reactivity.effect(() => {
            const { el, expressionFn } = ctx;
            let name = expressionFn(ctx);
            if (!name) return;

            const elVTASTyle = el.style as unknown as CSSStyleDeclaration;
            elVTASTyle.viewTransitionName = name;
        });
    },
};
