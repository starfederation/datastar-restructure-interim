import { nodeHTMLorSVGElement } from "../utils/dom";
import { HTMLorSVGElement } from "../utils/types";
import { DeepSignal, deepSignal, DeepState } from "../vendored/deepsignal";
import { computed, effect, Signal, signal } from "../vendored/preact-core";
import { apply } from "../vendored/ts-merge-patch";
import {
    PLUGIN_ACTION,
    PLUGIN_ATTRIBUTE,
    PLUGIN_PREPROCESSOR,
    PLUGIN_WATCHER,
} from "./client_only_consts";
import {
    ActionPlugin,
    ActionPlugins,
    AttribtueExpressionFunction,
    AttributeContext,
    AttributePlugin,
    DatastarPlugin,
    InitContext,
    OnRemovalFn,
    PreprocessorPlugin,
    Reactivity,
    WatcherPlugin,
} from "./types";
import { VERSION } from "./version";

const isPreprocessorPlugin = (p: DatastarPlugin): p is PreprocessorPlugin =>
    p.pluginType === PLUGIN_PREPROCESSOR;
const isWatcherPlugin = (p: DatastarPlugin): p is WatcherPlugin =>
    p.pluginType === PLUGIN_WATCHER;
const isAttributePlugin = (p: DatastarPlugin): p is AttributePlugin =>
    p.pluginType === PLUGIN_ATTRIBUTE;
const isActionPlugin = (p: DatastarPlugin): p is ActionPlugin =>
    p.pluginType === PLUGIN_ACTION;

const alreadyExistsErr = (type: string, name: string) =>
    new Error(`A ${type} named '${name}' already exists`);
export class Engine {
    plugins: AttributePlugin[] = [];
    store: DeepSignal<any> = deepSignal({ _dsPlugins: {} });
    preprocessors = new Array<PreprocessorPlugin>();
    actions: ActionPlugins = {};
    watchers = new Array<WatcherPlugin>();
    refs: Record<string, HTMLElement> = {};
    reactivity: Reactivity = {
        signal,
        computed,
        effect,
    };
    parentID = "";
    missingIDNext = 0;
    removals = new Map<Element, { id: string; set: Set<OnRemovalFn> }>();
    mergeRemovals = new Array<OnRemovalFn>();

    constructor() {
        // const observer = new MutationObserver(
        //     (_mutationList, _observer) => {
        //         this.sendDatastarEvent(
        //             "core",
        //             "dom",
        //             "mutation",
        //             document.body,
        //             document.body.outerHTML,
        //         );
        //     },
        // );

        // Start observing the target node for configured mutations
        // observer.observe(document.body, {
        //     attributes: true,
        //     childList: true,
        //     subtree: true,
        // });
    }

    get version() {
        return VERSION;
    }

    load(...pluginsToLoad: DatastarPlugin[]) {
        const allLoadedPlugins = new Set<DatastarPlugin>(this.plugins);

        pluginsToLoad.forEach((plugin) => {
            if (plugin.requiredPlugins) {
                for (
                    const requiredPluginType of plugin
                        .requiredPlugins
                ) {
                    if (
                        !allLoadedPlugins.has(requiredPluginType)
                    ) {
                        throw new Error(
                            `Plugin '${plugin.name}' requires plugin '${requiredPluginType}' to be loaded`,
                        );
                    }
                }
            }

            let globalInitializer: ((ctx: InitContext) => void) | undefined;
            if (isPreprocessorPlugin(plugin)) {
                if (this.preprocessors.includes(plugin)) {
                    throw alreadyExistsErr("Preprocessor", plugin.name);
                }
                this.preprocessors.push(plugin);
            } else if (isWatcherPlugin(plugin)) {
                if (this.watchers.includes(plugin)) {
                    throw alreadyExistsErr("Watcher", plugin.name);
                }
                this.watchers.push(plugin);
                globalInitializer = plugin.onGlobalInit;
            } else if (isActionPlugin(plugin)) {
                if (!!this.actions[plugin.name]) {
                    throw alreadyExistsErr("Action", plugin.name);
                }
                this.actions[plugin.name] = plugin;
            } else if (isAttributePlugin(plugin)) {
                if (this.plugins.includes(plugin)) {
                    throw alreadyExistsErr("Attribute", plugin.name);
                }
                this.plugins.push(plugin);
                globalInitializer = plugin.onGlobalInit;
            } else {
                throw new Error(`Unknown plugin type: ${plugin}`);
            }

            if (globalInitializer) {
                globalInitializer({
                    store: this.store,
                    upsertIfMissingFromStore: this.upsertIfMissingFromStore
                        .bind(this),
                    mergeSignals: this.mergeSignals.bind(this),
                    removeSignals: this.removeSignals.bind(this),
                    actions: this.actions,
                    reactivity: this.reactivity,
                    applyPlugins: this.applyPlugins.bind(this),
                    cleanupElementRemovals: this.cleanupElementRemovals.bind(
                        this,
                    ),
                    // sendDatastarEvent: this.sendDatastarEvent.bind(this),
                });
            }

            // this.sendDatastarEvent(
            //     "core",
            //     "plugins",
            //     "registration",
            //     "BODY",
            //     `On prefix ${plugin.name}`,
            // );

            allLoadedPlugins.add(plugin);
        });

        this.applyPlugins(document.body);
    }

    // private sendDatastarEvent(
    //     category: "core" | "plugin",
    //     subcategory: string,
    //     type: string,
    //     target: Element | Document | Window | string,
    //     message: string,
    //     opts: CustomEventInit = {
    //         bubbles: true,
    //         cancelable: true,
    //         composed: true,
    //     },
    // ) {
    //     const contents = Object.assign(
    //         {
    //             detail: {
    //                 time: new Date(),
    //                 category,
    //                 subcategory,
    //                 type,
    //                 target: elemToSelector(target),
    //                 message,
    //             },
    //         },
    //         opts,
    //     );
    //     const evt = new CustomEvent<DatastarEvent>(DATASTAR_EVENT, contents);
    //     // console.log("Sending Datastar event", evt);
    //     window.dispatchEvent(evt);
    // }

    private cleanupElementRemovals(element: Element) {
        const removalSet = this.removals.get(element);
        if (removalSet) {
            for (const removal of removalSet.set) {
                removal();
            }
            this.removals.delete(element);
        }
    }

    lastMarshalledStore = "";
    private mergeSignals<T extends object>(mergeSignals: T) {
        this.mergeRemovals.forEach((removal) => removal());
        this.mergeRemovals = this.mergeRemovals.slice(0);

        const revisedStore = apply(this.store.value, mergeSignals) as DeepState;
        this.store = deepSignal(revisedStore);

        const marshalledStore = JSON.stringify(this.store.value);
        if (marshalledStore === this.lastMarshalledStore) return;

        // this.sendDatastarEvent(
        //     "core",
        //     "store",
        //     "merged",
        //     "STORE",
        //     marshalledStore,
        // );
    }

    private removeSignals(...keys: string[]) {
        const revisedStore = { ...this.store.value };
        for (const key of keys) {
            const parts = key.split(".");
            let currentID = parts[0];
            let subStore = revisedStore;
            for (let i = 1; i < parts.length; i++) {
                const part = parts[i];
                if (!subStore[currentID]) {
                    subStore[currentID] = {};
                }
                subStore = subStore[currentID];
                currentID = part;
            }
            delete subStore[currentID];
        }
        this.store = deepSignal(revisedStore);
        this.applyPlugins(document.body);
    }

    private upsertIfMissingFromStore(path: string, value: any) {
        const parts = path.split(".");
        let subStore = this.store as any;
        for (let i = 0; i < parts.length - 1; i++) {
            const part = parts[i];
            if (!subStore[part]) {
                subStore[part] = {};
            }
            subStore = subStore[part];
        }
        const last = parts[parts.length - 1];
        if (!!subStore[last]) return;
        subStore[last] = this.reactivity.signal(value);
        // this.sendDatastarEvent("core", "store", "upsert", path, value);
    }

    signalByName<T>(name: string) {
        return (this.store as any)[name] as Signal<T>;
    }

    private applyPlugins(rootElement: Element) {
        const appliedProcessors = new Set<PreprocessorPlugin>();

        this.plugins.forEach((p, pi) => {
            this.walkDownDOM(rootElement, (el) => {
                if (!pi) this.cleanupElementRemovals(el);

                for (const rawKey in el.dataset) {
                    const rawExpression = `${el.dataset[rawKey]}` || "";
                    let expression = rawExpression;

                    if (!rawKey.startsWith(p.name)) continue;

                    if (el.id.length === 0) {
                        el.id = `ds-${this.parentID}-${this.missingIDNext++}`;
                    }

                    appliedProcessors.clear();

                    if (p.allowedTagRegexps) {
                        const lowerCaseTag = el.tagName.toLowerCase();
                        const allowed = [...p.allowedTagRegexps].some((r) =>
                            lowerCaseTag.match(r)
                        );
                        if (!allowed) {
                            throw new Error(
                                `'${el.tagName}' not allowed for '${rawKey}', allowed ${
                                    [
                                        [...p.allowedTagRegexps].map((t) =>
                                            `'${t}'`
                                        ),
                                    ].join(", ")
                                }`,
                            );
                        }
                    }

                    let keyRaw = rawKey.slice(p.name.length);
                    let [key, ...modifiersWithArgsArr] = keyRaw.split(".");
                    if (p.mustHaveEmptyKey && key.length > 0) {
                        throw new Error(`'${rawKey}' must have empty key`);
                    }
                    if (p.mustNotEmptyKey && key.length === 0) {
                        throw new Error(`'${rawKey}' must have non-empty key`);
                    }
                    if (key.length) {
                        key = key[0].toLowerCase() + key.slice(1);
                    }

                    const modifiersArr = modifiersWithArgsArr.map((m) => {
                        const [label, ...args] = m.split("_");
                        return { label, args };
                    });
                    if (p.allowedModifiers) {
                        for (const modifier of modifiersArr) {
                            if (!p.allowedModifiers.has(modifier.label)) {
                                throw new Error(
                                    `'${modifier.label}' is not allowed`,
                                );
                            }
                        }
                    }
                    const modifiers = new Map<string, string[]>();
                    for (const modifier of modifiersArr) {
                        modifiers.set(modifier.label, modifier.args);
                    }

                    if (p.mustHaveEmptyExpression && expression.length) {
                        throw new Error(
                            `'${rawKey}' must have empty expression`,
                        );
                    }
                    if (p.mustNotEmptyExpression && !expression.length) {
                        throw new Error(
                            `'${rawKey}' must have non-empty expression`,
                        );
                    }

                    const splitRegex = /;|\n/;

                    if (p.removeNewLines) {
                        expression = expression
                            .split("\n")
                            .map((p: string) => p.trim())
                            .join(" ");
                    }

                    const processors = [
                        ...(p.preprocessors?.pre || []),
                        ...this.preprocessors,
                        ...(p.preprocessors?.post || []),
                    ];
                    for (const processor of processors) {
                        if (appliedProcessors.has(processor)) continue;
                        appliedProcessors.add(processor);

                        const expressionParts = expression.split(splitRegex);
                        const revisedParts: string[] = [];

                        expressionParts.forEach((exp) => {
                            let revised = exp;
                            const matches = [
                                ...revised.matchAll(processor.regexp),
                            ];
                            if (matches.length) {
                                for (const match of matches) {
                                    if (!match.groups) continue;
                                    const { groups } = match;
                                    const { whole } = groups;
                                    revised = revised.replace(
                                        whole,
                                        processor.replacer(groups),
                                    );
                                }
                            }
                            revisedParts.push(revised);
                        });
                        // })

                        expression = revisedParts.join("; ");
                    }

                    const ctx: AttributeContext = {
                        store: () => this.store,
                        mergeSignals: this.mergeSignals.bind(this),
                        upsertIfMissingFromStore: this.upsertIfMissingFromStore
                            .bind(this),
                        removeSignals: this.removeSignals.bind(this),
                        applyPlugins: this.applyPlugins.bind(this),
                        cleanupElementRemovals: this.cleanupElementRemovals
                            .bind(this),
                        walkSignals: this.walkSignals.bind(this),
                        actions: this.actions,
                        reactivity: this.reactivity,
                        el,
                        rawKey,
                        key,
                        rawExpression,
                        expression,
                        expressionFn: () => {
                            throw new Error("Expression function not created");
                        },
                        modifiers,
                        // sendDatastarEvent: this.sendDatastarEvent.bind(this),
                    };

                    if (
                        !p.bypassExpressionFunctionCreation?.(ctx) &&
                        !p.mustHaveEmptyExpression && expression.length
                    ) {
                        const statements = expression
                            .split(splitRegex)
                            .map((s) => s.trim())
                            .filter((s) => s.length);
                        statements[statements.length - 1] = `return ${
                            statements[statements.length - 1]
                        }`;
                        const joined = statements.map((s) => `  ${s}`).join(
                            ";\n",
                        );
                        const fnContent = `
  try {
    const _datastarExpression = () => {
  ${joined}
    }
    const _datastarReturnVal = _datastarExpression()
    return _datastarReturnVal
  } catch (e) {
   const msg = \`
  Error evaluating Datastar expression:
  ${joined.replaceAll("`", "\\`")}

  Error: \${e.message}

  Check if the expression is valid before raising an issue.
  \`.trim()
   console.error(msg)
   debugger
  }
              `;

                        /*sendDatastarEvent(
                'core',
                'attributes',
                'expr_construction',
                ctx.el,
                `${rawKey}="${rawExpression}" becomes: ${joined}`,
              )*/
                        try {
                            const argumentNames = p.argumentNames || [];
                            const fn = new Function(
                                "ctx",
                                ...argumentNames,
                                fnContent,
                            ) as AttribtueExpressionFunction;
                            ctx.expressionFn = fn;
                        } catch (e) {
                            const err = new Error(
                                `Error creating expression function for '${fnContent}', error: ${e}`,
                            );
                            // this.sendDatastarEvent(
                            //     "core",
                            //     "attributes",
                            //     "expr_construction_err",
                            //     ctx.el,
                            //     String(err),
                            // );
                            console.error(err);
                            debugger;
                        }
                    }

                    const removal = p.onLoad(ctx);
                    if (removal) {
                        if (!this.removals.has(el)) {
                            this.removals.set(el, {
                                id: el.id,
                                set: new Set(),
                            });
                        }
                        this.removals.get(el)!.set.add(removal);
                    }
                }
            });
        });
    }

    private walkSignalsStore(
        store: any,
        callback: (name: string, signal: Signal<any>) => void,
    ) {
        const keys = Object.keys(store);
        for (let i = 0; i < keys.length; i++) {
            const key = keys[i];
            const value = store[key];
            const isSignal = value instanceof Signal;
            const hasChildren = typeof value === "object" &&
                Object.keys(value).length > 0;

            if (isSignal) {
                callback(key, value);
                continue;
            }

            if (!hasChildren) continue;

            this.walkSignalsStore(value, callback);
        }
    }

    private walkSignals(callback: (name: string, signal: Signal<any>) => void) {
        this.walkSignalsStore(this.store, callback);
    }

    private walkDownDOM(
        element: Element | null,
        callback: (el: HTMLorSVGElement) => void,
        siblingOffset = 0,
    ) {
        if (!element) return;
        const el = nodeHTMLorSVGElement(element);
        if (!el) return;

        callback(el);

        siblingOffset = 0;
        element = element.firstElementChild;
        while (element) {
            this.walkDownDOM(element, callback, siblingOffset++);
            element = element.nextElementSibling;
        }
    }
}
