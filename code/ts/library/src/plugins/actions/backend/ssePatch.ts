import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const PatchSSE: ActionPlugin = {
    pluginType: "action",
    name: "patch",
    method: fetcherActionMethod("PATCH"),
};
