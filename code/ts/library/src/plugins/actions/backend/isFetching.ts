// Authors: Delaney Gillilan
// Icon: carbon:fetch-upload
// Slug: Check if a fetch request is currently in progress
// Description: This action checks if a fetch request is currently in progress. This is useful for showing loading indicators or disabling buttons while a fetch request is in progress.

import { ActionPlugin } from "../../../engine";
export type IndicatorReference = { el: HTMLElement; count: number };

export const IsFetching: ActionPlugin = {
    pluginType: "action",
    name: "isFetching",
    method: (ctx, selector: string) => {
        const indicators = [...document.querySelectorAll(selector)];
        const store = ctx.store();
        const indicatorsVisible: IndicatorReference[] =
            store?._dsPlugins?.fetch.indicatorsVisible?.value || [];
        if (!!!indicators.length) return false;

        return indicators.some((indicator) => {
            return indicatorsVisible
                .filter((val) => !!val)
                .some((indicatorVisible) => {
                    return indicatorVisible.el.isSameNode(indicator) &&
                        indicatorVisible.count > 0;
                });
        });
    },
};
