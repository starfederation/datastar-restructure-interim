import { ActionPlugin } from "../../../engine";

export const ClipboardActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "clipboard",
    method: (_, text) => {
        if (!navigator.clipboard) {
            throw new Error("Clipboard API not available");
        }
        navigator.clipboard.writeText(text);
    },
};
