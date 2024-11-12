// Authors: Delaney Gillilan
// Icon: ic:baseline-get-app
// Description: Use a GET request to fetch data from a server using Server-Sent Events matching the Datastar SDK interface

import { ActionPlugin } from "../../../engine";
import { fetcherActionMethod } from "./sse";

export const GetSSE: ActionPlugin = {
    pluginType: "action",
    name: "get",
    method: fetcherActionMethod("GET"),
};
