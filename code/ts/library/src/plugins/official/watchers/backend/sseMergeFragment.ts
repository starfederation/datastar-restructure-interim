// Authors: Delaney Gillilan
// Icon: material-symbols:cloud-download
// Slug: Use Server-Sent Events to fetch data from a server using the Datastar SDK interface
// Description: Remember, SSE is just a regular SSE request but with the ability to send 0-inf messages to the client.

import { InitContext, WatcherPlugin } from "../../../../engine";
import {
    docWithViewTransitionAPI,
    supportsViewTransitions,
} from "../../../../utils/view-transitions";
import { idiomorph } from "../../../../vendored/idiomorph";
import {
    datastarSSEEventWatcher,
    DEFAULT_SETTLE_DURATION_RAW,
    SETTLING_CLASS,
    SWAPPING_CLASS,
} from "./sseShared";

const DEFAULT_MERGE_MODE: FragmentMergeMode = "morph";
const DEFAULT_USE_VIEW_TRANSITION = "false";

const FragmentMergeModes = {
    MorphElement: "morph",
    InnerElement: "inner",
    OuterElement: "outer",
    PrependElement: "prepend",
    AppendElement: "append",
    BeforeElement: "before",
    AfterElement: "after",
    UpsertAttributes: "upsertAttributes",
} as const;
export type FragmentMergeMode =
    (typeof FragmentMergeModes)[keyof typeof FragmentMergeModes];

const name = "mergeFragments";
export const MergeFragments: WatcherPlugin = {
    pluginType: "effect",
    name: name,
    onGlobalInit: async (ctx) => {
        const fragmentContainer = document.createElement("template");
        datastarSSEEventWatcher(ctx, name, ({
            fragment = "<div></div>",
            selector = "",
            mergeMode = DEFAULT_MERGE_MODE,
            settleDuration: settleDurationRaw = DEFAULT_SETTLE_DURATION_RAW,
            useViewTransition: useViewTransitionRaw =
                DEFAULT_USE_VIEW_TRANSITION,
        }) => {
            const settleDuration = parseInt(settleDurationRaw);
            const useViewTransition = useViewTransitionRaw === "true";

            fragmentContainer.innerHTML = fragment.trim();
            const fragments = [...fragmentContainer.content.children];
            fragments.forEach((fragment) => {
                if (!(fragment instanceof Element)) {
                    throw new Error(`No fragment found`);
                }

                const selectorOrID = selector ||
                    `#${fragment.getAttribute("id")}`;
                const targets = document.querySelectorAll(selectorOrID) ||
                    [];
                const allTargets = [...targets];
                if (!allTargets.length) {
                    throw new Error(`No targets found for ${selector}`);
                }

                if (supportsViewTransitions && useViewTransition) {
                    docWithViewTransitionAPI.startViewTransition(() =>
                        applyToTargets(
                            ctx,
                            mergeMode,
                            settleDuration,
                            fragment,
                            allTargets,
                        )
                    );
                } else {
                    applyToTargets(
                        ctx,
                        mergeMode,
                        settleDuration,
                        fragment,
                        allTargets,
                    );
                }
            });
        });
    },
};

function applyToTargets(
    ctx: InitContext,
    mergeMode: string,
    settleDuration: number,
    fragment: Element,
    capturedTargets: Element[],
) {
    for (const initialTarget of capturedTargets) {
        initialTarget.classList.add(SWAPPING_CLASS);
        const originalHTML = initialTarget.outerHTML;
        let modifiedTarget = initialTarget;
        switch (mergeMode) {
            case FragmentMergeModes.MorphElement:
                const result = idiomorph(
                    modifiedTarget,
                    fragment,
                    {
                        callbacks: {
                            beforeNodeRemoved: (
                                oldNode: Element,
                                _: Element,
                            ) => {
                                ctx.cleanupElementRemovals(
                                    oldNode,
                                );
                                return true;
                            },
                        },
                    },
                );
                if (!result?.length) {
                    throw new Error(`No morph result `);
                }
                const first = result[0] as Element;
                modifiedTarget = first;
                break;
            case FragmentMergeModes.InnerElement:
                // Replace the contents of the target element with the response
                modifiedTarget.innerHTML = fragment.innerHTML;
                break;
            case FragmentMergeModes.OuterElement:
                // Replace the entire target element with the response
                modifiedTarget.replaceWith(fragment);
                break;
            case FragmentMergeModes.PrependElement:
                // Insert the response before the first child of the target element
                modifiedTarget.prepend(fragment);
                break;
            case FragmentMergeModes.AppendElement:
                // Insert the response after the last child of the target element
                modifiedTarget.append(fragment);
                break;
            case FragmentMergeModes.BeforeElement:
                // Insert the response before the target element
                modifiedTarget.before(fragment);
                break;
            case FragmentMergeModes.AfterElement:
                // Insert the response after the target element
                modifiedTarget.after(fragment);
                break;
            case FragmentMergeModes.UpsertAttributes:
                // Upsert the attributes of the target element
                fragment.getAttributeNames().forEach(
                    (attrName) => {
                        const value = fragment.getAttribute(
                            attrName,
                        )!;
                        modifiedTarget.setAttribute(
                            attrName,
                            value,
                        );
                    },
                );
                break;
            default:
                throw new Error(
                    `Unknown merge type: ${mergeMode}`,
                );
        }
        ctx.cleanupElementRemovals(modifiedTarget);
        modifiedTarget.classList.add(SWAPPING_CLASS);

        ctx.applyPlugins(document.body);

        setTimeout(() => {
            initialTarget.classList.remove(SWAPPING_CLASS);
            modifiedTarget.classList.remove(SWAPPING_CLASS);
        }, settleDuration);

        const revisedHTML = modifiedTarget.outerHTML;

        if (originalHTML !== revisedHTML) {
            modifiedTarget.classList.add(SETTLING_CLASS);
            setTimeout(() => {
                modifiedTarget.classList.remove(
                    SETTLING_CLASS,
                );
            }, settleDuration);
        }
    }
}
