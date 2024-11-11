import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const PatchSSEActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "patch",
    method: fetcherActionMethod("PATCH"),
};
