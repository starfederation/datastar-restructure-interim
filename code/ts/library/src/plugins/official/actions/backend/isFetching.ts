// Authors: Delaney Gillilan
// Icon: carbon:fetch-upload
// Slug: Check if a fetch request is currently in progress
// Description: This action checks if a fetch request is currently in progress. This is useful for showing loading indicators or disabling buttons while a fetch request is in progress.

import { ActionPlugin } from "../../../../engine";
import { PLUGIN_ACTION } from "../../../../engine/client_only_consts";
import { IndicatorReference } from "./sseShared";

export const IsFetching: ActionPlugin = {
    pluginType: PLUGIN_ACTION,
    name: "isFetching",
    method: (ctx, selector: string) => {
        const indicators = [...document.querySelectorAll(selector)];
        const store = ctx.signals();
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
