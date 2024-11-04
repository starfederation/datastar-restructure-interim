import {
    AttributeContext,
    AttributePlugin,
    DATASTAR,
    DATASTAR_EVENT,
    DatastarEvent,
    RegexpGroups,
} from "library/engine";
import { storeFromPossibleContents } from "library/utils/signals";

// Setup the global store
export const StoreAttributePlugin: AttributePlugin = {
    pluginType: "attribute",
    prefix: "store",
    removeNewLines: true,
    preprocessors: {
        pre: [
            {
                pluginType: "preprocessor",
                name: "store",
                regexp: /(?<whole>.+)/g,
                replacer: (groups: RegexpGroups) => {
                    const { whole } = groups;
                    return `Object.assign({...ctx.store()}, ${whole})`;
                },
            },
        ],
    },
    allowedModifiers: new Set(["local", "session", "ifmissing"]),
    onLoad: (ctx: AttributeContext) => {
        let lastCachedMarshalled = ``;
        const localFn = ((_: CustomEvent<DatastarEvent>) => {
            const s = ctx.store();
            const marshalledStore = JSON.stringify(s);

            if (marshalledStore !== lastCachedMarshalled) {
                window.localStorage.setItem(DATASTAR, marshalledStore);
                lastCachedMarshalled = marshalledStore;
            }
        }) as EventListener;
        const hasLocal = ctx.modifiers.has("local");
        if (hasLocal) {
            window.addEventListener(DATASTAR_EVENT, localFn);
            const marshalledStore = window.localStorage.getItem(DATASTAR) ||
                "{}";
            const store = JSON.parse(marshalledStore);
            ctx.mergeStore(store);
        }

        const hasSession = ctx.modifiers.has("session");
        const sessionFn = ((_: CustomEvent<DatastarEvent>) => {
            const s = ctx.store();
            const marshalledStore = JSON.stringify(s);
            window.sessionStorage.setItem(DATASTAR, marshalledStore);
        }) as EventListener;
        if (hasSession) {
            window.addEventListener(DATASTAR_EVENT, sessionFn);
            const marshalledStore = window.sessionStorage.getItem(DATASTAR) ||
                "{}";
            const store = JSON.parse(marshalledStore);
            ctx.mergeStore(store);
        }

        const possibleMergeStore = ctx.expressionFn(ctx);
        const actualMergeStore = storeFromPossibleContents(
            ctx.store(),
            possibleMergeStore,
            ctx.modifiers.has("ifmissing"),
        );
        ctx.mergeStore(actualMergeStore);

        delete ctx.el.dataset[ctx.rawKey];

        return () => {
            if (hasLocal) {
                window.removeEventListener(DATASTAR_EVENT, localFn);
            }

            if (hasSession) {
                window.removeEventListener(DATASTAR_EVENT, sessionFn);
            }
        };
    },
};
