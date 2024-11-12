import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const DeleteSSE: ActionPlugin = {
    pluginType: "action",
    name: "delete",
    method: fetcherActionMethod("DELETE"),
};
