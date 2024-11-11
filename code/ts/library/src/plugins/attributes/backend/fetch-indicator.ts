import { AttributePlugin, DATASTAR } from "../../../engine";

export const INDICATOR_CLASS = `${DATASTAR}-indicator`;
export const INDICATOR_LOADING_CLASS = `${INDICATOR_CLASS}-loading`;

// Sets the fetch indicator selector
export const FetchIndicatorPlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "fetchIndicator",
    mustHaveEmptyKey: true,
    mustNotEmptyExpression: true,
    onGlobalInit: () => {
        const style = document.createElement("style");
        style.innerHTML = `
  .${INDICATOR_CLASS}{
   opacity:0;
   transition: opacity 300ms ease-out;
  }
  .${INDICATOR_LOADING_CLASS} {
   opacity:1;
   transition: opacity 300ms ease-in;
  }
  `;
        document.head.appendChild(style);
    },
    onLoad: (ctx) => {
        return ctx.reactivity.effect(() => {
            ctx.upsertIfMissingFromStore(
                "_dsPlugins.fetch.indicatorElements",
                {},
            );
            ctx.upsertIfMissingFromStore(
                "_dsPlugins.fetch.indicatorsVisible",
                [],
            );
            const c = ctx.reactivity.computed(() => `${ctx.expressionFn(ctx)}`);
            const s = ctx.store();

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
