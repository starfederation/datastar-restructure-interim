// Authors: Delaney Gillilan
// Icon: material-symbols:find-replace
// Slug: Use a PUT request to fetch data from a server using Server-Sent Events matching the Datastar SDK interface
// Description: Remember, SSE is just a regular SSE request but with the ability to send 0-inf messages to the client.

import { ActionPlugin } from "../../../../engine";
import { PLUGIN_ACTION } from "../../../../engine/client_only_consts";
import { PUT, sendSSERequest } from "./sseShared";

export const PutSSE: ActionPlugin = {
    pluginType: PLUGIN_ACTION,
    name: PUT,
    method: sendSSERequest(PUT),
};
