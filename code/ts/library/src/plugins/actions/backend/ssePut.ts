import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const PutSSE: ActionPlugin = {
    pluginType: "action",
    name: "put",
    method: fetcherActionMethod("PUT"),
};
