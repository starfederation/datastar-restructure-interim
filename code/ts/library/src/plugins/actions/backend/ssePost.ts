import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const PostSSE: ActionPlugin = {
    pluginType: "action",
    name: "post",
    method: fetcherActionMethod("POST"),
};
