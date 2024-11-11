import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const PutSSEActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "put",
    method: fetcherActionMethod("PUT"),
};
