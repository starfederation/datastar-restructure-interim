// Authors: Delaney Gillilan
// Icon: material-symbols:network-wifi
// Slug: Use a selector to show a loading indicator when fetching data from the server
// Description: This plugin allows you to use a selector to show a loading indicator when fetching data from the server.  Once you add this attribute the indicator will be hidden by default and shown when the fetch is in progress.

import { AttributePlugin, DATASTAR } from "../../../../engine";
import { PLUGIN_ATTRIBUTE } from "../../../../engine/client_only_consts";

export const INDICATOR_CLASS = `${DATASTAR}-indicator`;
export const INDICATOR_LOADING_CLASS = `${INDICATOR_CLASS}-loading`;

// Sets the fetch indicator selector
export const FetchIndicator: AttributePlugin = {
    pluginType: PLUGIN_ATTRIBUTE,
    name: "fetchIndicator",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,
    onGlobalInit: () => {
        const style = document.createElement("style");
        style.innerHTML = `.${INDICATOR_CLASS}{
opacity:0;
transition: opacity 300ms ease-out;
}
.${INDICATOR_LOADING_CLASS} {
opacity:1;
transition: opacity 300ms ease-in;
}`;
        document.head.appendChild(style);
    },
    onLoad: (ctx) => {
        return ctx.reactivity.effect(() => {
            ctx.upsertIfMissingSignals(
                "_dsPlugins.fetch.indicatorElements",
                {},
            );
            ctx.upsertIfMissingSignals(
                "_dsPlugins.fetch.indicatorsVisible",
                [],
            );
            const c = ctx.reactivity.computed(() => `${ctx.expressionFn(ctx)}`);
            const s = ctx.signals();

            const indicators = document.querySelectorAll(c.value);
            if (indicators.length === 0) {
                throw new Error(`No indicator found`);
            }
            indicators.forEach((indicator) => {
                indicator.classList.add(INDICATOR_CLASS);
            });

            s._dsPlugins.fetch.indicatorElements[ctx.el.id] = ctx.reactivity
                .signal(indicators);

            return () => {
                delete s._dsPlugins.fetch.indicatorElements[ctx.el.id];
            };
        });
    },
};
