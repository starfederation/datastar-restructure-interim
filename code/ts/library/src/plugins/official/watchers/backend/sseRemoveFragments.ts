// Authors: Delaney Gillilan
// Icon: material-symbols:settings-input-antenna
// Slug: Merge fine grain signals store data from a server using the Datastar SDK interface
// Description: Merge store data from a server using the Datastar SDK interface

import {
    DefaultFragmentsUseViewTransitions,
    DefaultSettleDurationMs,
    EventTypes,
    WatcherPlugin,
} from "../../../../engine";
import { PLUGIN_WATCHER } from "../../../../engine/client_only_consts";
import { isBoolString } from "../../../../utils/text";
import {
    docWithViewTransitionAPI,
    supportsViewTransitions,
} from "../../../../utils/view-transitions";
import { datastarSSEEventWatcher, SWAPPING_CLASS } from "./sseShared";

export const RemoveFragments: WatcherPlugin = {
    pluginType: PLUGIN_WATCHER,
    name: EventTypes.RemoveFragments,
    onGlobalInit: async () => {
        datastarSSEEventWatcher(EventTypes.RemoveFragments, ({
            selector,
            settleDuration: settleDurationRaw = `${DefaultSettleDurationMs}`,
            useViewTransition: useViewTransitionRaw =
                `${DefaultFragmentsUseViewTransitions}`,
        }) => {
            if (!!!selector.length) {
                throw new Error(
                    "No selector provided for remove-fragments",
                );
            }

            const settleDuration = parseInt(settleDurationRaw);
            const useViewTransition = isBoolString(useViewTransitionRaw);
            const removeTargets = document.querySelectorAll(selector);

            const applyToTargets = () => {
                for (const target of removeTargets) {
                    target.classList.add(SWAPPING_CLASS);
                }

                setTimeout(() => {
                    for (const target of removeTargets) {
                        target.remove();
                    }
                }, settleDuration);
            };

            if (supportsViewTransitions && useViewTransition) {
                docWithViewTransitionAPI.startViewTransition(() =>
                    applyToTargets()
                );
            } else {
                applyToTargets();
            }
        });
    },
};
