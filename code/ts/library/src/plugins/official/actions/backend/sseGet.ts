// Authors: Delaney Gillilan
// Icon: ic:baseline-get-app
// Slug: Use a GET request to fetch data from a server using Server-Sent Events matching the Datastar SDK interface
// Description: Remember, SSE is just a regular SSE request but with the ability to send 0-inf messages to the client.

import { ActionPlugin } from "../../../../engine";
import { PLUGIN_ACTION } from "../../../../engine/client_only_consts";
import { sendSSERequest } from "./sseShared";

export const GetSSE: ActionPlugin = {
    pluginType: PLUGIN_ACTION,
    name: "get",
    method: sendSSERequest("GET"),
};
