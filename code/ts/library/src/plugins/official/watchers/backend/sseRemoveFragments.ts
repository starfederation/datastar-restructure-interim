// Authors: Delaney Gillilan
// Icon: material-symbols:settings-input-antenna
// Slug: Merge fine grain signals store data from a server using the Datastar SDK interface
// Description: Merge store data from a server using the Datastar SDK interface

import { WatcherPlugin } from "../../../../engine";
import {
    docWithViewTransitionAPI,
    supportsViewTransitions,
} from "../../../../utils/view-transitions";
import {
    datastarSSEEventWatcher,
    DEFAULT_SETTLE_DURATION_RAW,
    SWAPPING_CLASS,
} from "./sseShared";

const name = "removeFragments";
export const RemoveFragments: WatcherPlugin = {
    pluginType: "effect",
    name,
    onGlobalInit: async (ctx) => {
        datastarSSEEventWatcher(ctx, name, ({
            selector,
            settleDuration: settleDurationRaw = DEFAULT_SETTLE_DURATION_RAW,
            useViewTransition: useViewTransitionRaw = "false",
        }) => {
            if (!!!selector.length) {
                throw new Error(
                    "No selector provided for remove-fragments",
                );
            }

            const settleDuration = parseInt(settleDurationRaw);
            const useViewTransition = useViewTransitionRaw === "true";
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
