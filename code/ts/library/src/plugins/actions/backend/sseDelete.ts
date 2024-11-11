import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const DeleteSSEActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "delete",
    method: fetcherActionMethod("DELETE"),
};
