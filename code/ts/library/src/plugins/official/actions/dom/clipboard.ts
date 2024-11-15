// Authors: Delaney Gillilan
// Icon: mdi:clipboard
// Slug: Copy text to the clipboard
// Description: This action copies text to the clipboard using the Clipboard API.

import { ActionPlugin } from "../../../engine";

export const Clipboard: ActionPlugin = {
    pluginType: "action",
    name: "clipboard",
    method: (_, text) => {
        if (!navigator.clipboard) {
            throw new Error("Clipboard API not available");
        }
        navigator.clipboard.writeText(text);
    },
};
