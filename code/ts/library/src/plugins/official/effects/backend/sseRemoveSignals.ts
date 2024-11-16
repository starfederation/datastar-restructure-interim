// Authors: Delaney Gillilan
// Icon: material-symbols:settings-input-antenna
// Slug: Merge fine grain signals store data from a server using the Datastar SDK interface
// Description: Merge store data from a server using the Datastar SDK interface

import { EffectPlugin } from "../../../../engine";
import { datastarSSEEventWatcher } from "./sseShared";

const name = "removeSignals";
export const RemoveSignals: EffectPlugin = {
    pluginType: "effect",
    name,
    onGlobalInit: async (ctx) => {
        datastarSSEEventWatcher(ctx, name, ({ paths: pathsRaw = "" }) => {
            // replace all whitespace with a single space
            pathsRaw = pathsRaw.replaceAll(/\s+/g, " ");
            if (!!!pathsRaw?.length) {
                throw new Error(
                    "No paths provided for remove-signals",
                );
            }

            const paths = pathsRaw.split(" ");
            ctx.removeFromStore(...paths);
        });
    },
};
