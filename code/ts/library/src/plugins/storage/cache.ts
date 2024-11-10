import {
    AttributePlugin,
    DATASTAR,
    DATASTAR_EVENT,
    DatastarEvent,
} from "library/src/engine";
import { remoteSignals } from "library/src/utils/signals";

export const CacheStoreAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "cache",
    allowedModifiers: new Set(["local", "session", "remote"]),
    onLoad: (ctx) => {
        // console.log("CacheStoreAttributePlugin onLoad ", ctx);

        const key = ctx.key || DATASTAR;

        const keys = new Set<string>();
        let expression = ctx.expression;
        if (expression.trim() !== "") {
            debugger;
            expression = ctx.expressionFn(ctx);

            if (typeof expression === "string") {
                const parts = expression.split(" ");
                for (const part of parts) {
                    keys.add(part);
                }
            }
        }

        // let lastCachedMarshalled = ``;

        const hasLocal = ctx.modifiers.has("local");
        const hasSession = ctx.modifiers.has("session");
        const hasAny = hasLocal || hasSession;
        const useRemote = ctx.modifiers.has("remote");

        const storeUpdateHandler = ((_: CustomEvent<DatastarEvent>) => {
            // console.log("CacheStoreAttributePlugin storeUpdateHandler ");

            if (!hasAny) return;

            let s = ctx.store();
            if (useRemote) {
                s = remoteSignals(s);
            }
            if (keys.size > 0) {
                const newStore: Record<string, any> = {};
                for (const key of keys) {
                    const parts = key.split(".");

                    // if nested make sure the parent exists in newStore then fill in the value from the store
                    let newSubstore = newStore;
                    let subStore = s;
                    for (let i = 0; i < parts.length - 1; i++) {
                        const part = parts[i];
                        if (!newSubstore[part]) {
                            newSubstore[part] = {};
                        }
                        newSubstore = newSubstore[part];
                        subStore = subStore[part];
                    }

                    const lastPart = parts[parts.length - 1];
                    newSubstore[lastPart] = subStore[lastPart];
                }
                s = newStore;
            }

            const marshalledStore = JSON.stringify(s);

            // if (marshalledStore === lastCachedMarshalled) return;

            if (hasLocal) {
                window.localStorage.setItem(key, marshalledStore);
            }
            if (hasSession) {
                window.sessionStorage.setItem(key, marshalledStore);
            }

            // lastCachedMarshalled = marshalledStore;
        }) as EventListener;

        if (hasAny) {
            window.addEventListener(DATASTAR_EVENT, storeUpdateHandler);

            let marshalledStore: string | null = null;

            if (hasLocal) {
                marshalledStore = window.localStorage.getItem(key);
            } else if (hasSession) {
                marshalledStore = window.sessionStorage.getItem(key);
            }

            if (!!marshalledStore) {
                const store = JSON.parse(marshalledStore);
                for (const key in store) {
                    const value = store[key];
                    ctx.upsertIfMissingFromStore(key, value);
                }
            }
        }

        return () => {
            if (hasAny) {
                window.removeEventListener(DATASTAR_EVENT, storeUpdateHandler);
            }
        };
    },
};
