import { HTMLorSVGElement } from "../utils/types";
import { DeepState } from "../vendored/deepsignal";
import { ReadonlySignal, Signal } from "../vendored/preact-core";

export type InitExpressionFunction = (
  ctx: InitContext,
  ...args: any
) => any;
export type AttribtueExpressionFunction = (
  ctx: AttributeContext,
  ...args: any
) => any;
export type Reactivity = {
  signal: <T>(value: T) => Signal<T>;
  computed: <T>(fn: () => T) => ReadonlySignal<T>;
  effect: (cb: () => void) => OnRemovalFn;
};

export type InitContext = {
  store: any;
  mergeStore: (store: DeepState) => void;
  actions: Readonly<ActionPlugins>;
  reactivity: Reactivity;
  applyPlugins: (target: Element) => void;
  cleanupElementRemovals: (el: Element) => void;
  sendDatastarEvent: SendDatastarEvent;
};

export type AttributeContext = InitContext & {
  upsertIfMissingFromStore: (path: string, value: any) => void;
  removeFromStore: (...paths: string[]) => void;
  walkSignals: (cb: (name: string, signal: Signal<any>) => void) => void;
  el: Readonly<HTMLorSVGElement>;
  key: Readonly<string>;
  rawKey: Readonly<string>;
  rawExpression: Readonly<string>;
  expression: Readonly<string>;
  expressionFn: AttribtueExpressionFunction;
  modifiers: Map<string, string[]>;
};

export type OnRemovalFn = () => void;

export interface DatastarPlugin {
  pluginType: "preprocessor" | "attribute" | "effect" | "action"; // The type of plugin
  name: string; // The name of the plugin
  requiredPlugins?: Set<DatastarPlugin>; // If not provided, no plugins are required
}

// A plugin accesible via a `data-${name}` attribute on an element
export interface AttributePlugin extends DatastarPlugin {
  pluginType: "attribute";
  onGlobalInit?: (ctx: InitContext) => void; // Called once on registration of the plugin
  onLoad: (ctx: AttributeContext) => OnRemovalFn | void; // Return a function to be called on removal
  allowedModifiers?: Set<string>; // If not provided, all modifiers are allowed
  mustHaveEmptyExpression?: boolean; // The contents of the data-* attribute must be empty
  mustNotEmptyExpression?: boolean; // The contents of the data-* attribute must not be empty
  mustHaveEmptyKey?: boolean; // The key of the data-* attribute must be empty after the prefix
  mustNotEmptyKey?: boolean; // The key of the data-* attribute must not be empty after the prefix
  allowedTagRegexps?: Set<string>; // If not provided, all tags are allowed
  disallowedTags?: Set<string>; // If not provided, no tags are disallowed
  preprocessors?: {
    pre?: PreprocessorPlugin[];
    post?: PreprocessorPlugin[];
  };
  removeNewLines?: boolean; // If true, the expression is not split by commas
  bypassExpressionFunctionCreation?: (ctx: AttributeContext) => boolean; // If true, the expression function is not created
  argumentNames?: Readonly<string[]>; // The names of the arguments passed to the expression function
}

export type RegexpGroups = Record<string, string>;

// A plugin that runs on the global scope that can effect the contents of a Datastar expression
export interface PreprocessorPlugin extends DatastarPlugin {
  pluginType: "preprocessor";
  regexp: RegExp;
  replacer: (groups: RegexpGroups) => string;
}

export type PreprocessorPlugins = Record<string, PreprocessorPlugin>;

export type ActionMethod = (ctx: AttributeContext, ...args: any[]) => any;

export interface ActionPlugin extends DatastarPlugin {
  pluginType: "action";
  method: ActionMethod;
}

export type ActionPlugins = Record<string, ActionPlugin>;

// A plugin that runs on the global scope of the DastaStar instance
export interface EffectPlugin extends DatastarPlugin {
  pluginType: "effect";
  onGlobalInit?: (ctx: InitContext) => void;
}

export type SendDatastarEvent = (
  category: "core" | "plugin",
  subcategory: string,
  type: string,
  target: Element | Document | Window | string,
  message: string,
  opts?: CustomEventInit,
) => void;

export interface DatastarEvent {
  time: Date;
  category: "core" | "plugin";
  subcategory: string;
  type: string;
  target: string;
  message: string;
}
