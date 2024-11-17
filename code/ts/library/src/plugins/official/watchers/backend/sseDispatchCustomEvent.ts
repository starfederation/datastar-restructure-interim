// Authors: Delaney Gillilan
// Icon: material-symbols:mail
// Slug: Dispatch a custom event to a specified elements from the backend
// Description: Dispatch a custom event to a specified elements from the backend

import { WatcherPlugin } from "../../../../engine";
import { datastarSSEEventWatcher } from "./sseShared";

const name = "dispatchCustomEvent";
export const DispatchCustomEvent: WatcherPlugin = {
    pluginType: "effect",
    name,
    onGlobalInit: async (ctx) => {
        datastarSSEEventWatcher(
            ctx,
            name,
            ({
                eventName,
                selector = "document",
                bubbles: bubblesRaw = "true",
                cancelable: cancelableRaw = "true",
                composed: composedRaw = "true",
                detailJson = "{}",
            }) => {
                if (!eventName.length) {
                    throw new Error(
                        "No event name provided for dispatch-custom-event",
                    );
                }
                const elements = selector === "document"
                    ? [document]
                    : document.querySelectorAll(selector);

                const composed = composedRaw.trim() === "true";
                const bubbles = bubblesRaw.trim() === "true";
                const cancelable = cancelableRaw.trim() === "true";
                const detail = JSON.parse(detailJson);

                const event = new CustomEvent(eventName, {
                    bubbles,
                    cancelable,
                    composed,
                    detail,
                });

                elements.forEach((element) => {
                    element.dispatchEvent(event);
                });
            },
        );
    },
};
