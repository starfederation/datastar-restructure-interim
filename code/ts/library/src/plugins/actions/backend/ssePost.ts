import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const PostSSEActionPlugin: ActionPlugin = {
    pluginType: "action",
    name: "post",
    method: fetcherActionMethod("POST"),
};
