import {
    AttributePlugin,
    DATASTAR,
    DATASTAR_EVENT,
    DatastarEvent,
} from "library/src/engine";
import { remoteSignals } from "library/src/utils/signals";

export const PersistAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "persist",
    allowedModifiers: new Set(["local", "session", "remote"]),

    onLoad: (ctx) => {
        const key = ctx.key || DATASTAR;
        const expression = ctx.expression;
        const keys = new Set<string>();

        if (expression.trim() !== "") {
            const value = ctx.expressionFn(ctx);
            const parts = value.split(" ");
            for (const part of parts) {
                keys.add(part);
            }
        }

        let lastMarshalled = "";
        const storageType = ctx.modifiers.has("session") ? "session" : "local";
        const useRemote = ctx.modifiers.has("remote");

        const storeUpdateHandler = ((_: CustomEvent<DatastarEvent>) => {
            let store = ctx.store();
            if (useRemote) {
                store = remoteSignals(store);
            }
            if (keys.size > 0) {
                const newStore: Record<string, any> = {};
                for (const key of keys) {
                    let newSubstore = newStore;
                    let subStore = store;
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
                store = newStore;
            }

            const marshalledStore = JSON.stringify(store);

            if (marshalledStore === lastMarshalled) {
                return;
            }

            if (storageType === "session") {
                window.sessionStorage.setItem(key, marshalledStore);
            } else {
                window.localStorage.setItem(key, marshalledStore);
            }

            lastMarshalled = marshalledStore;
        }) as EventListener;

        window.addEventListener(DATASTAR_EVENT, storeUpdateHandler);

        let marshalledStore: string | null = null;

        if (storageType === "session") {
            marshalledStore = window.sessionStorage.getItem(key);
        } else {
            marshalledStore = window.localStorage.getItem(key);
        }

        if (!!marshalledStore) {
            const store = JSON.parse(marshalledStore);
            for (const key in store) {
                const value = store[key];
                ctx.upsertIfMissingFromStore(key, value);
            }
        }

        return () => {
            window.removeEventListener(DATASTAR_EVENT, storeUpdateHandler);
        };
    },
};
